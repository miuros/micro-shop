package biz

import (
	v1 "comment-srv/api/comment/v1"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Comment struct {
	Id         int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ProductId  int    `json:"productId" gorm:"column:product_id;not null"`
	UserUuid   string `json:"userUuid" gorm:"column:user_uuid;not null"`
	Content    string `json:"content" gorm:"column:content;not null"`
	ToUserUuid string `json:"toUserUuid" gorm:"column:to_user_uuid;omitempty"`
	CreateAt   string `json:"createAt" gorm:"column:create_at"`
	UpdateAt   string `json:"updateAt" gorm:"column:update_at"`
	DeleteAt   string `json:"deleteAt" gorm:"column:delete_at"`
	IsDeleted  int    `json:"isDeleted" gorm:"column:is_deleted;default 0"`
}

type QueryCm struct {
	Offset    int `json:"offset"`
	Limit     int `json:"limit"`
	ProductId int `json:"productId"`
}

type CommentRepo interface {
	CreateCm(context.Context, *Comment) (*Comment, error)
	UpdateCm(context.Context, *Comment) (*Comment, error)
	DeleteCm(context.Context, *Comment) (*Comment, error)
	GetCm(context.Context, *Comment) (*Comment, error)
	ListCm(context.Context, *QueryCm) ([]Comment, error)
}

type CmUseCase struct {
	repo   CommentRepo
	logger *log.Helper
}

func NewCmUseCase(repo CommentRepo, logger log.Logger) *CmUseCase {
	return &CmUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/comment")),
	}
}

func (cu *CmUseCase) CreateCm(ctx context.Context, req *v1.CreateCmReq) (*v1.CreateCmReply, error) {
	var cm = &Comment{
		Id:         0,
		ProductId:  int(req.Cm.ProductId),
		UserUuid:   req.Cm.UserUuid,
		ToUserUuid: req.Cm.ToUserUuid,
		Content:    req.Cm.Content,
		CreateAt:   time.Now().Format("2006-01-02:15-04"),
		IsDeleted:  0,
	}
	cm, err := cu.repo.CreateCm(ctx, cm)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCmReply{Cm: createCm(cm)}, nil
}

func (cu *CmUseCase) UpdateCm(ctx context.Context, req *v1.UpdateCmReq) (*v1.UpdateCmReply, error) {
	var cm = &Comment{
		Id:       int(req.Cm.Id),
		Content:  req.Cm.Content,
		UpdateAt: time.Now().String(),
	}
	cm, err := cu.repo.UpdateCm(ctx, cm)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateCmReply{Cm: createCm(cm)}, nil
}

func (cu *CmUseCase) DeleteCm(ctx context.Context, req *v1.DeleteCmReq) (*v1.DeleteCmReply, error) {
	var cm = &Comment{
		Id:       int(req.Id),
		DeleteAt: time.Now().Format("2006-01-02:15-04"),
	}
	_, err := cu.repo.DeleteCm(ctx, cm)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteCmReply{}, nil
}

func (cu *CmUseCase) GetCm(ctx context.Context, req *v1.GetCmReq) (*v1.GetCmReply, error) {
	var cm = &Comment{Id: int(req.Id)}
	cm, err := cu.repo.GetCm(ctx, cm)
	if err != nil {
		return nil, err
	}
	return &v1.GetCmReply{Cm: createCm(cm)}, nil
}

func (cu *CmUseCase) ListCm(ctx context.Context, req *v1.ListCmReq) (*v1.ListCmReply, error) {
	var qc = &QueryCm{
		Offset:    int(req.Limit * (req.Page - 1)),
		Limit:     int(req.Limit),
		ProductId: int(req.ProductId),
	}
	cmList, err := cu.repo.ListCm(ctx, qc)
	if err != nil {
		return nil, err
	}
	if len(cmList) == 0 {
		return &v1.ListCmReply{CmList: make([]*v1.Comment, 0)}, nil
	}
	var res = make([]*v1.Comment, len(cmList))
	for i, item := range cmList {
		res[i] = createCm(&item)
	}
	return &v1.ListCmReply{CmList: res}, nil
}

func createCm(cm *Comment) *v1.Comment {
	return &v1.Comment{
		Id:         int64(cm.Id),
		UserUuid:   cm.UserUuid,
		Content:    cm.Content,
		ToUserUuid: cm.ToUserUuid,
		ProductId:  int64(cm.ProductId),
		CreateAt:   cm.CreateAt,
		DeleteAt:   cm.DeleteAt,
		UpdateAt:   cm.UpdateAt,
		IsDeleted:  int64(cm.IsDeleted),
	}
}
