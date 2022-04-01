package service

import (
	"context"
	"fmt"
	v1 "pd-srv/api/v1"
)

func (ps *PdService) CreateCg(ctx context.Context, req *v1.CreateCgReq) (*v1.CreateCgReply, error) {
	if len(req.Cg.Name) == 0 {
		return &v1.CreateCgReply{Cg: &v1.Category{}}, fmt.Errorf("name is nil")
	}
	res, err := ps.cgCase.CreateCg(ctx, req)
	if err != nil {
		ps.logger.Errorf("create category error:%s", err.Error())
		return &v1.CreateCgReply{Cg: &v1.Category{}}, fmt.Errorf("failed to create category,internal server error")
	}
	return res, nil
}

func (ps *PdService) UpdateCg(ctx context.Context, req *v1.UpdateCgReq) (*v1.UpdateCgReply, error) {
	if req.Cg.Id < 1 {
		return &v1.UpdateCgReply{Cg: &v1.Category{}}, fmt.Errorf("id is wrong")
	}
	if len(req.Cg.Name) == 0 {
		return &v1.UpdateCgReply{Cg: &v1.Category{}}, fmt.Errorf("name is nil")
	}
	res, err := ps.cgCase.UpdateCg(ctx, req)
	if err != nil {
		ps.logger.Errorf("update category error:%s", err.Error())
		return &v1.UpdateCgReply{Cg: &v1.Category{}}, fmt.Errorf("failed to update category,internal server error")
	}
	return res, nil
}

func (ps *PdService) DeleteCg(ctx context.Context, req *v1.DeleteCgReq) (*v1.DeleteCgReply, error) {
	if req.Id < 1 {
		return &v1.DeleteCgReply{}, fmt.Errorf("category id is wrong")
	}
	res, err := ps.cgCase.DeleteCg(ctx, req)
	if err != nil {
		ps.logger.Errorf("delete category error:%s", err.Error())
		return nil, err
	}
	return res, nil

}

func (ps *PdService) GetCg(ctx context.Context, req *v1.GetCgReq) (*v1.GetCgReply, error) {
	if req.Id < 1 {
		return &v1.GetCgReply{Cg: &v1.Category{}}, fmt.Errorf("category id is wrong")
	}
	res, err := ps.cgCase.GetCg(ctx, req)
	if err != nil {
		ps.logger.Errorf("get category error:%s", err.Error())
		return &v1.GetCgReply{Cg: &v1.Category{}}, fmt.Errorf("failed to get category ,internal server error")
	}
	return res, nil
}

func (ps *PdService) ListCg(ctx context.Context, req *v1.ListCgReq) (*v1.ListCgReply, error) {
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}
	res, err := ps.cgCase.ListCg(ctx, req)
	if err != nil {
		ps.logger.Errorf("list category error:%s", err.Error())
		return &v1.ListCgReply{CgList: make([]*v1.Category, 0)}, fmt.Errorf("failed to list category ,internal server error")
	}
	return res, nil
}
