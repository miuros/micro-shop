package service

import (
	"context"
	v1 "od-srv/api/order/v1"
)

func (os *OdSrv) GetPd(ctx context.Context, id int64) (*v1.Product, error) {
	res, err := os.pdSrv.GetPd(ctx, &v1.GetPdReq{Id: id})
	if err != nil {
		return nil, err
	}
	return res.Pd, nil

}

func (os *OdSrv) GetSp(ctx context.Context, id int64) (*v1.Shop, error) {
	res, err := os.pdSrv.GetShop(ctx, &v1.GetShopReq{Id: id})
	if err != nil {
		return nil, err

	}
	return res.Sp, nil
}
