package service

import (
	"context"
	"google.golang.org/grpc"
	"log"
	odv1 "micro-shop/api/od/v1"
	"os"
)

type OdSrv struct {
	or     odv1.OdSrvClient
	logger *log.Logger
}

func NewOdSrv(logger *log.Logger, conn grpc.ClientConnInterface) *OdSrv {
	logger.SetPrefix("service/order")
	logger.SetOutput(os.Stdout)
	logger.SetFlags(log.Ltime | log.Lshortfile)
	or := odv1.NewOdSrvClient(conn)
	return &OdSrv{
		or:     or,
		logger: logger,
	}
}

func (os *OdSrv) CreateOd(ctx context.Context, req *odv1.CreateOrderReq) (*odv1.CreateOrderReply, error) {
	res, err := os.or.CreateOd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) UpdateOd(ctx context.Context, req *odv1.UpdateOrderReq) (*odv1.UpdateOrderReply, error) {
	res, err := os.or.UpdateOd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) DeleteOd(ctx context.Context, req *odv1.DeleteOrderReq) (*odv1.DeleteOrderReply, error) {
	res, err := os.or.DeleteOd(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (os *OdSrv) GetOd(ctx context.Context, req *odv1.GetOrderReq) (*odv1.GetOrderReply, error) {
	res, err := os.or.GetOd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) ListOd(ctx context.Context, req *odv1.ListOrderReq) (*odv1.ListOrderReply, error) {
	res, err := os.or.ListOd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) ListOdByCateId(ctx context.Context, req *odv1.ListOrderByCateIdReq) (*odv1.ListOrderByCateIdReply, error) {
	res, err := os.or.ListOdByCateId(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) ListOdForSp(ctx context.Context, req *odv1.ListOdForSReq) (*odv1.ListOdForSpReply, error) {
	res, err := os.or.ListOdForShopper(ctx, req)
	if err != nil {
		return nil, err

	}
	return res, nil
}

func (os *OdSrv) GetSc(ctx context.Context, req *odv1.GetStockReq) (*odv1.GetStockReply, error) {
	res, err := os.or.GetStock(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) CreateSc(ctx context.Context, req *odv1.CreateStockReq) (*odv1.CreateStockReply, error) {
	res, err := os.or.CreateStock(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) UpdateSc(ctx context.Context, req *odv1.UpdateStockReq) (*odv1.UpdateStockReply, error) {
	res, err := os.or.UpdateStock(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) DeleteSc(ctx context.Context, req *odv1.DeleteStockReq) (*odv1.DeleteStockReply, error) {
	res, err := os.or.DeleteStock(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) CreateCate(ctx context.Context, req *odv1.CreateCateReq) (*odv1.CreateCateReply, error) {
	res, err := os.or.CreateCate(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (os *OdSrv) PayOd(ctx context.Context, req *odv1.PayOdReq) (*odv1.PayOdReply, error) {
	res, err := os.or.PayOd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
