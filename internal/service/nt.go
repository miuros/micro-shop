package service

import (
	"context"
	"google.golang.org/grpc"
	"log"
	ntv1 "micro-shop/api/notice/v1"
	"os"
)

type NtSrv struct {
	logger *log.Logger
	nr     ntv1.NtSrvClient
}

func NewNtSrv(logger *log.Logger, conn grpc.ClientConnInterface) *NtSrv {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("service/notice")
	logger.SetFlags(log.Lshortfile | log.Ltime)
	nr := ntv1.NewNtSrvClient(conn)
	return &NtSrv{
		nr:     nr,
		logger: logger,
	}
}

func (ns *NtSrv) CreateNt(ctx context.Context, req *ntv1.CreateNtReq) (*ntv1.CreateNtReply, error) {
	res, err := ns.nr.CreateNt(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ns *NtSrv) UpdateNtStatus(ctx context.Context, req *ntv1.UpdateStatusReq) (*ntv1.UpdateStatusReply, error) {
	res, err := ns.nr.UpdateStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ns *NtSrv) DeleteNt(ctx context.Context, req *ntv1.DeleteNtReq) (*ntv1.DeleteNtReply, error) {
	res, err := ns.nr.DeleteNt(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ns *NtSrv) ListNt(ctx context.Context, req *ntv1.ListNtReq) (*ntv1.ListNtReply, error) {
	res, err := ns.nr.ListNt(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ns *NtSrv) GetNt(ctx context.Context, req *ntv1.GetNtReq) (*ntv1.GetNtReply, error) {
	res, err := ns.nr.GetNt(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
