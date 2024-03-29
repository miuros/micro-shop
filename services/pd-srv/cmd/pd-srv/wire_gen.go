// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"pd-srv/internal/biz"
	"pd-srv/internal/conf"
	"pd-srv/internal/data"
	"pd-srv/internal/server"
	"pd-srv/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	productRepo := data.NewProductRepo(dataData, logger)
	pdUseCase := biz.NewPdUseCase(productRepo, logger)
	bannerRepo := data.NewBannerRepo(dataData, logger)
	bnUseCase := biz.NewBnUseCase(bannerRepo, logger)
	cartRepo := data.NewCartRepo(dataData, logger)
	cartUseCase := biz.NewCartUseCase(cartRepo, logger)
	shopRepo := data.NewShopRepo(dataData, logger)
	shopUseCase := biz.NewShopUseCase(shopRepo, logger)
	cgRepo := data.NewCgRepo(logger, dataData)
	cgUseCase := biz.NewCgUseCase(logger, cgRepo)
	pdService := service.NewPdService(pdUseCase, bnUseCase, cartUseCase, shopUseCase, cgUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, pdService, logger)
	grpcServer := server.NewGRPCServer(confServer, pdService, logger)
	app := newApp(logger, confData, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
