package data

import (
	"comment-srv/internal/biz"
	"comment-srv/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCmRepo)

// Data .
type Data struct {
	maria *gorm.DB
	redis redis.Cmdable
}

func autoMigrate(db *gorm.DB) *gorm.DB {
	db = db.AutoMigrate(&biz.Comment{})
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

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	maria, err := newMaria(c)
	if err != nil {
		return nil, nil, err
	}
	redisCli := newRedis(c)
	cleanup := func() {
		maria.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		maria: maria,
		redis: redisCli,
	}, cleanup, nil
}
