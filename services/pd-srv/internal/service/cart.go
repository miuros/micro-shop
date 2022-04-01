package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	v1 "pd-srv/api/v1"
)

func (ps *PdService) CreateCart(ctx context.Context, req *v1.CreateCartReq) (*v1.CreateCartReply, error) {
	if len(req.C.UserUuid) == 0 {

	}
	res, err := ps.cartCase.CreateCart(ctx, req)
	if err != nil {
		ps.logger.Errorf("create cart error")
		return &v1.CreateCartReply{C: &v1.Cart{}}, fmt.Errorf("failed to create cart,internal server error")
	}
	return res, nil
}

func (ps *PdService) UpdateCart(ctx context.Context, req *v1.UpdateCartReq) (*v1.UpdateCartReply, error) {
	if req.C.Num < 1 {
		return &v1.UpdateCartReply{C: &v1.Cart{}}, nil
	}
	_, err := ps.cartCase.GetCart(ctx, &v1.GetCartReq{Id: req.C.Id, UserUuid: req.C.UserUuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.UpdateCartReply{C: &v1.Cart{}}, fmt.Errorf("no such a cart")
		}
		ps.logger.Errorf("get cart error:%s", err.Error())
	}
	res, err := ps.cartCase.UpdateCart(ctx, req)
	if err != nil {
		ps.logger.Errorf("update cart error:%s", err.Error())
		return &v1.UpdateCartReply{C: &v1.Cart{}}, fmt.Errorf("failed to update cart,internal server error")
	}
	return res, nil
}

func (ps *PdService) DeleteCart(ctx context.Context, req *v1.DeleteCartReq) (*v1.DeleteCartReply, error) {
	if req.Id < 1 || len(req.UserUuid) < 20 {
		return &v1.DeleteCartReply{}, fmt.Errorf("wrong id or uuid")
	}
	_, err := ps.cartCase.GetCart(ctx, &v1.GetCartReq{Id: req.Id, UserUuid: req.UserUuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeleteCartReply{}, fmt.Errorf("no such a cart")
		}
		ps.logger.Errorf("get cart error:%s", err.Error())
		return &v1.DeleteCartReply{}, fmt.Errorf("failed to delete cart,internal server error")
	}
	res, err := ps.cartCase.DeleteCart(ctx, req)
	if err != nil {
		ps.logger.Errorf("delete cart error:%s", err.Error())
		return &v1.DeleteCartReply{}, fmt.Errorf("failed to delete cart ,internal server error")
	}
	return res, nil
}

func (ps *PdService) GetCart(ctx context.Context, req *v1.GetCartReq) (*v1.GetCartReply, error) {
	if req.Id < 1 || len(req.UserUuid) < 20 {
		return &v1.GetCartReply{C: &v1.Cart{}}, fmt.Errorf("wrong with id or uuid")
	}
	res, err := ps.cartCase.GetCart(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetCartReply{C: &v1.Cart{}}, fmt.Errorf("no such a cart")
		}
		ps.logger.Errorf("get cart error:%s", err.Error())
		return &v1.GetCartReply{C: &v1.Cart{}}, fmt.Errorf("failed get cart ,internal server error")
	}
	return res, nil
}

func (ps *PdService) ListCart(ctx context.Context, req *v1.ListCartReq) (*v1.ListCartReply, error) {
	if req.Limit < 1 {
		req.Limit = 1
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if len(req.UserUuid) < 20 {
		return &v1.ListCartReply{CartList: make([]*v1.Cart, 0)}, fmt.Errorf("wrong with user uuid")
	}
	caList, err := ps.cartCase.ListCart(ctx, req)
	if err != nil {
		ps.logger.Errorf("list cart error:%s", err.Error())
		return &v1.ListCartReply{CartList: make([]*v1.Cart, 0)}, fmt.Errorf("failed to get cart list ,internal server error")
	}
	return caList, nil
}
