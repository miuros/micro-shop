package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	v1 "pd-srv/api/v1"
)

func (ps *PdService) CreateBn(ctx context.Context, req *v1.CreateBnReq) (*v1.CreateBnReply, error) {
	if len(req.Bn.ImageUrl) == 0 || len(req.Bn.RedirectUrl) == 0 {
		return &v1.CreateBnReply{Bn: &v1.Banner{}}, fmt.Errorf("wrong with image url or redirect url")
	}
	res, err := ps.bnCase.CreateBn(ctx, req)
	if err != nil {
		ps.logger.Errorf("create banner error:%s", err.Error())
		return &v1.CreateBnReply{Bn: &v1.Banner{}}, fmt.Errorf("failed to create banner ,internal server error")
	}
	return res, nil

}

func (ps *PdService) UpdateBn(ctx context.Context, req *v1.UpdateBnReq) (*v1.UpdateBnReply, error) {
	if req.Bn.Id < 1 {
		return &v1.UpdateBnReply{Bn: &v1.Banner{}}, fmt.Errorf("wrong with id")
	}
	res, err := ps.bnCase.UpdateBn(ctx, req)
	if err != nil {
		ps.logger.Errorf("update banner error:%s", err.Error())
		return &v1.UpdateBnReply{Bn: &v1.Banner{}}, fmt.Errorf("failed to update banner ,internal server error")
	}
	return res, nil
}

func (ps *PdService) DeleteBn(ctx context.Context, req *v1.DeleteBnReq) (*v1.DeleteBnReply, error) {
	if req.Id < 1 {
		return &v1.DeleteBnReply{}, fmt.Errorf("wrong with id")
	}
	res, err := ps.bnCase.DeleteBn(ctx, req)
	if err != nil {
		ps.logger.Errorf("delete banner error:%s", err.Error())
		return &v1.DeleteBnReply{}, fmt.Errorf("failed to delete banner ,internal server error")
	}
	return res, nil
}

func (ps *PdService) GetBn(ctx context.Context, req *v1.GetBnReq) (*v1.GetBnReply, error) {
	if req.Id < 1 {
		return &v1.GetBnReply{Bn: &v1.Banner{}}, fmt.Errorf("wrong with id")
	}
	res, err := ps.bnCase.GetBn(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetBnReply{Bn: &v1.Banner{}}, fmt.Errorf("no such a banner")
		}
		ps.logger.Errorf("get banner error:%s", err.Error())
		return &v1.GetBnReply{Bn: &v1.Banner{}}, fmt.Errorf("failed to get banner ,internal server error")
	}
	return res, nil
}

func (ps *PdService) ListBn(ctx context.Context, req *v1.ListBnReq) (*v1.ListBnReply, error) {
	if req.Limit < 1 {
		req.Limit = 1
	}
	if req.Page < 1 {
		req.Page = 1
	}
	bnList, err := ps.bnCase.ListBn(ctx, req)
	if err != nil {
		ps.logger.Errorf("list banner error:%s", err.Error())
		return &v1.ListBnReply{BnList: make([]*v1.Banner, 0)}, fmt.Errorf("failed to list banner ,internal server error")
	}
	return bnList, nil
}
