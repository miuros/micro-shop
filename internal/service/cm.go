package service

import (
	"context"
	"google.golang.org/grpc"
	"log"
	cmv1 "micro-shop/api/cm/v1"
	"os"
)

type CmSrv struct {
	logger *log.Logger
	cr     cmv1.CmServiceClient
}

func NewCmSrv(logger *log.Logger, conn grpc.ClientConnInterface) *CmSrv {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("service/comment")
	logger.SetFlags(log.Lshortfile | log.Ltime)
	cr := cmv1.NewCmServiceClient(conn)
	return &CmSrv{
		cr:     cr,
		logger: logger,
	}
}

func (cs *CmSrv) CreateCm(ctx context.Context, req *cmv1.CreateCmReq) (*cmv1.CreateCmReply, error) {
	res, err := cs.cr.CreateCm(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cs *CmSrv) UpdateCm(ctx context.Context, req *cmv1.UpdateCmReq) (*cmv1.UpdateCmReply, error) {
	res, err := cs.cr.UpdateCm(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cs *CmSrv) DeleteCm(ctx context.Context, req *cmv1.DeleteCmReq) (*cmv1.DeleteCmReply, error) {
	res, err := cs.cr.DeleteCm(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cs *CmSrv) GetCm(ctx context.Context, req *cmv1.GetCmReq) (*cmv1.GetCmReply, error) {
	res, err := cs.cr.GetCm(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cs *CmSrv) ListCm(ctx context.Context, req *cmv1.ListCmReq) (*cmv1.ListCmReply, error) {
	res, err := cs.cr.ListCm(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
