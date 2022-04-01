package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "nt-srv/api/notify/v1"
	"nt-srv/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewNtSrv)

type NtSrv struct {
	ntCase *biz.NtUseCase
	logger *log.Helper
	v1.UnsafeNtSrvServer
}

var _ v1.NtSrvServer = (*NtSrv)(nil)

func NewNtSrv(ntCase *biz.NtUseCase, logger log.Logger) *NtSrv {
	return &NtSrv{
		ntCase:            ntCase,
		logger:            log.NewHelper(log.With(logger, "module", "service/notice")),
		UnsafeNtSrvServer: &v1.UnimplementedNtSrvServer{},
	}
}
