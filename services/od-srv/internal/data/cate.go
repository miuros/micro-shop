package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"od-srv/internal/biz"
	"time"
)

type cateRepo struct {
	data   *Data
	logger *log.Helper
}

func (cr *cateRepo) CreateCate(ctx context.Context, cate *biz.Cate) (*biz.Cate, error) {
	err := cr.data.maria.Create(cate).Error
	if err != nil {
		return nil, err
	}
	return cate, nil
}

func (cr *cateRepo) UpdateCate(ctx context.Context, cate *biz.Cate) (*biz.Cate, error) {
	err := cr.data.maria.Model(&biz.Cate{}).Where("id=? and user_uuid=?", cate.Id, cate.UserUuid).Update(cate).Error
	if err != nil {
		return nil, err
	}
	return cate, nil
}

func (cr *cateRepo) DeleteCate(ctx context.Context, cate *biz.Cate) (*biz.Cate, error) {
	err := cr.data.maria.Model(&biz.Cate{}).Where("id=?  and user_uuid=?", cate.Id, cate.UserUuid).Update(&biz.Cate{IsDeleted: 1, DeleteAt: time.Now().String()}).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (cr *cateRepo) GetCate(ctx context.Context, cate *biz.Cate) (*biz.Cate, error) {
	err := cr.data.maria.Model(&biz.Cate{}).Where("id=?  and user_uuid=?", cate.Id, cate.UserUuid).First(cate).Error
	if err != nil {
		return nil, err
	}
	return cate, nil
}

func (cr *cateRepo) ListCate(ctx context.Context, qc *biz.QueryCate) ([]biz.Cate, error) {
	var caList []biz.Cate
	err := cr.data.maria.Model(&biz.Cate{}).Where("user_uuid=?", qc.UserUuid).Offset(qc.Offset).Limit(qc.Limit).Find(&caList).Error
	if err != nil {
		return nil, err
	}
	return caList, nil
}

var _ biz.CateRepo = (*cateRepo)(nil)

func NewCateRepo(data *Data, logger log.Logger) biz.CateRepo {
	return &cateRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/cate")),
	}
}
