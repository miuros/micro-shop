package service

import (
	v1 "comment-srv/api/comment/v1"
	"comment-srv/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewCmService)

type CmService struct {
	v1.UnsafeCmServiceServer
	cmCase *biz.CmUseCase
	logger *log.Helper
}

func NewCmService(cmCase *biz.CmUseCase, logger log.Logger) *CmService {
	return &CmService{
		cmCase:                cmCase,
		logger:                log.NewHelper(log.With(logger, "module", "service/comment")),
		UnsafeCmServiceServer: &v1.UnimplementedCmServiceServer{},
	}
}

var _ v1.CmServiceServer = (*CmService)(nil)
