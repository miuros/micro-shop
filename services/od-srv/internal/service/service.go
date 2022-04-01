package service

import (
	"context"
	"fmt"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	ggrpc "google.golang.org/grpc"
	v1 "od-srv/api/order/v1"
	"od-srv/internal/biz"
	"od-srv/internal/conf"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewOdSrv)

type OdSrv struct {
	logger *log.Helper
	odCase *biz.OdUseCase
	caCase *biz.CateUseCase
	scCase *biz.ScUseCase
	pdSrv  v1.PdServiceClient
	v1.UnsafeOdSrvServer
}

func NewOdSrv(odCase *biz.OdUseCase, caCase *biz.CateUseCase, scCase *biz.ScUseCase, confData *conf.Data, logger log.Logger) *OdSrv {
	client, err := api.NewClient(&api.Config{
		Address: confData.Consul.Addr,
	})
	if err != nil {
		panic(err)
	}
	cli := consul.New(client)
	srv, _, err := client.Agent().Service(confData.PdSrvName, nil)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(context.Background(), grpc.WithEndpoint(fmt.Sprintf("%s:%d", srv.Address, srv.Port)), grpc.WithDiscovery(cli), grpc.WithOptions(ggrpc.WithInsecure()))
	if err != nil {
		panic(err)
	}
	pdClient := v1.NewPdServiceClient(conn)

	return &OdSrv{
		odCase:            odCase,
		caCase:            caCase,
		scCase:            scCase,
		pdSrv:             pdClient,
		logger:            log.NewHelper(log.With(logger, "module", "service/order")),
		UnsafeOdSrvServer: &v1.UnimplementedOdSrvServer{},
	}
}

var _ v1.OdSrvServer = (*OdSrv)(nil)
