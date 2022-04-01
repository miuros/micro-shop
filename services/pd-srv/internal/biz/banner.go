package biz

import (
	"context"
	v1 "pd-srv/api/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type BnUseCase struct {
	repo   BannerRepo
	logger *log.Helper
}

type Banner struct {
	Id          int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:name;omitempty"`
	ImageUrl    string `json:"imageUrl" gorm:"column:image_url;not null"`
	RedirectUrl string `json:"redirectUrl" gorm:"column:redirect_url;not null"`
}

type BannerRepo interface {
	CreateBn(context.Context, *Banner) (*Banner, error)
	UpdateBn(context.Context, *Banner) (*Banner, error)
	DeleteBn(context.Context, *Banner) (*Banner, error)
	GetBn(context.Context, int) (*Banner, error)
	ListBn(context.Context, int) ([]Banner, error)
}

func NewBnUseCase(repo BannerRepo, logger log.Logger) *BnUseCase {
	return &BnUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/banner")),
	}
}

func (bu *BnUseCase) CreateBn(ctx context.Context, req *v1.CreateBnReq) (*v1.CreateBnReply, error) {
	var bn = &Banner{
		Id:          0,
		Name:        req.Bn.Name,
		ImageUrl:    req.Bn.ImageUrl,
		RedirectUrl: req.Bn.RedirectUrl,
	}
	res, err := bu.repo.CreateBn(ctx, bn)
	if err != nil {
		return nil, err
	}
	return &v1.CreateBnReply{Bn: createBn(res)}, nil
}

func (bu *BnUseCase) UpdateBn(ctx context.Context, req *v1.UpdateBnReq) (*v1.UpdateBnReply, error) {
	var bn = &Banner{
		Id:          int(req.Bn.Id),
		Name:        req.Bn.Name,
		ImageUrl:    req.Bn.ImageUrl,
		RedirectUrl: req.Bn.RedirectUrl,
	}
	res, err := bu.repo.UpdateBn(ctx, bn)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateBnReply{Bn: createBn(res)}, nil
}

func (bu *BnUseCase) DeleteBn(ctx context.Context, req *v1.DeleteBnReq) (*v1.DeleteBnReply, error) {
	var bn = &Banner{
		Id: int(req.Id),
	}
	_, err := bu.repo.DeleteBn(ctx, bn)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteBnReply{}, nil
}

func (bu *BnUseCase) GetBn(ctx context.Context, req *v1.GetBnReq) (*v1.GetBnReply, error) {
	res, err := bu.repo.GetBn(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &v1.GetBnReply{Bn: createBn(res)}, nil
}

func (bu *BnUseCase) ListBn(ctx context.Context, req *v1.ListBnReq) (*v1.ListBnReply, error) {
	res, err := bu.repo.ListBn(ctx, int(req.Limit))
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &v1.ListBnReply{BnList: make([]*v1.Banner, 0)}, nil
	}
	var bnList = make([]*v1.Banner, len(res))
	for i, item := range res {
		bnList[i] = createBn(&item)
	}
	return &v1.ListBnReply{BnList: bnList}, nil
}

func createBn(bn *Banner) *v1.Banner {
	return &v1.Banner{
		Id:          int64(bn.Id),
		Name:        bn.Name,
		ImageUrl:    bn.ImageUrl,
		RedirectUrl: bn.RedirectUrl,
	}
}
