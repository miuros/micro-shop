package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	v1 "pd-srv/api/v1"
)

func (ps *PdService) verifyShopByPdId(ctx context.Context, pid int64, sid int64, uuid string) error {
	truePd, err := ps.pdCase.GetPd(ctx, &v1.GetPdReq{Id: int64(pid)})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no such a product")
		}
		ps.logger.Errorf("get product error:%s", err.Error())
		return fmt.Errorf("get product error")
	}
	if sid != truePd.Pd.ShopId {
		return fmt.Errorf("shop id is wrong")
	}
	trueSp, err := ps.shopCase.GetShop(ctx, &v1.GetShopReq{Id: int64(truePd.Pd.ShopId)})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no such a shop")
		}
		ps.logger.Errorf("failed to get shop info:%s", err.Error())
		return fmt.Errorf("failed to get shop info")
	}
	if trueSp.Sp.UserUuid != uuid {
		return fmt.Errorf("uuid is wrong")
	}
	return nil
}

func (ps *PdService) CreatePd(ctx context.Context, req *v1.CreatePdReq) (*v1.CreatePdReply, error) {
	trueSp, err := ps.shopCase.GetShop(ctx, &v1.GetShopReq{Id: int64(req.Pd.ShopId)})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.CreatePdReply{Pd: &v1.Product{}}, fmt.Errorf("no such a shop")
		}
		ps.logger.Errorf("failed to get shop info:%s", err.Error())
		return &v1.CreatePdReply{Pd: &v1.Product{}}, fmt.Errorf("failed to get shop info")
	}
	if trueSp.Sp.UserUuid != req.UserUuid {
		return &v1.CreatePdReply{Pd: &v1.Product{}}, fmt.Errorf("uuid is wrong")
	}

	res, err := ps.pdCase.CreatePd(ctx, req)
	if err != nil {
		ps.logger.Errorf("create product error:%s", err.Error())
		return &v1.CreatePdReply{Pd: &v1.Product{}}, fmt.Errorf("failed to create product")
	}
	return res, nil
}

func (ps *PdService) UpdatePd(ctx context.Context, req *v1.UpdatePdReq) (*v1.UpdatePdReply, error) {

	err := ps.verifyShopByPdId(ctx, req.Pd.Id, req.Pd.ShopId, req.UserUuid)
	if err != nil {
		return &v1.UpdatePdReply{Pd: &v1.Product{}}, err
	}
	res, err := ps.pdCase.UpdatePd(ctx, req)
	if err != nil {
		ps.logger.Errorf("update product error:%s", err.Error())
		return &v1.UpdatePdReply{Pd: &v1.Product{}}, fmt.Errorf("failed to update product")
	}
	return res, nil
}

func (ps *PdService) DeletePd(ctx context.Context, req *v1.DeletePdReq) (*v1.DeletePdReply, error) {
	truePd, err := ps.pdCase.GetPd(ctx, &v1.GetPdReq{Id: req.Id})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeletePdReply{}, fmt.Errorf("no such a prodcut")
		}
		ps.logger.Errorf("get pd error:%s", err.Error())
		return &v1.DeletePdReply{}, fmt.Errorf("failed to delete product ,can not find product")
	}
	trueSp, err := ps.shopCase.GetShop(ctx, &v1.GetShopReq{Id: int64(truePd.Pd.ShopId)})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeletePdReply{}, fmt.Errorf("no such a product")
		}
		ps.logger.Errorf("get shop error:%s", err.Error())
		return &v1.DeletePdReply{}, fmt.Errorf("failed to delete product,can not find shop")
	}
	if req.UserUuid != trueSp.Sp.UserUuid {
		return &v1.DeletePdReply{}, fmt.Errorf("uuid is wrong")
	}
	_, err = ps.pdCase.DeletePd(ctx, &v1.DeletePdReq{Id: req.Id, UserUuid: req.UserUuid})
	if err != nil {
		ps.logger.Errorf("delete product error:%s", err.Error())
		return &v1.DeletePdReply{}, fmt.Errorf("failed to delete product,internal server error")
	}
	return &v1.DeletePdReply{}, nil

}

func (ps *PdService) GetPd(ctx context.Context, req *v1.GetPdReq) (*v1.GetPdReply, error) {
	res, err := ps.pdCase.GetPd(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetPdReply{Pd: &v1.Product{}}, fmt.Errorf("no such a product")
		}
		ps.logger.Errorf("failed to get product,internal server error")
	}
	return res, nil
}

func (ps *PdService) ListPd(ctx context.Context, req *v1.ListPdReq) (*v1.ListPdReply, error) {
	if req.Limit < 0 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}
	pdList, err := ps.pdCase.ListPd(ctx, req)
	if err != nil {
		ps.logger.Errorf("list product error:%s", err.Error())
		return &v1.ListPdReply{PdList: make([]*v1.Product, 0)}, fmt.Errorf("failed to get product")
	}
	return pdList, nil
}

func (ps *PdService) ListPdByCi(ctx context.Context, req *v1.ListPdByCiReq) (*v1.ListPdByCiReply, error) {
	if req.CategoryId < 1 {
		return &v1.ListPdByCiReply{PdList: make([]*v1.Product, 0)}, fmt.Errorf("request category id is nil")
	}
	if req.Limit < 1 {
		req.Limit = 5
	}
	if req.Page < 1 {
		req.Page = 1
	}
	res, err := ps.pdCase.ListPdByCgId(ctx, req)
	if err != nil {
		ps.logger.Errorf("list product by category id error:%s", err.Error())
		return &v1.ListPdByCiReply{PdList: make([]*v1.Product, 0)}, fmt.Errorf("failed to list product,internal server error")
	}
	return res, nil
}

func (ps *PdService) FindPdByName(ctx context.Context, req *v1.ListPdReq) (*v1.ListPdReply, error) {
	return ps.ListPd(ctx, req)
}
