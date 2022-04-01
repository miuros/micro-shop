package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "pd-srv/api/v1"
	"pd-srv/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewPdService)

type PdService struct {
	pdCase   *biz.PdUseCase
	cartCase *biz.CartUseCase
	bnCase   *biz.BnUseCase
	shopCase *biz.ShopUseCase
	cgCase   *biz.CgUseCase
	v1.UnsafePdServiceServer
	logger *log.Helper
}

func (ps *PdService) ListForSp(ctx context.Context, req *v1.ListForSpReq) (*v1.ListForSpReply, error) {

	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}

	if req.ShopId < 1 {
		return &v1.ListForSpReply{PdList: make([]*v1.Product, 0)}, fmt.Errorf("shop id is nil")
	}
	res, err := ps.pdCase.ListForSp(ctx, req)
	if err != nil {
		ps.logger.Errorf("list product for shop error:%s", err.Error())
		return &v1.ListForSpReply{PdList: make([]*v1.Product, 0)}, fmt.Errorf("internal server error")
	}
	return res, nil
}

func NewPdService(pdCase *biz.PdUseCase, bnCase *biz.BnUseCase, cartCase *biz.CartUseCase, shopCase *biz.ShopUseCase, cgCase *biz.CgUseCase, logger log.Logger) *PdService {
	return &PdService{
		pdCase:                pdCase,
		shopCase:              shopCase,
		cartCase:              cartCase,
		bnCase:                bnCase,
		cgCase:                cgCase,
		UnsafePdServiceServer: &v1.UnimplementedPdServiceServer{},
		logger:                log.NewHelper(log.With(logger, "module", "pd/service")),
	}
}

var _ v1.PdServiceServer = (*PdService)(nil)
