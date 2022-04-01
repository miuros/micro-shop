package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"od-srv/internal/biz"
	"od-srv/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCateRepo, NewOdRepo, NewScRepo)

// Data .
type Data struct {
	maria  *gorm.DB
	redis  redis.Cmdable
	rabbit *amqp.Connection
}

func autoMigrate(db *gorm.DB) *gorm.DB {
	db = db.AutoMigrate(&biz.Order{}, &biz.Cate{}, &biz.Stock{})
	return db
}

func newRedis(c *conf.Data) redis.Cmdable {
	var cli = redis.NewClient(&redis.Options{
		Network:      c.Redis.Network,
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			return cn.Ping(ctx).Err()
		},
	})
	return cli

}

func newMaria(cc *conf.Data) (*gorm.DB, error) {
	db, err := gorm.Open(cc.Database.Driver, cc.Database.Source)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db = autoMigrate(db)
	return db, nil
}

func newRabbit(c *conf.Data) (*amqp.Connection, error) {
	return amqp.Dial(c.Rabbit.Addr)
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	maria, err := newMaria(c)
	if err != nil {
		return nil, nil, err
	}
	redisCli := newRedis(c)
	rabbit, err := newRabbit(c)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		rabbit.Close()
		maria.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{
		maria:  maria,
		redis:  redisCli,
		rabbit: rabbit,
	}, cleanup, nil
}
