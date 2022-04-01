package biz

import (
	"context"
	"fmt"
	v1 "pd-srv/api/v1"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type PdUseCase struct {
	repo   ProductRepo
	logger *log.Helper
}

type Product struct {
	Id          int     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"column:name;not null"`
	OriginPrice float64 `json:"originPrice" gorm:"column:origin_price;not null"`
	ImageUrl    string  `json:"imageUrl" gorm:"column:image_url;omitempty"`
	SellPrice   float64 `json:"sellPrice" gorm:"column:sell_price;not null"`
	Desc        string  `json:"desc" gorm:"column:desc;omitempty"`
	Tags        string  `json:"tags" gorm:"column:tags;omitempty"`
	CategoryId  uint64  `json:"categoryId" gorm:"category_id;not null"`
	ShopId      int     `json:"shopId" gorm:"column:shop_id;not null"`
	Extra       string  `json:"extra" gorm:"column:extra;omitempty"`
	CreateAt    string  `json:"createAt" gorm:"column:create_at"`
	DeleteAt    string  `json:"deleteAt" gorm:"column:deleted_at"`
	IsDeleted   int     `json:"isDeleted" gorm:"column:is_deleted;default 0"`
}

type Query struct {
	Offset int
	Limit  int
	Name   string
}

type QueryForSp struct {
	Offset uint64
	Limit  uint64
	ShopId uint64
}
type QueryByCgId struct {
	Offset uint64
	Limit  uint64
	CgId   uint64
}

type ProductRepo interface {
	CreatePd(context.Context, *Product) (*Product, error)
	UpdatePd(context.Context, *Product) (*Product, error)
	DeletePd(context.Context, int) (*Product, error)
	GetPd(context.Context, int) (*Product, error)
	ListPd(context.Context, *Query) ([]Product, error)
	ListForSp(ctx context.Context, sp *QueryForSp) ([]Product, error)
	ListPdByCgId(context.Context, *QueryByCgId) ([]Product, error)
}

func NewPdUseCase(repo ProductRepo, logger log.Logger) *PdUseCase {
	return &PdUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz")),
	}
}

func (pr *PdUseCase) CreatePd(ctx context.Context, pdv1 *v1.CreatePdReq) (*v1.CreatePdReply, error) {
	if pdv1.Pd.Name == "" {
		return nil, fmt.Errorf("the name is not writed")
	}
	var pd = &Product{
		Id:          0,
		Name:        pdv1.Pd.Name,
		ImageUrl:    pdv1.Pd.ImageUrl,
		SellPrice:   float64(pdv1.Pd.SellPrice),
		OriginPrice: float64(pdv1.Pd.OriginPrice),
		Desc:        pdv1.Pd.Desc,
		Tags:        pdv1.Pd.Tags,
		CategoryId:  uint64(pdv1.Pd.CategoryId),
		ShopId:      int(pdv1.Pd.ShopId),
		Extra:       pdv1.Pd.Extra,
		CreateAt:    time.Now().Format("2006-01-02:15-04"),
		IsDeleted:   0,
	}
	res, err := pr.repo.CreatePd(ctx, pd)
	if err != nil {
		return nil, err
	}
	return &v1.CreatePdReply{
		Pd: createPdV1(res),
	}, nil
}

func (pu *PdUseCase) ListForSp(ctx context.Context, q *v1.ListForSpReq) (*v1.ListForSpReply, error) {
	var qs = &QueryForSp{
		Limit:  q.Limit,
		Offset: q.Limit * (q.Page - 1),
		ShopId: q.ShopId,
	}
	res, err := pu.repo.ListForSp(ctx, qs)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListForSpReply{PdList: make([]*v1.Product, 0)}, nil
	}
	var pdList = make([]*v1.Product, len(res))
	for i, item := range res {
		pdList[i] = createPdV1(&item)
	}
	return &v1.ListForSpReply{PdList: pdList}, nil
}

func (pr *PdUseCase) UpdatePd(ctx context.Context, pdv1 *v1.UpdatePdReq) (*v1.UpdatePdReply, error) {
	var pd = &Product{
		Id:          int(pdv1.Pd.Id),
		Name:        pdv1.Pd.Name,
		CategoryId:  uint64(pdv1.Pd.CategoryId),
		ImageUrl:    pdv1.Pd.ImageUrl,
		SellPrice:   float64(pdv1.Pd.SellPrice),
		OriginPrice: float64(pdv1.Pd.OriginPrice),
		Desc:        pdv1.Pd.Desc,
		Tags:        pdv1.Pd.Tags,
		Extra:       pdv1.Pd.Extra,
	}

	res, err := pr.repo.UpdatePd(ctx, pd)
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePdReply{
		Pd: createPdV1(res),
	}, nil
}

func (pr *PdUseCase) DeletePd(ctx context.Context, pdv1 *v1.DeletePdReq) (*v1.DeletePdReply, error) {
	_, err := pr.repo.DeletePd(ctx, int(pdv1.Id))
	if err != nil {
		return nil, err
	}
	return &v1.DeletePdReply{}, nil
}

func (pr *PdUseCase) GetPd(ctx context.Context, pdv1 *v1.GetPdReq) (*v1.GetPdReply, error) {
	res, err := pr.repo.GetPd(ctx, int(pdv1.GetId()))
	if err != nil {
		return nil, err

	}
	return &v1.GetPdReply{
		Pd: createPdV1(res),
	}, nil
}

func (pr *PdUseCase) ListPd(ctx context.Context, pdv1 *v1.ListPdReq) (*v1.ListPdReply, error) {
	var q = &Query{
		Limit:  int(pdv1.Limit),
		Offset: int(pdv1.Limit * (pdv1.Page - 1)),
		Name:   pdv1.Name,
	}
	res, err := pr.repo.ListPd(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListPdReply{PdList: make([]*v1.Product, 0)}, nil
	}
	var pdList = make([]*v1.Product, len(res))
	for i, item := range res {
		pdList[i] = createPdV1(&item)
	}
	return &v1.ListPdReply{
		PdList: pdList,
	}, nil
}

func (pr *PdUseCase) ListPdByCgId(ctx context.Context, req *v1.ListPdByCiReq) (*v1.ListPdByCiReply, error) {
	var q = &QueryByCgId{
		Limit:  uint64(req.Limit),
		Offset: uint64(req.Limit * (req.Page - 1)),
		CgId:   uint64(req.CategoryId),
	}
	res, err := pr.repo.ListPdByCgId(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListPdByCiReply{PdList: make([]*v1.Product, 0)}, nil
	}
	var pdList = make([]*v1.Product, len(res))
	for i, item := range res {
		pdList[i] = createPdV1(&item)
	}
	return &v1.ListPdByCiReply{PdList: pdList}, nil
}

func createPdV1(pd *Product) *v1.Product {
	return &v1.Product{
		Id:          int64(pd.Id),
		Name:        pd.Name,
		ImageUrl:    pd.ImageUrl,
		SellPrice:   float32(pd.SellPrice),
		OriginPrice: float32(pd.OriginPrice),
		Extra:       pd.Extra,
		Tags:        pd.Tags,
		CreateAt:    pd.CreateAt,
		CategoryId:  int64(pd.CategoryId),
		DeleteAt:    pd.DeleteAt,
		ShopId:      int64(pd.ShopId),
		IsDeleted:   int64(pd.IsDeleted),
	}
}
