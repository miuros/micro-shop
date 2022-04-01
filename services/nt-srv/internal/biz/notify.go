package biz

import (
	"context"
	v1 "nt-srv/api/notify/v1"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type NtUseCase struct {
	ntRepo NtRepo
	logger *log.Helper
}

type QueryNotice struct {
	Offset   uint64
	Limit    uint64
	Status   uint64
	Type     string
	UserUuid string
}

type Notice struct {
	Id         uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserUuid   string `json:"userUuid" gorm:"column:user_uuid;not null"`
	UserName   string `json:"userName" gorm:"column:user_name;not null"`
	ToUserUuid string `json:"toUserUuid" gorm:"column:to_user_uuid;not null"`
	Content    string `json:"content" gorm:"column:content;not null"`
	Type       string `json:"type" gorm:"column:type;not null"`
	CreateAt   string `json:"createAt" gorm:"create_at;not null"`
	Status     uint64 `json:"status" gorm:"column:status;default 0"`
	IsDeleted  uint64 `json:"isDeleted" gorm:"column:is_deleted;default 0"`
}

type NtRepo interface {
	CreateNt(context.Context, *Notice) (*Notice, error)
	UpdateStatus(context.Context, *Notice) error
	DeleteNt(context.Context, *Notice) (*Notice, error)
	GetNt(context.Context, *Notice) (*Notice, error)
	ListNt(context.Context, *QueryNotice) ([]Notice, int, error)
}

func NewNtUseCase(ntRepo NtRepo, logger log.Logger) *NtUseCase {
	return &NtUseCase{
		ntRepo: ntRepo,
		logger: log.NewHelper(log.With(logger, "module", "biz/notify")),
	}
}

func (nu *NtUseCase) CreateNt(ctx context.Context, req *v1.Notice) (*v1.Notice, error) {
	var nt = &Notice{
		Id:         0,
		UserUuid:   req.UserUuid,
		UserName:   req.UserName,
		ToUserUuid: req.ToUserUuid,
		Content:    req.Content,
		Type:       req.Type,
		Status:     1,
		CreateAt:   time.Now().Format("2006-01-02:15-04"),
		IsDeleted:  0,
	}
	res, err := nu.ntRepo.CreateNt(ctx, nt)
	if err != nil {
		return nil, err
	}
	return replaceNt(res), nil
}

func (nu *NtUseCase) GetNt(ctx context.Context, req *v1.GetNtReq) (*v1.Notice, error) {
	var nt = &Notice{
		Id:       uint64(req.Id),
		UserUuid: req.UserUuid,
		Type:     req.Type,
	}
	res, err := nu.ntRepo.GetNt(ctx, nt)
	if err != nil {
		return nil, err
	}
	return replaceNt(res), nil
}

func replaceNt(nt *Notice) *v1.Notice {
	return &v1.Notice{
		Id:         int64(nt.Id),
		ToUserUuid: nt.ToUserUuid,
		UserName:   nt.UserName,
		UserUuid:   nt.UserUuid,
		Content:    nt.Content,
		Status:     int64(nt.Status),
		Type:       nt.Type,
		CreateAt:   nt.CreateAt,
		IsDeleted:  int64(nt.IsDeleted),
	}
}

func (nu *NtUseCase) UpdateStatus(ctx context.Context, req *v1.UpdateStatusReq) error {
	var nt = &Notice{
		Id:       uint64(req.Id),
		Type:     req.Type,
		UserUuid: req.UserUuid,
		Status:   0,
	}
	return nu.ntRepo.UpdateStatus(ctx, nt)
}

func (nu *NtUseCase) DeleteNt(ctx context.Context, req *v1.DeleteNtReq) (*v1.DeleteNtReply, error) {
	if req.Type == "chat" {
		return &v1.DeleteNtReply{}, nil
	}
	var nt = &Notice{
		Id:       uint64(req.Id),
		UserUuid: req.UserUuid,
		Type:     req.Type,
	}
	_, err := nu.ntRepo.DeleteNt(ctx, nt)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteNtReply{}, nil
}

func (nu *NtUseCase) ListNt(ctx context.Context, req *v1.ListNtReq) (*v1.ListNtReply, error) {
	var qn = &QueryNotice{
		Type:     req.Type,
		Offset:   uint64(req.Limit * (req.Page - 1)),
		Limit:    uint64(req.Limit),
		Status:   uint64(req.Status),
		UserUuid: req.UserUuid,
	}

	res, num, err := nu.ntRepo.ListNt(ctx, qn)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListNtReply{NtList: make([]*v1.Notice, 0), Num: 0}, nil
	}
	var ntList = make([]*v1.Notice, len(res))
	for i, item := range res {
		ntList[i] = replaceNt(&item)
	}
	return &v1.ListNtReply{NtList: ntList, Num: int64(num)}, nil
}
