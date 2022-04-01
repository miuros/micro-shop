package service

import (
	v1 "comment-srv/api/comment/v1"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
)

func (cs *CmService) CreateCm(ctx context.Context, req *v1.CreateCmReq) (*v1.CreateCmReply, error) {
	if err := verifyNil(req.Cm); err != nil {
		return &v1.CreateCmReply{Cm: &v1.Comment{}}, err
	}
	res, err := cs.cmCase.CreateCm(ctx, req)
	if err != nil {
		cs.logger.Errorf("create comment error:%s", err.Error())
		return &v1.CreateCmReply{Cm: &v1.Comment{}}, fmt.Errorf("create comment error,internal server error")
	}
	return res, nil
}

func (cs *CmService) UpdateCm(ctx context.Context, req *v1.UpdateCmReq) (*v1.UpdateCmReply, error) {
	if err := verifyNil(req.Cm); err != nil {
		return &v1.UpdateCmReply{Cm: &v1.Comment{}}, err
	}
	trueCm, err := cs.cmCase.GetCm(ctx, &v1.GetCmReq{Id: req.Cm.Id})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.UpdateCmReply{Cm: &v1.Comment{}}, fmt.Errorf("no such a comment")
		}
		cs.logger.Errorf("get comment error:%s", err.Error())
		return &v1.UpdateCmReply{Cm: &v1.Comment{}}, fmt.Errorf("failed to update,internal server error")
	}
	if trueCm.Cm.UserUuid != req.Cm.UserUuid {
		return &v1.UpdateCmReply{Cm: &v1.Comment{}}, fmt.Errorf("wrong request,user uuid is not right")
	}
	if trueCm.Cm.ProductId != req.Cm.ProductId {
		return &v1.UpdateCmReply{Cm: &v1.Comment{}}, fmt.Errorf("wrong request,product id is not right")
	}
	res, err := cs.cmCase.UpdateCm(ctx, req)
	if err != nil {
		cs.logger.Errorf("update comment error:%s", err.Error())
		return &v1.UpdateCmReply{Cm: &v1.Comment{}}, fmt.Errorf("failed to update comment ,internal server error")
	}
	return res, nil
}

func (cs *CmService) DeleteCm(ctx context.Context, req *v1.DeleteCmReq) (*v1.DeleteCmReply, error) {
	if req.Id < 1 {
		return &v1.DeleteCmReply{}, fmt.Errorf("id is wrong")
	}
	trueCm, err := cs.cmCase.GetCm(ctx, &v1.GetCmReq{Id: req.Id})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeleteCmReply{}, fmt.Errorf("no such a comment")
		}
		cs.logger.Errorf("get comment error:%s", err.Error())
		return &v1.DeleteCmReply{}, fmt.Errorf("failed to delete comment,internal server error")
	}
	if trueCm.Cm.UserUuid != req.UserUuid {
		return &v1.DeleteCmReply{}, fmt.Errorf("user uuid is wrong")
	}
	res, err := cs.cmCase.DeleteCm(ctx, req)
	if err != nil {
		cs.logger.Errorf("delete comment error:%s", err.Error())
		return &v1.DeleteCmReply{}, fmt.Errorf("failed to delete comment,internal server error")
	}
	return res, nil
}

func (cs *CmService) GetCm(ctx context.Context, req *v1.GetCmReq) (*v1.GetCmReply, error) {
	if req.Id < 1 {
		return &v1.GetCmReply{Cm: &v1.Comment{}}, fmt.Errorf("id is wrong")
	}
	res, err := cs.cmCase.GetCm(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetCmReply{Cm: &v1.Comment{}}, fmt.Errorf("no such a comment")
		}
		cs.logger.Errorf("get comment error: %s", err.Error())
		return &v1.GetCmReply{Cm: &v1.Comment{}}, fmt.Errorf("failed to get comment ,internal server error")
	}
	return res, nil
}

func (cs *CmService) ListCm(ctx context.Context, req *v1.ListCmReq) (*v1.ListCmReply, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.ProductId < 1 {
		return &v1.ListCmReply{CmList: make([]*v1.Comment, 0)}, fmt.Errorf("product id is wrong")
	}

	cmList, err := cs.cmCase.ListCm(ctx, req)
	if err != nil {
		cs.logger.Errorf("list comment error: %s", err.Error())
		return &v1.ListCmReply{CmList: make([]*v1.Comment, 0)}, nil
	}
	return cmList, nil
}

func verifyNil(req *v1.Comment) error {

	if req.ProductId < 1 {
		return fmt.Errorf("product id is wrong")
	}
	if len(req.UserUuid) == 0 {

		return fmt.Errorf("user uuid is nil")
	}
	if len(req.Content) == 0 {

		return fmt.Errorf("content is nil")
	}
	return nil
}
