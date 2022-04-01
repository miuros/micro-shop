package service

import (
	"context"
	"fmt"
	v1 "nt-srv/api/notify/v1"

	"github.com/jinzhu/gorm"
)

func (ns *NtSrv) CreateNt(ctx context.Context, req *v1.CreateNtReq) (*v1.CreateNtReply, error) {
	if len(req.N.Type) == 0 {
		return &v1.CreateNtReply{}, fmt.Errorf("type is nil")
	}
	if req.N.Type != "notice" && req.N.Type != "chat" {
		return &v1.CreateNtReply{}, fmt.Errorf("type is wrong")
	}
	if len(req.N.UserUuid) == 0 {
		return &v1.CreateNtReply{}, fmt.Errorf("user uuid is nil")
	}
	if len(req.N.UserName) == 0 {
		return &v1.CreateNtReply{}, fmt.Errorf("name is nil")
	}
	if len(req.N.ToUserUuid) == 0 {
		return &v1.CreateNtReply{}, fmt.Errorf("to user uuid is nil")
	}
	if len(req.N.Content) == 0 {
		return &v1.CreateNtReply{}, fmt.Errorf("content is nil")
	}
	_, err := ns.ntCase.CreateNt(ctx, req.N)
	if err != nil {
		ns.logger.Errorf("create notify error:%s", err.Error())
		return &v1.CreateNtReply{}, fmt.Errorf("failed to create notice,internal server error")
	}
	return &v1.CreateNtReply{}, nil
}
func (ns *NtSrv) GetNt(ctx context.Context, req *v1.GetNtReq) (*v1.GetNtReply, error) {
	if req.Id < 1 {
		return &v1.GetNtReply{}, fmt.Errorf("id is nil")
	}
	if len(req.UserUuid) == 0 {
		return &v1.GetNtReply{}, fmt.Errorf("user uuid is nil")
	}
	if len(req.Type) == 0 {
		return &v1.GetNtReply{}, fmt.Errorf("type is nil")
	}
	if req.Type != "chat" && req.Type != "notice" {
		return &v1.GetNtReply{}, fmt.Errorf("type iss wrong")
	}
	res, err := ns.ntCase.GetNt(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetNtReply{}, fmt.Errorf("wrong id")
		}
		ns.logger.Errorf("get notice error:%s", err.Error())
		return &v1.GetNtReply{}, fmt.Errorf("failed to get notice")
	}
	return &v1.GetNtReply{Nt: res}, nil

}

func (ns *NtSrv) UpdateStatus(ctx context.Context, req *v1.UpdateStatusReq) (*v1.UpdateStatusReply, error) {
	if req.Id < 1 {
		return &v1.UpdateStatusReply{}, fmt.Errorf("id is nil")
	}
	if len(req.UserUuid) == 0 {
		return &v1.UpdateStatusReply{}, fmt.Errorf("user uuid is nil")
	}
	if len(req.Type) == 0 {
		return &v1.UpdateStatusReply{}, fmt.Errorf("type is nil")
	}
	if req.Type != "chat" && req.Type != "notice" {
		return &v1.UpdateStatusReply{}, fmt.Errorf("type iss wrong")
	}
	err := ns.ntCase.UpdateStatus(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.UpdateStatusReply{}, fmt.Errorf("wrong id")
		}
		ns.logger.Errorf("update status error:%s", err.Error())
		return &v1.UpdateStatusReply{}, fmt.Errorf("failed to update status")
	}
	return &v1.UpdateStatusReply{}, nil

}

func (ns *NtSrv) ListNt(ctx context.Context, req *v1.ListNtReq) (*v1.ListNtReply, error) {
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if len(req.UserUuid) == 0 {
		return &v1.ListNtReply{NtList: make([]*v1.Notice, 0), Num: 0}, fmt.Errorf("user uuid is nil")
	}
	if req.Type != "chat" && req.Type != "notice" {
		return &v1.ListNtReply{NtList: make([]*v1.Notice, 0), Num: 0}, fmt.Errorf("type is wrong")
	}
	if req.Status != 2 && req.Status != 1 {
		req.Status = 0
	}
	res, err := ns.ntCase.ListNt(ctx, req)
	if err != nil {
		ns.logger.Errorf("list notify error:%s", err.Error())
		return &v1.ListNtReply{NtList: make([]*v1.Notice, 0), Num: 0}, fmt.Errorf("failed to get notice,internal server error")
	}
	return res, nil
}

func (ns *NtSrv) DeleteNt(ctx context.Context, req *v1.DeleteNtReq) (*v1.DeleteNtReply, error) {
	if req.Id < 1 {
		return &v1.DeleteNtReply{}, fmt.Errorf("id error")

	}
	if len(req.UserUuid) == 0 {
		return &v1.DeleteNtReply{}, fmt.Errorf("user id is nil")
	}
	if req.Type != "chat" && req.Type != "notice" {
		return &v1.DeleteNtReply{}, fmt.Errorf("type is error")
	}
	res, err := ns.ntCase.DeleteNt(ctx, req)
	if err != nil {
		return &v1.DeleteNtReply{}, err
	}
	return res, nil
}
