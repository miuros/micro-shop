package data

import (
	"context"
	"pd-srv/internal/biz"
	"pd-srv/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewBannerRepo, NewProductRepo, NewCartRepo, NewShopRepo, NewCgRepo)

// Data .
type Data struct {
	maria *gorm.DB
	redis redis.Cmdable
}

func newMaria(c *conf.Data) (*gorm.DB, error) {
	db, err := gorm.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db = migrate(db)
	return db, nil
}

func migrate(db *gorm.DB) *gorm.DB {
	return db.AutoMigrate(&biz.Product{}, &biz.Cart{}, &biz.Shop{}, &biz.Banner{}, &biz.Category{})
}

func newRedis(c *conf.Data) redis.Cmdable {
	cli := redis.NewClient(&redis.Options{
		DB:           0,
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			return cn.Ping(ctx).Err()
		},
	})

	return cli
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	maria, err := newMaria(c)
	if err != nil {
		return nil, nil, err
	}
	redisCli := newRedis(c)
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		maria.Close()
	}
	return &Data{
		maria: maria,
		redis: redisCli,
	}, cleanup, nil
}
