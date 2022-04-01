package data

import (
	"context"
	"fmt"
	"pd-srv/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type shopRepo struct {
	data   *Data
	logger *log.Helper
}

func NewShopRepo(data *Data, logger log.Logger) biz.ShopRepo {
	return &shopRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/shop")),
	}
}

func (sr *shopRepo) CreateShop(ctx context.Context, s *biz.Shop) (*biz.Shop, error) {
	err := sr.data.maria.Create(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (sr *shopRepo) UpdateShop(ctx context.Context, s *biz.Shop) (*biz.Shop, error) {
	err := sr.data.maria.Model(&biz.Shop{}).Where("id=?", s.Id).Update(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (sr *shopRepo) DeleteShop(ctx context.Context, s *biz.Shop) (*biz.Shop, error) {
	s.DeleteAt = time.Now().Format("2006-01-02:15-04")
	s.IsDeleted = 1
	err := sr.data.maria.Model(&biz.Shop{}).Where("id=?", s.Id).Update(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (sr *shopRepo) GetShop(ctx context.Context, id int) (*biz.Shop, error) {
	var s = new(biz.Shop)
	err := sr.data.maria.Model(&biz.Shop{}).Where("id=?", id).First(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (sr *shopRepo) ListShop(ctx context.Context, q *biz.Query) ([]biz.Shop, error) {
	var shopList []biz.Shop
	if len(q.Name) != 0 {
		sr.data.maria = sr.data.maria.Model(&biz.Shop{}).Where("name like ?", fmt.Sprintf("%%%s%%", q.Name))
	}
	err := sr.data.maria.Model(&biz.Shop{}).Offset(q.Offset).Limit(q.Limit).Find(&shopList).Error
	if err != nil {
		return nil, err

	}

	return shopList, nil

}

func (sr *shopRepo) GetSpByUuid(ctx context.Context, sp *biz.Shop) (*biz.Shop, error) {
	err := sr.data.maria.Model(&biz.Shop{}).Where("user_uuid=?", sp.UserUuid).First(sp).Error
	if err != nil {
		return nil, err
	}
	return sp, nil
}
