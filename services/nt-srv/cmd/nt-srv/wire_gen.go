// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"nt-srv/internal/biz"
	"nt-srv/internal/conf"
	"nt-srv/internal/data"
	"nt-srv/internal/server"
	"nt-srv/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	ntRepo := data.NewNtRepo(dataData, logger)
	ntUseCase := biz.NewNtUseCase(ntRepo, logger)
	ntSrv := service.NewNtSrv(ntUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, ntSrv, logger)
	grpcServer := server.NewGRPCServer(confServer, ntSrv, logger)
	app := newApp(logger, confData, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}