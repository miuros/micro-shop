package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "od-srv/api/order/v1"
)

type ScRepo interface {
	CreateSc(context.Context, *Stock) (*Stock, error)
	UpdateSc(context.Context, *Stock) (*Stock, error)
	DeleteSc(context.Context, *Stock) (*Stock, error)
	AddSc(context.Context, *Stock) (*Stock, error)
	GetSc(context.Context, *Stock) (*Stock, error)
}

type ScUseCase struct {
	repo   ScRepo
	logger *log.Helper
}

type Stock struct {
	Id        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ProductId int64 `json:"productId" gorm:"column:product_id;unique;not null"`
	Storage   int64 `json:"storage" gorm:"column:storage;not null"`
	Sale      int64 `json:"sale" gorm:"column:sale;not null;default 0"`
	IsDeleted int64 `json:"isDeleted" gorm:"column:is_deleted;default 0"`
}

func NewScUseCase(repo ScRepo, logger log.Logger) *ScUseCase {
	return &ScUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/stock")),
	}
}

func (su *ScUseCase) CreateSc(ctx context.Context, req *v1.CreateStockReq) (*v1.CreateStockReply, error) {
	var sc = &Stock{
		Id:        0,
		ProductId: req.Stock.ProductId,
		Storage:   req.Stock.Storage,
		Sale:      0,
		IsDeleted: 0,
	}
	res, err := su.repo.CreateSc(ctx, sc)
	if err != nil {
		return nil, err
	}
	return &v1.CreateStockReply{Stock: createSc(res)}, nil
}

func createSc(sc *Stock) *v1.StockInfo {
	return &v1.StockInfo{
		Id:        sc.Id,
		ProductId: sc.ProductId,
		Storage:   sc.Storage,
		Sale:      sc.Sale,
	}
}

func (su *ScUseCase) UpdateSc(ctx context.Context, req *v1.UpdateStockReq) (*v1.UpdateStockReply, error) {
	var sc = &Stock{Id: req.Stock.Id, ProductId: req.Stock.ProductId, Storage: req.Stock.Storage}
	res, err := su.repo.UpdateSc(ctx, sc)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateStockReply{Stock: createSc(res)}, nil
}

func (su *ScUseCase) DeleteSc(ctx context.Context, req *v1.DeleteStockReq) (*v1.DeleteStockReply, error) {
	var sc = &Stock{ProductId: req.ProductId}
	_, err := su.repo.DeleteSc(ctx, sc)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteStockReply{}, nil
}

func (su *ScUseCase) GetSc(ctx context.Context, req *v1.GetStockReq) (*v1.GetStockReply, error) {
	var sc = &Stock{ProductId: req.ProductId}
	res, err := su.repo.GetSc(ctx, sc)
	if err != nil {
		return nil, err
	}
	return &v1.GetStockReply{Stock: createSc(res)}, nil
}
