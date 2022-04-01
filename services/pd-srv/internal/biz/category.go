package biz

import (
	"context"
	v1 "pd-srv/api/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type Category struct {
	Id   uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name"`
}

type QueryCg struct {
	Offset uint64
	Limit  uint64
}

type CgRepo interface {
	CreateCg(context.Context, *Category) (*Category, error)
	UpdateCg(context.Context, *Category) (*Category, error)
	DeleteCg(context.Context, *Category) (*Category, error)
	GetCg(context.Context, *Category) (*Category, error)
	ListCg(context.Context, *QueryCg) ([]Category, error)
}

type CgUseCase struct {
	cr     CgRepo
	logger *log.Helper
}

func NewCgUseCase(logger log.Logger, cr CgRepo) *CgUseCase {
	return &CgUseCase{
		cr:     cr,
		logger: log.NewHelper(log.With(logger, "biz/category")),
	}
}

func (cu *CgUseCase) CreateCg(ctx context.Context, req *v1.CreateCgReq) (*v1.CreateCgReply, error) {
	var cg = &Category{
		Id:   0,
		Name: req.Cg.Name,
	}
	res, err := cu.cr.CreateCg(ctx, cg)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCgReply{Cg: transformCg(res)}, nil
}

func (cu *CgUseCase) UpdateCg(ctx context.Context, req *v1.UpdateCgReq) (*v1.UpdateCgReply, error) {
	var cg = &Category{
		Id:   uint64(req.Cg.Id),
		Name: req.Cg.Name,
	}
	res, err := cu.cr.UpdateCg(ctx, cg)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateCgReply{Cg: transformCg(res)}, nil
}

func (cu *CgUseCase) DeleteCg(ctx context.Context, req *v1.DeleteCgReq) (*v1.DeleteCgReply, error) {
	var cg = &Category{
		Id: uint64(req.Id),
	}
	_, err := cu.cr.DeleteCg(ctx, cg)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteCgReply{}, nil
}

func (cu *CgUseCase) GetCg(ctx context.Context, req *v1.GetCgReq) (*v1.GetCgReply, error) {
	var cg = &Category{Id: uint64(req.Id)}
	res, err := cu.cr.GetCg(ctx, cg)
	if err != nil {
		return nil, err

	}
	return &v1.GetCgReply{Cg: transformCg(res)}, nil
}

func (cu *CgUseCase) ListCg(ctx context.Context, req *v1.ListCgReq) (*v1.ListCgReply, error) {
	var qc = &QueryCg{
		Limit:  uint64(req.Limit),
		Offset: uint64(req.Limit * (req.Page - 1)),
	}
	res, err := cu.cr.ListCg(ctx, qc)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListCgReply{CgList: make([]*v1.Category, 0)}, nil
	}
	var cgList = make([]*v1.Category, len(res))
	for i, item := range res {
		cgList[i] = transformCg(&item)
	}
	return &v1.ListCgReply{CgList: cgList}, nil
}

func transformCg(res *Category) *v1.Category {
	return &v1.Category{
		Id:   uint64(res.Id),
		Name: res.Name,
	}
}
