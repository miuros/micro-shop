package service

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pdv1 "micro-shop/api/pd/v1"
	"os"
)

type PdSrv struct {
	pr     pdv1.PdServiceClient
	logger *log.Logger
}

func NewPdSrv(logger *log.Logger, conn grpc.ClientConnInterface) *PdSrv {
	pr := pdv1.NewPdServiceClient(conn)
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("service/product:")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	return &PdSrv{
		pr:     pr,
		logger: logger,
	}
}

func (ps *PdSrv) CreatePd(ctx context.Context, req *pdv1.CreatePdReq) (*pdv1.CreatePdReply, error) {
	res, err := ps.pr.CreatePd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) UpdatePd(ctx context.Context, req *pdv1.UpdatePdReq) (*pdv1.UpdatePdReply, error) {
	res, err := ps.pr.UpdatePd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) DeletePd(ctx context.Context, req *pdv1.DeletePdReq) (*pdv1.DeletePdReply, error) {
	res, err := ps.pr.DeletePd(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ps *PdSrv) GetPd(ctx context.Context, req *pdv1.GetPdReq) (*pdv1.GetPdReply, error) {
	res, err := ps.pr.GetPd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListPd(ctx context.Context, req *pdv1.ListPdReq) (*pdv1.ListPdReply, error) {
	res, err := ps.pr.ListPd(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListPdForSp(ctx context.Context, req *pdv1.ListForSpReq) (*pdv1.ListForSpReply, error) {
	res, err := ps.pr.ListForSp(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListPdByCg(ctx context.Context, req *pdv1.ListPdByCiReq) (*pdv1.ListPdByCiReply, error) {
	res, err := ps.pr.ListPdByCi(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) FindPdByName(ctx context.Context, req *pdv1.ListPdReq) (*pdv1.ListPdReply, error) {
	res, err := ps.pr.FindPdByName(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) CreateSp(ctx context.Context, req *pdv1.CreateShopReq) (*pdv1.CreateShopReply, error) {
	res, err := ps.pr.CreateShop(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) UpdateSp(ctx context.Context, req *pdv1.UpdateShopReq) (*pdv1.UpdateShopReply, error) {
	res, err := ps.pr.UpdateShop(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) DeleteSp(ctx context.Context, req *pdv1.DeleteShopReq) (*pdv1.DeleteShopReply, error) {
	res, err := ps.pr.DeleteShop(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) GetSp(ctx context.Context, req *pdv1.GetShopReq) (*pdv1.GetShopReply, error) {
	res, err := ps.pr.GetShop(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) GetSpByUuid(ctx context.Context, req *pdv1.GetSpByUuidReq) (*pdv1.GetSpByUuidReply, error) {
	res, err := ps.pr.GetSpByUuid(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListSp(ctx context.Context, req *pdv1.ListShopReq) (*pdv1.ListShopReply, error) {
	res, err := ps.pr.ListShop(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) GetCart(ctx context.Context, req *pdv1.GetCartReq) (*pdv1.GetCartReply, error) {
	res, err := ps.pr.GetCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) CreateCt(ctx context.Context, req *pdv1.CreateCartReq) (*pdv1.CreateCartReply, error) {
	res, err := ps.pr.CreateCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) UpdateCt(ctx context.Context, req *pdv1.UpdateCartReq) (*pdv1.UpdateCartReply, error) {
	res, err := ps.pr.UpdateCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) DeleteCt(ctx context.Context, req *pdv1.DeleteCartReq) (*pdv1.DeleteCartReply, error) {
	res, err := ps.pr.DeleteCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListCt(ctx context.Context, req *pdv1.ListCartReq) (*pdv1.ListCartReply, error) {
	res, err := ps.pr.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) GetBn(ctx context.Context, req *pdv1.GetBnReq) (*pdv1.GetBnReply, error) {
	res, err := ps.pr.GetBn(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) CreateBn(ctx context.Context, req *pdv1.CreateBnReq) (*pdv1.CreateBnReply, error) {
	res, err := ps.pr.CreateBn(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) UpdateBn(ctx context.Context, req *pdv1.UpdateBnReq) (*pdv1.UpdateBnReply, error) {
	res, err := ps.pr.UpdateBn(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) DeleteBn(ctx context.Context, req *pdv1.DeleteBnReq) (*pdv1.DeleteBnReply, error) {
	res, err := ps.pr.DeleteBn(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListBn(ctx context.Context, req *pdv1.ListBnReq) (*pdv1.ListBnReply, error) {
	res, err := ps.pr.ListBn(ctx, req)
	if err != nil {
		return nil, err

	}
	return res, nil
}

func (ps *PdSrv) CreateCg(ctx context.Context, req *pdv1.CreateCgReq) (*pdv1.CreateCgReply, error) {
	res, err := ps.pr.CreateCg(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) UpdateCg(ctx context.Context, req *pdv1.UpdateCgReq) (*pdv1.UpdateCgReply, error) {
	res, err := ps.pr.UpdateCg(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) DeleteCg(ctx context.Context, req *pdv1.DeleteCgReq) (*pdv1.DeleteCgReply, error) {
	res, err := ps.pr.DeleteCg(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) GetCg(ctx context.Context, req *pdv1.GetCgReq) (*pdv1.GetCgReply, error) {
	res, err := ps.pr.GetCg(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ps *PdSrv) ListCg(ctx context.Context, req *pdv1.ListCgReq) (*pdv1.ListCgReply, error) {
	res, err := ps.pr.ListCg(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
