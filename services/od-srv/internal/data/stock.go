package data

import (
	"context"
	"fmt"
	"od-srv/internal/biz"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/gorm"
)

type scRepo struct {
	data   *Data
	logger *log.Helper
}

func (sr *scRepo) CreateSc(ctx context.Context, sc *biz.Stock) (*biz.Stock, error) {
	err := sr.data.maria.Transaction(func(tx *gorm.DB) error {
		return tx.Create(sc).Error
	})
	if err != nil {
		return nil, err
	}
	err = sr.updateRedisStock(ctx, sc)
	if err != nil {
		sr.data.maria.Rollback()

	}
	return sc, nil
}

func storageKey(productId int64) string {
	return "product_storage_" + strconv.FormatInt(productId, 10)
}

func saleKey(productId int64) string {
	return "product_sale_" + strconv.FormatInt(productId, 10)
}

func (sr *scRepo) updateRedisStock(ctx context.Context, sc *biz.Stock) error {
	pipe := sr.data.redis.TxPipeline()

	pipe.Set(ctx, storageKey(sc.ProductId), sc.Storage, -1)

	pipe.Set(ctx, saleKey(sc.ProductId), sc.Sale, -1)
	_, err := pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
	}
	return err
}

func (sr *scRepo) UpdateSc(ctx context.Context, sc *biz.Stock) (*biz.Stock, error) {
	err := sr.data.maria.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).Update(sc).Error
	})
	if err != nil {
		return nil, err
	}
	err = sr.updateRedisStock(ctx, sc)
	if err != nil {
		sr.data.maria.Rollback()
		return nil, err
	}
	return sc, nil
}

func (sr *scRepo) DeleteSc(ctx context.Context, sc *biz.Stock) (*biz.Stock, error) {
	err := sr.data.maria.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).Update(sc).Error
	})

	if err != nil {
		return nil, err
	}
	err = sr.deleteRedisStock(ctx, sc)
	if err != nil {
		sr.data.maria.Rollback()
	}
	return sc, nil
}

func (sr *scRepo) deleteRedisStock(ctx context.Context, sc *biz.Stock) error {
	pipe := sr.data.redis.TxPipeline()
	pipe.Del(ctx, storageKey(sc.ProductId))
	pipe.Del(ctx, saleKey(sc.ProductId))
	_, err := pipe.Exec(ctx)
	if err != nil {
		pipe.Discard()
	}
	return err
}

func (sr *scRepo) AddSc(ctx context.Context, sc *biz.Stock) (*biz.Stock, error) {
	err := sr.data.maria.Transaction(func(tx *gorm.DB) error {
		var old = new(biz.Stock)
		err := tx.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).First(old).Error
		if err != nil {
			return err
		}
		if old.IsDeleted == 1 {
			return fmt.Errorf("this product is deleted")
		}
		sc.Storage = sc.Storage + old.Storage
		err = tx.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).Update(sc).Error
		return err
	})
	if err != nil {
		return nil, err
	}
	err = sr.updateRedisStock(ctx, sc)
	if err != nil {
		sr.data.maria.Rollback()
		return nil, err
	}

	return sc, nil
}

func (sr *scRepo) GetSc(ctx context.Context, sc *biz.Stock) (*biz.Stock, error) {
	err := sr.data.maria.Model(&biz.Stock{}).Where("product_id=?", sc.ProductId).First(sc).Error
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func NewScRepo(data *Data, logger log.Logger) biz.ScRepo {
	return &scRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/stock")),
	}
}

var _ biz.ScRepo = (*scRepo)(nil)
