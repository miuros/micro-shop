package biz

import (
	"context"
	v1 "od-srv/api/order/v1"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type OdRepo interface {
	CreateOd(context.Context, *Order) (*Order, error)
	UpdateOd(context.Context, *Order) (*Order, error)
	DeleteOd(context.Context, *Order) (*Order, error)
	GetOd(context.Context, *Order) (*Order, error)
	ListOd(context.Context, *QueryOd) ([]Order, error)
	ListOdByCateId(context.Context, *Cate) ([]Order, error)
	ListOdForShopper(context.Context, *QueryOdForShopper) ([]Order, error)
	PayOd(context.Context, *Order) error
}

type QueryOd struct {
	Offset   uint   `json:"offset"`
	Limit    uint   `json:"limit"`
	UserUuid string `json:"userUuid"`
	Status   int64  `json:"status"`
}

type QueryOdForShopper struct {
	Offset uint64
	Limit  uint64
	ShopId uint64
	Status uint64
}

type Order struct {
	Id        int64   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ProductId int64   `json:"productId" gorm:"column:product_id;not null"`
	Number    int64   `json:"number" gorm:"column:number;not null"`
	PayType   int64   `json:"payType" gorm:"column:pay_type;not null"`
	Status    int64   `json:"status" gorm:"column:status;not null"`
	CateId    int64   `json:"cateId" gorm:"column:cate_id;not null"`
	AddressId int64   `json:"addressId" gorm:"column:address_id;not null"`
	IsDeleted int64   `json:"isDeleted" gorm:"column:is_deleted;default 0"`
	Price     float64 `json:"price" gorm:"column:price;not null"`
	UserUuid  string  `json:"userUuid" gorm:"column:user_uuid;not null"`
	PayTime   string  `json:"payTime" gorm:"column:pay_time;omitempty"`
	CreateAt  string  `json:"createAt" gorm:"column:create_at;omitempty"`
	UpdateAt  string  `json:"updateAt" gorm:"column:update_at;omitempty"`
	DeleteAt  string  `json:"deleteAt" gorm:"column:delete_at;omitempty"`
}

type OdUseCase struct {
	repo   OdRepo
	logger *log.Helper
}

func NewOdUseCate(repo OdRepo, logger log.Logger) *OdUseCase {
	return &OdUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/order")),
	}
}

func (ou *OdUseCase) CreateOd(ctx context.Context, req *v1.Item) (*v1.Item, error) {

	od := &Order{
		Id:        0,
		ProductId: req.ProductId,
		Number:    req.Number,
		Price:     float64(req.Price),
		Status:    1,
		AddressId: req.AddressId,
		UserUuid:  req.UserUuid,
		PayType:   0,
		CateId:    req.CateId,
		IsDeleted: 0,
		CreateAt:  time.Now().Format("2006-01-02:15-04"),
	}
	res, err := ou.repo.CreateOd(ctx, od)
	if err != nil {
		return nil, err
	}
	return createOd(res), nil
}

func createOd(od *Order) *v1.Item {
	return &v1.Item{
		Id:        od.Id,
		UserUuid:  od.UserUuid,
		Number:    od.Number,
		Price:     float32(od.Price),
		ProductId: od.ProductId,
		PayType:   od.PayType,
		Status:    od.Status,
		AddressId: od.AddressId,
		IsDeleted: od.IsDeleted,
		CateId:    od.CateId,
		CreateAt:  od.CreateAt,
		UpdateAt:  od.UpdateAt,
		DeletedAt: od.DeleteAt,
	}
}

func (ou *OdUseCase) UpdateOd(ctx context.Context, req *v1.UpdateOrderReq) (*v1.Item, error) {
	var od = &Order{
		Id:        int64(req.Id),
		AddressId: req.AddressId,
		UserUuid:  req.UserUuid,
	}
	res, err := ou.repo.UpdateOd(ctx, od)
	if err != nil {
		return nil, err
	}
	return createOd(res), nil
}

func (ou *OdUseCase) DeleteOd(ctx context.Context, req *v1.Item) (*v1.Item, error) {
	var od = &Order{
		Id:        req.Id,
		UserUuid:  req.UserUuid,
		ProductId: req.ProductId,
		Number:    req.Number,
		Status:    3,
		DeleteAt:  time.Now().Format("2006-01-02:15-04"),
		IsDeleted: 1,
	}
	res, err := ou.repo.DeleteOd(ctx, od)
	if err != nil {
		return nil, err

	}
	return createOd(res), nil
}

func (ou *OdUseCase) GetOd(ctx context.Context, req *v1.Item) (*v1.Item, error) {
	var od = &Order{
		Id:       req.Id,
		UserUuid: req.UserUuid,
	}
	res, err := ou.repo.GetOd(ctx, od)
	if err != nil {
		return nil, err
	}
	return createOd(res), nil
}

func (ou *OdUseCase) ListOd(ctx context.Context, req *v1.ListOrderReq) (*v1.ListOrderReply, error) {
	var qo = &QueryOd{
		Offset:   uint(req.Limit * (req.Page - 1)),
		Limit:    uint(req.Limit),
		UserUuid: req.UserUuid,
	}
	res, err := ou.repo.ListOd(ctx, qo)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListOrderReply{ItemList: make([]*v1.Item, 0)}, nil
	}
	var itemList = make([]*v1.Item, len(res))
	for i, item := range res {
		itemList[i] = createOd(&item)
	}
	return &v1.ListOrderReply{ItemList: itemList}, nil
}

func (ou *OdUseCase) ListOdByCateId(ctx context.Context, req *v1.ListOrderByCateIdReq) (*v1.ListOrderByCateIdReply, error) {
	var cate = &Cate{Id: req.CateId, UserUuid: req.UserUuid}
	res, err := ou.repo.ListOdByCateId(ctx, cate)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListOrderByCateIdReply{ItemList: make([]*v1.Item, 0), Cate: &v1.Cate{Id: cate.Id, UserUuid: cate.UserUuid}}, nil
	}
	var itemList = make([]*v1.Item, len(res))
	for i, item := range res {
		itemList[i] = createOd(&item)
	}
	return &v1.ListOrderByCateIdReply{ItemList: itemList, Cate: &v1.Cate{Id: cate.Id, UserUuid: cate.UserUuid}}, nil
}

func (ou *OdUseCase) ListOdForSp(ctx context.Context, req *v1.ListOdForSReq) (*v1.ListOdForSpReply, error) {
	var query = &QueryOdForShopper{
		ShopId: req.ShopId,
		Offset: req.Limit * (req.Page - 1),
		Limit:  req.Limit,
		Status: req.Status,
	}
	res, err := ou.repo.ListOdForShopper(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListOdForSpReply{OdList: make([]*v1.Item, 0)}, nil
	}
	var odList = make([]*v1.Item, len(res))
	for i, item := range res {
		odList[i] = createOd(&item)
	}
	return &v1.ListOdForSpReply{OdList: odList}, nil
}

func (ou *OdUseCase) PayOd(ctx context.Context, req *v1.Item) (*v1.PayOdReply, error) {
	var od = &Order{
		Id:        int64(req.Id),
		ProductId: int64(req.ProductId),
		Number:    int64(req.Number),
		UserUuid:  req.UserUuid,
		Status:    2,
	}
	err := ou.repo.PayOd(ctx, od)
	if err != nil {
		return nil, err
	}
	return &v1.PayOdReply{}, nil
}
