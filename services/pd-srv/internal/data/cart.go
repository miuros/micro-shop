package data

import (
	"context"
	"pd-srv/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type cartRepo struct {
	data   *Data
	logger *log.Helper
}

func (cr *cartRepo) CreateCart(ctx context.Context, ca *biz.Cart) (*biz.Cart, error) {
	err := cr.data.maria.Create(ca).Error
	if err != nil {
		return nil, err
	}
	return ca, err
}

func (cr *cartRepo) UpdateCart(ctx context.Context, ca *biz.Cart) (*biz.Cart, error) {
	err := cr.data.maria.Model(&biz.Cart{}).Where("id=? and user_uuid=?", ca.Id, ca.UserUuid).Update(ca).Error
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (cr *cartRepo) DeleteCart(ctx context.Context, ca *biz.Cart) (*biz.Cart, error) {
	err := cr.data.maria.Model(&biz.Cart{}).Where("id=?  and user_uuid=?", ca.Id, ca.UserUuid).Delete(ca).Error
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (cr *cartRepo) GetCart(ctx context.Context, ca *biz.Cart) (*biz.Cart, error) {
	err := cr.data.maria.Model(&biz.Cart{}).Where("id=? and user_uuid=?", ca.Id, ca.UserUuid).First(ca).Error
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (cr *cartRepo) ListCart(ctx context.Context, q *biz.QueryCart) ([]biz.Cart, error) {
	var cartList []biz.Cart
	err := cr.data.maria.Model(&biz.Cart{}).Where("user_uuid=?", q.UserUuid).Offset(q.Offset).Limit(q.Limit).Find(&cartList).Error
	if err != nil {
		return nil, err
	}
	return cartList, nil
}

func NewCartRepo(data *Data, logger log.Logger) biz.CartRepo {
	return &cartRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/cart")),
	}
}
