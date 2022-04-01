package data

import (
	"context"
	"fmt"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	ggrpc "google.golang.org/grpc"
	"time"
	odv1 "user-srv/api/od/v1"
	"user-srv/internal/biz"
	"user-srv/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewAddressRepo)

// Data .
type Data struct {
	maria  *gorm.DB
	rabbit *amqp.Connection
	redis  redis.Cmdable
	odSrv  odv1.OdSrvClient
	mail   string
	host   string
	passwd string
}

func newMaria(c *conf.Data) (*gorm.DB, error) {
	db, err := gorm.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	return db, nil

}

func newRabbit(c *conf.Data) (*amqp.Connection, error) {
	return amqp.Dial(c.Rabbit.Addr)
}

func newRedisCmdable(c *conf.Data) redis.Cmdable {
	redisCli := redis.NewClient(&redis.Options{
		Network: c.Redis.Network,
		Addr:    c.Redis.Addr,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			select {
			case <-ctx.Done():
				return fmt.Errorf("context had canceled")
			default:
				return cn.Ping(ctx).Err()
			}
		},
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		PoolSize:     10,
		DialTimeout:  time.Second * 2,
	})

	return redisCli
}

func Migrate(db *gorm.DB) *gorm.DB {
	return db.AutoMigrate(&biz.User{}, &biz.AddressInfo{})
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	maria, err := newMaria(c)
	if err != nil {
		return nil, nil, err
	}
	maria = Migrate(maria)
	rabbit, err := newRabbit(c)
	if err != nil {
		return nil, nil, err
	}

	redisCli := newRedisCmdable(c)

	client, err := api.NewClient(&api.Config{
		Address: c.Consul.Addr,
	})
	if err != nil {
		panic(err)
	}
	src, _, err := client.Agent().Service(c.OdEndpoint, nil)
	if err != nil {
		panic(err)
	}
	cli := consul.New(client)
	conn, err := grpc.Dial(context.Background(), grpc.WithEndpoint(fmt.Sprintf("%s:%d", src.Address, src.Port)), grpc.WithDiscovery(cli), grpc.WithOptions(ggrpc.WithInsecure()))
	if err != nil {
		panic(err)
	}
	odClient := odv1.NewOdSrvClient(conn)
	cleanup := func() {
		_ = rabbit.Close()
		log.NewHelper(logger).Info("closing the data resources")
		_ = maria.Close()
	}
	return &Data{maria: maria, redis: redisCli, rabbit: rabbit, host: c.Host, passwd: c.Pwd, mail: c.User, odSrv: odClient}, cleanup, nil
}
