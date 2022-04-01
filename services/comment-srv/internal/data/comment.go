package data

import (
	"comment-srv/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type cmRepo struct {
	data   *Data
	logger *log.Helper
}

func NewCmRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &cmRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/comment")),
	}
}

func (cr *cmRepo) CreateCm(ctx context.Context, cm *biz.Comment) (*biz.Comment, error) {
	err := cr.data.maria.Create(cm).Error
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (cr *cmRepo) UpdateCm(ctx context.Context, cm *biz.Comment) (*biz.Comment, error) {
	err := cr.data.maria.Model(&biz.Comment{}).Where("id=?", cm.Id).Updates(cm).Error
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (cr *cmRepo) DeleteCm(ctx context.Context, cm *biz.Comment) (*biz.Comment, error) {
	err := cr.data.maria.Model(&biz.Comment{}).Where("id=?", cm.Id).Updates(cm).Error
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (cr *cmRepo) GetCm(ctx context.Context, cm *biz.Comment) (*biz.Comment, error) {
	err := cr.data.maria.Model(&biz.Comment{}).Where("id=?", cm.Id).First(cm).Error
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (cr *cmRepo) ListCm(ctx context.Context, qc *biz.QueryCm) ([]biz.Comment, error) {
	var cmList []biz.Comment
	err := cr.data.maria.Model(&biz.Comment{}).Where("product_id=?", qc.ProductId).Offset(qc.Offset).Limit(qc.Limit).Find(&cmList).Error
	if err != nil {
		return nil, err
	}
	return cmList, nil
}

var _ biz.CommentRepo = (*cmRepo)(nil)
