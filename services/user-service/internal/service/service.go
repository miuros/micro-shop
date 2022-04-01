package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"user-srv/api/user/v1"
	"user-srv/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService)

var _ v1.UserServer = (*UserService)(nil)

type UserService struct {
	addressCase *biz.AddressUseCase
	userCase    *biz.UserUseCase
	v1.UnsafeUserServer
	logger *log.Helper
}

func NewUserService(ad *biz.AddressUseCase, u *biz.UserUseCase, logger log.Logger) *UserService {

	us := &UserService{
		addressCase:      ad,
		UnsafeUserServer: v1.UnimplementedUserServer{},
		userCase:         u,
		logger:           log.NewHelper(log.With(logger, "module", "service")),
	}
	return us
}
