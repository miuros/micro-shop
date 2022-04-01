package biz

import (
	"context"
	v1 "pd-srv/api/v1"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type ShopUseCase struct {
	repo   ShopRepo
	logger *log.Helper
}

type ShopRepo interface {
	CreateShop(context.Context, *Shop) (*Shop, error)
	UpdateShop(context.Context, *Shop) (*Shop, error)
	DeleteShop(context.Context, *Shop) (*Shop, error)
	GetShop(context.Context, int) (*Shop, error)
	ListShop(context.Context, *Query) ([]Shop, error)
	GetSpByUuid(context.Context, *Shop) (*Shop, error)
}

type Shop struct {
	Id        int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"column:name;unique"`
	UserUuid  string `json:"userUuid" gorm:"column:user_uuid;not null"`
	ImageUrl  string `json:"imageUrl" gorm:"column:image_url;omitempty"`
	Address   string `json:"address" gorm:"column:address;not null"`
	CreateAt  string `json:"createAt" gorm:"column:create_at;omitempty"`
	DeleteAt  string `json:"deleteAt" gorm:"column:delete_at;omitempty"`
	IsDeleted int    `json:"isDeleted" gorm:"column:is_deleted;default 0"`
}

func NewShopUseCase(repo ShopRepo, logger log.Logger) *ShopUseCase {
	return &ShopUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/shop")),
	}
}

func (su *ShopUseCase) CreateShop(ctx context.Context, spv1 *v1.CreateShopReq) (*v1.CreateShopReply, error) {
	var sp = &Shop{
		Id:        0,
		Name:      spv1.Sp.Name,
		UserUuid:  spv1.Sp.UserUuid,
		ImageUrl:  spv1.Sp.ImageUrl,
		Address:   spv1.Sp.Address,
		CreateAt:  time.Now().Format("2006-01-02:15-04"),
		IsDeleted: 0,
	}
	res, err := su.repo.CreateShop(ctx, sp)
	if err != nil {
		return nil, err
	}
	return &v1.CreateShopReply{Sp: createShop(res)}, nil
}

func (su *ShopUseCase) UpdateShop(ctx context.Context, spv1 *v1.UpdateShopReq) (*v1.UpdateShopReply, error) {
	var sp = &Shop{
		Id:       int(spv1.Sp.Id),
		Name:     spv1.Sp.Name,
		ImageUrl: spv1.Sp.ImageUrl,
		Address:  spv1.Sp.Address,
	}
	res, err := su.repo.UpdateShop(ctx, sp)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateShopReply{Sp: createShop(res)}, nil
}

func (su *ShopUseCase) DeleteShop(ctx context.Context, id int) (*v1.DeleteShopReply, error) {
	_, err := su.repo.DeleteShop(ctx, &Shop{Id: id})
	if err != nil {
		return nil, err
	}
	return &v1.DeleteShopReply{}, nil
}

func (su *ShopUseCase) GetShop(ctx context.Context, req *v1.GetShopReq) (*v1.GetShopReply, error) {
	res, err := su.repo.GetShop(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &v1.GetShopReply{
		Sp: createShop(res),
	}, nil
}

func (su *ShopUseCase) ListShop(ctx context.Context, qv1 *v1.ListShopReq) (*v1.ListShopReply, error) {
	var q = &Query{
		Limit:  int(qv1.Limit),
		Offset: int(qv1.Limit * (qv1.Page - 1)),
		Name:   qv1.Name,
	}
	res, err := su.repo.ListShop(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListShopReply{SpList: make([]*v1.Shop, 0)}, nil
	}
	var spList = make([]*v1.Shop, len(res))
	for i, item := range res {
		spList[i] = createShop(&item)
	}
	return &v1.ListShopReply{SpList: spList}, nil
}

func (su *ShopUseCase) GetSpByUuid(ctx context.Context, req *v1.GetSpByUuidReq) (*v1.GetSpByUuidReply, error) {
	var sp = &Shop{
		UserUuid: req.UserUuid,
	}
	res, err := su.repo.GetSpByUuid(ctx, sp)
	if err != nil {
		return nil, err
	}
	return &v1.GetSpByUuidReply{Sp: createShop(res)}, nil
}

func createShop(sp *Shop) *v1.Shop {
	return &v1.Shop{
		Id:        int64(sp.Id),
		Name:      sp.Name,
		ImageUrl:  sp.ImageUrl,
		Address:   sp.Address,
		CreateAt:  sp.CreateAt,
		UserUuid:  sp.UserUuid,
		IsDeleted: int64(sp.IsDeleted),
		DeleteAt:  sp.DeleteAt,
	}
}
