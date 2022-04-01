package data

import (
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"nt-srv/internal/biz"
	"nt-srv/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewNtRepo)

// Data .
type Data struct {
	maria *gorm.DB
}

func newMaria(c *conf.Data) (*gorm.DB, error) {
	db, err := gorm.Open(c.Database.Driver, c.Database.Source)

	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db.AutoMigrate(&biz.Notice{})
	return db, nil

}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	maria, err := newMaria(c)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		maria.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		maria: maria,
	}, cleanup, nil
}
