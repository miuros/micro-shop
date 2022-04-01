package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"pd-srv/internal/biz"
)

type cgRepo struct {
	data   *Data
	logger *log.Helper
}

func NewCgRepo(logger log.Logger, data *Data) biz.CgRepo {
	return &cgRepo{
		logger: log.NewHelper(log.With(logger, "data/category")),
		data:   data,
	}
}

func (cr *cgRepo) CreateCg(ctx context.Context, cg *biz.Category) (*biz.Category, error) {
	err := cr.data.maria.Create(cg).Error
	if err != nil {
		return nil, err
	}
	return cg, nil
}

func (cr *cgRepo) UpdateCg(ctx context.Context, cg *biz.Category) (*biz.Category, error) {
	err := cr.data.maria.Model(&biz.Category{}).Where("id=?", cg.Id).Update(cg).Error
	if err != nil {
		return nil, err
	}
	return cg, nil
}

func (cr *cgRepo) DeleteCg(ctx context.Context, cg *biz.Category) (*biz.Category, error) {
	err := cr.data.maria.Model(&biz.Category{}).Where("id=?", cg.Id).Delete(cg).Error
	if err != nil {
		return nil, err
	}
	return cg, nil
}

func (cr *cgRepo) GetCg(ctx context.Context, cg *biz.Category) (*biz.Category, error) {
	err := cr.data.maria.Model(&biz.Category{}).Where("id=?", cg.Id).First(cg).Error
	if err != nil {
		return nil, err
	}
	return cg, nil
}

func (cr *cgRepo) ListCg(ctx context.Context, qc *biz.QueryCg) ([]biz.Category, error) {
	var cgList []biz.Category
	err := cr.data.maria.Model(&biz.Category{}).Offset(qc.Offset).Limit(qc.Limit).Find(&cgList).Error
	if err != nil {
		return nil, err
	}
	return cgList, nil
}

var _ biz.CgRepo = (*cgRepo)(nil)
