package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "od-srv/api/order/v1"
	"time"
)

type QueryCate struct {
	Offset   int64
	Limit    int64
	UserUuid string
}

type CateRepo interface {
	CreateCate(context.Context, *Cate) (*Cate, error)
	UpdateCate(context.Context, *Cate) (*Cate, error)
	DeleteCate(context.Context, *Cate) (*Cate, error)
	GetCate(context.Context, *Cate) (*Cate, error)
	ListCate(context.Context, *QueryCate) ([]Cate, error)
}

type Cate struct {
	Id        int64   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Status    int64   `json:"status"  gorm:"column:status;not null"`
	IsDeleted int64   `json:"isDeleted" gorm:"column:is_deleted;default 0"`
	AddressId int64   `json:"addressId" gorm:"column:address_id;not null"`
	Price     float64 `json:"price" gorm:"column:price;not null"`
	UserUuid  string  `json:"userUuid" gorm:"column:user_uuid;not null"`
	DeleteAt  string  `json:"deleteAt" gorm:"column:delete_at;omitempty"`
}

type CateUseCase struct {
	repo   CateRepo
	logger *log.Helper
}

func NewCateUseCase(repo CateRepo, logger log.Logger) *CateUseCase {
	return &CateUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/cate")),
	}
}

func (cu *CateUseCase) CreateCate(ctx context.Context, req *v1.Cate) (*v1.Cate, error) {
	var cate = &Cate{
		Id:        0,
		UserUuid:  req.UserUuid,
		Status:    1,
		IsDeleted: 0,
		AddressId: req.AddressId,
		Price:     float64(req.Price),
	}
	res, err := cu.repo.CreateCate(ctx, cate)
	if err != nil {
		return nil, err
	}
	return createCate(res), nil
}

func createCate(cate *Cate) *v1.Cate {
	return &v1.Cate{
		Id:        cate.Id,
		UserUuid:  cate.UserUuid,
		IsDeleted: cate.IsDeleted,
		Status:    cate.Status,
		Price:     float32(cate.Price),
		AddressId: cate.AddressId,
		DeleteAt:  cate.DeleteAt,
	}
}

func (cu *CateUseCase) UpdateCate(ctx context.Context, req *v1.UpdateCateReq) (*v1.UpdateCateReply, error) {
	var cate = &Cate{
		Id:        req.Cate.Id,
		UserUuid:  req.Cate.UserUuid,
		AddressId: req.Cate.AddressId,
	}
	res, err := cu.repo.UpdateCate(ctx, cate)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateCateReply{Cate: createCate(res)}, nil
}

func (cu *CateUseCase) DeleteCate(ctx context.Context, req *v1.DeleteCateReq) (*v1.DeleteCateReply, error) {
	var cate = &Cate{
		Id:        req.Id,
		UserUuid:  req.UserUuid,
		IsDeleted: 1,
		DeleteAt:  time.Now().String(),
	}
	_, err := cu.repo.DeleteCate(ctx, cate)
	if err != nil {
		return &v1.DeleteCateReply{}, err
	}
	return &v1.DeleteCateReply{}, nil
}

func (cu *CateUseCase) GetCate(ctx context.Context, req *v1.GetCateReq) (*v1.GetCateReply, error) {
	var cate = &Cate{Id: req.Id, UserUuid: req.UserUuid}
	res, err := cu.repo.GetCate(ctx, cate)
	if err != nil {
		return nil, err
	}
	return &v1.GetCateReply{Cate: createCate(res)}, nil
}

func (cu *CateUseCase) ListCate(ctx context.Context, req *v1.ListCateReq) (*v1.ListCateReply, error) {
	var qc = &QueryCate{
		UserUuid: req.UserUuid,
		Offset:   int64(req.Limit * (req.Page - 1)),
		Limit:    int64(req.Limit),
	}
	res, err := cu.repo.ListCate(ctx, qc)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListCateReply{CateList: make([]*v1.Cate, 0)}, nil
	}
	var caList = make([]*v1.Cate, 0, len(res))
	for i, item := range res {
		caList[i] = createCate(&item)
	}
	return &v1.ListCateReply{CateList: caList}, nil
}
