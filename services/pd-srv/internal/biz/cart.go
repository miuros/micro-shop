package biz

import (
	"context"
	v1 "pd-srv/api/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type CartUseCase struct {
	repo   CartRepo
	logger *log.Helper
}

type CartRepo interface {
	CreateCart(context.Context, *Cart) (*Cart, error)
	UpdateCart(context.Context, *Cart) (*Cart, error)
	DeleteCart(context.Context, *Cart) (*Cart, error)
	GetCart(context.Context, *Cart) (*Cart, error)
	ListCart(context.Context, *QueryCart) ([]Cart, error)
}

type QueryCart struct {
	Offset   int
	Limit    int
	UserUuid string
}

type Cart struct {
	Id          int     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserUuid    string  `json:"userUuid" gorm:"column:user_uuid;not null"`
	ImageUrl    string  `json:"imageUrl" gorm:"column:image_url;omitempty"`
	ProductId   int     `json:"productId" gorm:"column:product_id;not null"`
	ProductName string  `json:"productName" gorm:"column:product_name;not null"`
	ShopId      int     `json:"shopId" gorm:"column:shop_id;not null"`
	ShopName    string  `json:"shopName" gorm:"column:shop_name;not null"`
	Num         int     `json:"num" gorm:"column:num;not null"`
	Price       float64 `json:"price" gorm:"column:price;not null"`
}

func NewCartUseCase(repo CartRepo, logger log.Logger) *CartUseCase {
	return &CartUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/cart")),
	}
}

func (cu *CartUseCase) CreateCart(ctx context.Context, req *v1.CreateCartReq) (*v1.CreateCartReply, error) {
	var ca = &Cart{
		Id:          0,
		UserUuid:    req.C.UserUuid,
		ImageUrl:    req.C.ImageUrl,
		ProductId:   int(req.C.ProductId),
		ProductName: req.C.ProductName,
		ShopName:    req.C.ShopName,
		ShopId:      int(req.C.ShopId),
		Num:         int(req.C.Num),
		Price:       float64(req.C.Price),
	}
	res, err := cu.repo.CreateCart(ctx, ca)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCartReply{
		C: createCart(res),
	}, nil
}

func (cu *CartUseCase) UpdateCart(ctx context.Context, req *v1.UpdateCartReq) (*v1.UpdateCartReply, error) {
	var ca = &Cart{
		Id:       int(req.C.Id),
		UserUuid: req.C.UserUuid,
		Num:      int(req.C.Num),
	}
	res, err := cu.repo.GetCart(ctx, ca)
	if err != nil {
		return nil, err
	}
	ca = &Cart{Id: int(req.C.Id), UserUuid: req.C.UserUuid, Num: int(req.C.Num), Price: float64(res.Price / float64(res.Num) * float64(req.C.Num))}
	res, err = cu.repo.UpdateCart(ctx, ca)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateCartReply{
		C: createCart(res),
	}, nil
}

func (cu *CartUseCase) DeleteCart(ctx context.Context, req *v1.DeleteCartReq) (*v1.DeleteCartReply, error) {
	var ca = &Cart{
		Id:       int(req.Id),
		UserUuid: req.UserUuid,
	}
	_, err := cu.repo.DeleteCart(ctx, ca)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteCartReply{}, nil
}

func (cu *CartUseCase) GetCart(ctx context.Context, req *v1.GetCartReq) (*v1.GetCartReply, error) {
	var ca = &Cart{
		Id:       int(req.Id),
		UserUuid: req.UserUuid,
	}
	res, err := cu.repo.GetCart(ctx, ca)
	if err != nil {
		return nil, err
	}
	return &v1.GetCartReply{C: createCart(res)}, nil
}

func (cu *CartUseCase) ListCart(ctx context.Context, req *v1.ListCartReq) (*v1.ListCartReply, error) {
	var qc = &QueryCart{
		UserUuid: req.UserUuid,
		Limit:    int(req.Limit),
		Offset:   int(req.Limit * (req.Page - 1)),
	}
	res, err := cu.repo.ListCart(ctx, qc)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListCartReply{CartList: make([]*v1.Cart, 0)}, nil
	}
	var caList = make([]*v1.Cart, len(res))
	for i, item := range res {
		caList[i] = createCart(&item)
	}
	return &v1.ListCartReply{CartList: caList}, nil
}

func createCart(ca *Cart) *v1.Cart {
	return &v1.Cart{
		Id:          int64(ca.Id),
		UserUuid:    ca.UserUuid,
		ImageUrl:    ca.ImageUrl,
		ProductId:   int64(ca.ProductId),
		ProductName: ca.ProductName,
		ShopId:      int64(ca.ShopId),
		ShopName:    ca.ShopName,
		Num:         int64(ca.Num),
		Price:       float32(ca.Price),
	}
}
