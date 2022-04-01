package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	v1 "pd-srv/api/v1"
)

func (ps *PdService) CreateShop(ctx context.Context, req *v1.CreateShopReq) (*v1.CreateShopReply, error) {
	if len(req.Sp.Name) == 0 {

	}
	if len(req.Sp.UserUuid) == 0 {

	}
	if len(req.Sp.Address) == 0 {

	}
	res, err := ps.shopCase.CreateShop(ctx, req)
	if err != nil {
		ps.logger.Errorf("create shop error:%s", err.Error())
		return &v1.CreateShopReply{Sp: &v1.Shop{}}, fmt.Errorf("failed to create shop,internal server error")
	}
	return res, nil
}

func (ps *PdService) UpdateShop(ctx context.Context, req *v1.UpdateShopReq) (*v1.UpdateShopReply, error) {
	trueSp, err := ps.shopCase.GetShop(ctx, &v1.GetShopReq{Id: req.Sp.Id})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.UpdateShopReply{Sp: &v1.Shop{}}, fmt.Errorf("no such a shop")
		}
		ps.logger.Errorf("get shop error:%s", err.Error())
		return &v1.UpdateShopReply{Sp: &v1.Shop{}}, fmt.Errorf("failed to get shop ,internal server error")
	}
	if trueSp.Sp.UserUuid != req.Sp.UserUuid {
		return &v1.UpdateShopReply{Sp: &v1.Shop{}}, fmt.Errorf("user uuid is wrong")
	}
	res, err := ps.shopCase.UpdateShop(ctx, req)
	if err != nil {
		ps.logger.Errorf("update shop error:%s", err.Error())
		return &v1.UpdateShopReply{Sp: &v1.Shop{}}, fmt.Errorf("failed to update shop,internal server error")
	}
	return res, nil
}

func (ps *PdService) DeleteShop(ctx context.Context, req *v1.DeleteShopReq) (*v1.DeleteShopReply, error) {
	trueSp, err := ps.shopCase.GetShop(ctx, &v1.GetShopReq{Id: req.Id})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeleteShopReply{}, fmt.Errorf("no such a shop")
		}
		ps.logger.Errorf("get shop error:%s", err.Error())
		return &v1.DeleteShopReply{}, fmt.Errorf("failed to get shop,internal server error")
	}
	if trueSp.Sp.UserUuid != req.UserUuid {
		return &v1.DeleteShopReply{}, fmt.Errorf("user uuid is wrong")
	}
	res, err := ps.shopCase.DeleteShop(ctx, int(req.Id))
	if err != nil {
		ps.logger.Errorf("delete shop error:%s", err.Error())
		return &v1.DeleteShopReply{}, fmt.Errorf("failed to delete shop")
	}
	return res, nil

}

func (ps *PdService) GetShop(ctx context.Context, req *v1.GetShopReq) (*v1.GetShopReply, error) {
	res, err := ps.shopCase.GetShop(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetShopReply{Sp: &v1.Shop{}}, fmt.Errorf("no such a shop")
		}
		ps.logger.Errorf("get shop error:%s", err.Error())
		return &v1.GetShopReply{Sp: &v1.Shop{}}, fmt.Errorf("failed to get shop,internal server error")
	}
	return res, nil
}

func (ps *PdService) ListShop(ctx context.Context, req *v1.ListShopReq) (*v1.ListShopReply, error) {
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}
	spList, err := ps.shopCase.ListShop(ctx, req)
	if err != nil {
		ps.logger.Errorf("list shop error:%s", err.Error())
		return &v1.ListShopReply{SpList: make([]*v1.Shop, 0)}, fmt.Errorf("failed get shop list,internal server error")
	}
	return spList, nil
}

func (ps *PdService) GetSpByUuid(ctx context.Context, req *v1.GetSpByUuidReq) (*v1.GetSpByUuidReply, error) {
	if len(req.UserUuid) == 0 {
		return &v1.GetSpByUuidReply{Sp: &v1.Shop{}}, fmt.Errorf("user uuid is nil")
	}
	res, err := ps.shopCase.GetSpByUuid(ctx, req)
	if err != nil {
		ps.logger.Errorf("get shop by user uuid error:%s", err.Error())
		return &v1.GetSpByUuidReply{Sp: &v1.Shop{}}, fmt.Errorf("failed to get shop,internal server error")

	}
	return res, nil
}
