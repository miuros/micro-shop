package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"pd-srv/internal/biz"
)

type bannerRepo struct {
	data   *Data
	logger *log.Helper
}

func (br *bannerRepo) GetBn(ctx context.Context, i int) (*biz.Banner, error) {
	var bn = new(biz.Banner)
	err := br.data.maria.Model(&biz.Banner{}).Where("id=?", i).First(bn).Error
	if err != nil {
		return nil, err
	}
	return bn, nil
}

func (br *bannerRepo) CreateBn(ctx context.Context, bn *biz.Banner) (*biz.Banner, error) {
	err := br.data.maria.Create(bn).Error
	if err != nil {
		return nil, err
	}
	return bn, nil
}

func (br *bannerRepo) UpdateBn(ctx context.Context, bn *biz.Banner) (*biz.Banner, error) {
	err := br.data.maria.Model(&biz.Banner{}).Where("id=?", bn.Id).Update(bn).Error
	if err != nil {
		return nil, err
	}
	return bn, nil
}

func (br *bannerRepo) DeleteBn(ctx context.Context, bn *biz.Banner) (*biz.Banner, error) {
	err := br.data.maria.Model(&biz.Banner{}).Where("id=?", bn.Id).Delete(bn).Error
	if err != nil {
		return nil, err
	}
	return bn, nil
}

func (br *bannerRepo) ListBn(ctx context.Context, limit int) ([]biz.Banner, error) {
	var bnList []biz.Banner
	err := br.data.maria.Model(&biz.Banner{}).Limit(limit).Find(&bnList).Error
	if err != nil {
		return nil, err
	}
	return bnList, nil
}

func NewBannerRepo(data *Data, logger log.Logger) biz.BannerRepo {
	return &bannerRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/banner")),
	}
}
