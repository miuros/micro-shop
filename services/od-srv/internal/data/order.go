package data

import (
	"context"
	"encoding/json"
	"fmt"
	"od-srv/internal/biz"
	"strconv"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

var (
	ERRStorageNotEnough = fmt.Errorf("storage is not enough")
)

type odRepo struct {
	data    *Data
	logger  *log.Helper
	queue   chan *biz.Order
	channel *amqp.Channel
	locker  sync.Mutex
}

func (or *odRepo) ListOdForShopper(ctx context.Context, query *biz.QueryOdForShopper) ([]biz.Order, error) {
	var odList []biz.Order
	var err error
	if query.Status != 0 {
		err = or.data.maria.Select("`order`.*").Joins("left join product on `order`.product_id=product.id").Where("product.shop_id=? and status=?", query.ShopId, query.Status).Offset(query.Offset).Limit(query.Limit).Find(&odList).Error
	} else {
		err = or.data.maria.Select("`order`.*").Joins("left join product on `order`.product_id=product.id").Where("product.shop_id=?", query.ShopId).Offset(query.Offset).Limit(query.Limit).Find(&odList).Error
	}
	if err != nil {
		return nil, err
	}
	return odList, nil
}

func (or *odRepo) CreateOd(ctx context.Context, od *biz.Order) (*biz.Order, error) {
	or.locker.Lock()
	defer or.locker.Unlock()

	err := or.data.maria.Transaction(func(tx *gorm.DB) error {
		var sc = new(biz.Stock)
		err := tx.Model(&biz.Stock{}).Where("product_id=?", od.ProductId).First(sc).Error
		if err != nil {
			return err
		}
		if sc.Storage < od.Number {
			return ERRStorageNotEnough
		}
		sc.Storage -= od.Number
		err = tx.Model(&biz.Stock{}).Where("product_id=?", od.ProductId).Update(sc).Error
		if err != nil {
			return err
		}
		err = or.data.maria.Create(od).Error
		return err
	})
	if err != nil {
		return nil, err
	}
	/*
		pipe := or.data.redis.TxPipeline()
		storage, err := pipe.Get(ctx, storageKey(od.ProductId)).Int64()
		if err != nil {
			return nil, err
		}
		if storage < od.Number {
			return nil, ERRStorageNotEnough
		}
		err = pipe.DecrBy(ctx, storageKey(od.ProductId), int64(od.Number)).Err()
		if err != nil {
			return nil, err
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			pipe.Discard()
		}

		data, err := json.Marshal(od)
		if err != nil {
			//pipe.Discard()
			return nil, err
		}
		err = or.channel.Publish("", "create_order", true, true, amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
		if err != nil {
			//pipe.Discard()
			return nil, err
		}
	*/
	or.queue <- od
	return od, nil

}

func (or *odRepo) createOd(od *biz.Order) (*biz.Order, error) {
	/*
		pipe := or.data.redis.Pipeline()
		stock, err := pipe.Get(ctx, stockKey(od.ProductId)).Int()
		if err != nil {
			return nil, err
		}
		if stock < 1 {
			return nil, fmt.Errorf("%d 's stock is not enough", od.ProductId)
		}
		pipe.DecrBy(ctx, stockKey(od.ProductId), int64(od.Number))
	*/
	err := or.data.maria.Transaction(func(tx *gorm.DB) error {
		return or.data.maria.Create(od).Error
	})
	/*
		_, pipeErr := pipe.Exec(ctx)
		if pipeErr != nil {
			or.data.maria.Rollback()
			pipe.Discard()
			return nil, err
		}
	*/

	if err != nil {
		or.data.maria.Rollback()
		//pipe.Discard()
		return nil, err
	}
	or.queue <- od

	return od, err
}

func (or *odRepo) orderRabbitQueue() error {

	//_, err = channel.QueueDeclare("create_order", true, false, false, false, nil)
	//if err != nil {
	//	return err
	//}
	var args = amqp.Table{
		"x-dead-letter-exchange":    "order_exchange",
		"x-dead-letter-routing-key": "order",
	}
	err := or.channel.ExchangeDeclare("order_dead_exchange", amqp.ExchangeDirect, true, false, false, false, args)
	if err != nil {
		return err
	}
	deadQueue, err := or.channel.QueueDeclare("order_dead_queue", true, false, false, false, args)
	if err != nil {
		return err
	}
	err = or.channel.QueueBind(deadQueue.Name, "order", "order_dead_exchange", false, nil)
	if err != nil {
		return err
	}

	err = or.channel.ExchangeDeclare("order_exchange", amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return err
	}
	q, err := or.channel.QueueDeclare("order_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = or.channel.QueueBind(q.Name, "order", "order_exchange", false, nil)
	if err != nil {
		return err
	}

	var arg = amqp.Table{
		"x-dead-letter-exchange":    "mail_exchange",
		"x-dead-letter-routing-key": "mail",
	}
	err = or.channel.ExchangeDeclare("mail_dead_exchange", amqp.ExchangeDirect, true, false, false, false, arg)
	if err != nil {
		return err
	}
	mailDeadQueue, err := or.channel.QueueDeclare("mail_dead_queue", true, false, false, false, arg)
	if err != nil {
		return err
	}
	err = or.channel.QueueBind(mailDeadQueue.Name, "mail", "mail_dead_exchange", false, nil)
	if err != nil {
		return err
	}

	err = or.channel.ExchangeDeclare("mail_exchange", amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return err
	}
	mailQueue, err := or.channel.QueueDeclare("mail_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = or.channel.QueueBind(mailQueue.Name, "mail", "mail_exchange", false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case od := <-or.queue:
				data, err := json.Marshal(od)
				if err != nil {
					or.logger.Errorf("json marshal error:%s", err.Error())
					continue
				}
				err = or.channel.Publish("order_dead_exchange", "order", false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
					Expiration:  strconv.Itoa(3600000),
				})
				if err != nil {
					or.logger.Errorf("channel publish error:%s", err.Error())
				}
				err = or.channel.Publish("mail_dead_exchange", "mail", false, false, amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
					Expiration:  strconv.Itoa(3000000),
				})
				if err != nil {
					or.logger.Errorf("channel publish error:%s", err.Error())
				}
				//default:
			}
		}
	}()
	return nil
}

func (or *odRepo) rcvDatedOrder() error {
	msg, err := or.channel.Consume("order_queue", "", true, true, false, false, nil)
	if err != nil {
		return err
	}
	ctx := context.Background()
	go func() {
		for data := range msg {
			var od = new(biz.Order)
			err := json.Unmarshal(data.Body, od)
			if err != nil {
				or.logger.Errorf("json unmarshal error:%s", err.Error())
				continue
			}
			or.data.maria.Transaction(func(tx *gorm.DB) error {
				var res = new(biz.Order)
				err := tx.Model(&biz.Order{}).Where("id=?", od.Id).First(res).Error
				if err != nil {
					return err
				}
				if res.Status == 2 || res.IsDeleted == 1 {
					return nil
				}
				err = tx.Model(&biz.Order{}).Where("id=? and user_uuid=?", res.Id, res.UserUuid).Update(&biz.Order{Status: 3, IsDeleted: 1, DeleteAt: time.Now().Format("2006-01-02:15-04")}).Error
				if err != nil {
					return err
				}
				var sc = new(biz.Stock)
				err = tx.Model(&biz.Stock{}).Where("product_id=?", res.ProductId).First(sc).Error
				if err != nil {
					return err
				}
				sc.Storage += od.Number

				err = tx.Model(&biz.Stock{}).Where("product_id=?", res.ProductId).Update(sc).Error
				or.data.redis.Set(ctx, storageKey(od.ProductId), sc.Storage, -1)
				return err
			})
		}
	}()
	return nil
}

func (or *odRepo) deleteOdUpdateStock() {
	err := or.data.maria.Transaction(func(tx *gorm.DB) error {
		var odList []biz.Order
		err := tx.Model(&biz.Order{}).Where("create_at < ? and status =1 and is_deleted=0", time.Now().Add(-time.Hour).Format("2006-01-02:15-04")).Find(&odList).Error
		if err != nil {
			return err
		}
		for _, item := range odList {
			var sc = new(biz.Stock)
			err = tx.Model(&biz.Stock{}).Where("product_id=?", item.ProductId).First(&sc).Error
			if err != nil {
				return err
			}
			sc.Storage += item.Number
			err = tx.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).Update(sc).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		or.logger.Errorf("delete out dated order:%s", err.Error())
	}
}

func (or *odRepo) checkStatus() {
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			or.deleteOdUpdateStock()
			//default:

		}
	}
}

//update address id
func (or *odRepo) UpdateOd(ctx context.Context, od *biz.Order) (*biz.Order, error) {
	err := or.data.maria.Model(&biz.Order{}).Where("id=? and user_uuid=?", od.Id, od.UserUuid).Update(&biz.Order{UpdateAt: time.Now().Format("2006-01-02:15-04"), AddressId: od.AddressId}).Error
	if err != nil {
		return nil, err
	}

	return od, nil
}

func (or *odRepo) DeleteOd(ctx context.Context, od *biz.Order) (*biz.Order, error) {
	err := or.data.maria.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&biz.Order{}).Where("id=? and user_uuid=?", od.Id, od.UserUuid).Update(od).Error
		if err != nil {
			return err
		}
		var sc = new(biz.Stock)
		err = tx.Model(&biz.Stock{}).Where("product_id=?", od.ProductId).First(sc).Error
		if err != nil {
			return err
		}
		sc.Storage += od.Number

		err = tx.Model(&biz.Stock{}).Where("product_id=?", od.ProductId).Update(sc).Error
		or.data.redis.Set(ctx, storageKey(od.ProductId), sc.Storage, -1)
		return err
	})
	if err != nil {
		return nil, err
	}

	return od, nil
}

/*
func (or *odRepo) DeleteOd(ctx context.Context, cate *biz.Cate) ([]biz.Order, error) {
	var odList []biz.Order
	err := or.data.maria.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&biz.Order{}).Where("cate_id=?", cate.Id).Update(&biz.Order{IsDeleted: 1, DeleteAt: time.Now().String()}).Error
		if err != nil {
			return err
		}
		err = tx.Model(&biz.Order{}).Where("cate_id=?", cate.Id).Find(&odList).Error
		if err != nil {
			return err
		}
		for _, item := range odList {
			var sc = new(biz.Stock)

			err = tx.Model(&biz.Stock{}).Where("product_id=?", item.ProductId).First(sc).Error
			if err != nil {
				return err
			}
			sc.Storage += item.Number
			err = tx.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).Update(sc).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return odList, nil

}
*/

func (or *odRepo) GetOd(ctx context.Context, od *biz.Order) (*biz.Order, error) {
	err := or.data.maria.Model(&biz.Order{}).Where("id=? and user_uuid=?", od.Id, od.UserUuid).Find(od).Error
	if err != nil {
		return nil, err
	}

	return od, nil
}

func (or *odRepo) ListOd(ctx context.Context, qo *biz.QueryOd) ([]biz.Order, error) {
	var odList []biz.Order
	var err error

	if qo.Status > 0 {
		if len(qo.UserUuid) == 0 {

			err = or.data.maria.Model(&biz.Order{}).Where("status=?", qo.Status).Offset(qo.Offset).Limit(qo.Limit).Find(&odList).Error
		} else {
			err = or.data.maria.Model(&biz.Order{}).Where("user_uuid=? and status=?", qo.UserUuid, qo.Status).Offset(qo.Offset).Limit(qo.Limit).Find(&odList).Error
		}
	} else {
		if len(qo.UserUuid) == 0 {

			err = or.data.maria.Model(&biz.Order{}).Offset(qo.Offset).Limit(qo.Limit).Find(&odList).Error
		} else {
			err = or.data.maria.Model(&biz.Order{}).Where("user_uuid=? ", qo.UserUuid).Offset(qo.Offset).Limit(qo.Limit).Find(&odList).Error
		}
	}
	if err != nil {
		return nil, err
	}
	return odList, nil
}

func (or *odRepo) ListOdByCateId(ctx context.Context, cate *biz.Cate) ([]biz.Order, error) {
	var odList []biz.Order
	var err error
	if cate.Status == 0 {

		err = or.data.maria.Model(&biz.Order{}).Where("cate_id=? and user_uuid=? and status=?", cate.Id, cate.UserUuid, cate.Status).Find(&odList).Error
	} else {

		err = or.data.maria.Model(&biz.Order{}).Where("cate_id=? and user_uuid=?", cate.Id, cate.UserUuid).Find(&odList).Error
	}
	if err != nil {
		return nil, err
	}
	return odList, nil
}

func NewOdRepo(data *Data, logger log.Logger) biz.OdRepo {
	or := &odRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/order")),
		queue:  make(chan *biz.Order, 100),
		locker: sync.Mutex{},
	}

	channel, err := or.data.rabbit.Channel()
	if err != nil {
		panic(err)
	}
	or.channel = channel
	err = or.orderRabbitQueue()
	if err != nil {
		panic(err)
	}
	go or.checkStatus()
	err = or.rcvDatedOrder()
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	return or
}

/*
func (or *odRepo) rcvCreateOrder() error {
	channel, err := or.data.rabbit.Channel()
	if err != nil {
		return err
	}
	msg, err := channel.Consume("create_order", "", true, false, true, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for m := range msg {
			var od = new(biz.Order)
			err := json.Unmarshal(m.Body, od)
			if err != nil {
				or.logger.Errorf("json unmarshal error:%s", err.Error())
				continue
			}
			_, err = or.createOd(od)
			if err != nil {
				or.logger.Errorf("create order error:%s", err.Error())
				continue
			}

		}
	}()
	return nil
}
*/

func (or *odRepo) PayOd(ctx context.Context, od *biz.Order) error {
	err := or.data.maria.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&biz.Order{}).Where("id=? and user_uuid=?", od.Id, od.UserUuid).Update(&biz.Order{Status: 2, PayTime: time.Now().Format("2006-01-02:15-04"), PayType: 1}).Error
		if err != nil {
			return err
		}
		var sc = new(biz.Stock)
		err = tx.Model(&biz.Stock{}).Where("product_id=?", od.ProductId).First(sc).Error
		if err != nil {
			return err
		}
		sc.Sale += od.Number
		err = tx.Model(&biz.Stock{}).Where("product_id=?", od.ProductId).Update(sc).Error
		return err
	})
	return err
}
