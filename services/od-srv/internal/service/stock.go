package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	v1 "od-srv/api/order/v1"
)

func (os *OdSrv) getUserUuidByPdId(ctx context.Context, id int64) (string, error) {
	pd, err := os.GetPd(ctx, id)
	if err != nil {
		os.logger.Errorf("get product error:%s", err.Error())
		return "", err
	}
	sp, err := os.GetSp(ctx, pd.ShopId)
	if err != nil {
		os.logger.Errorf("get shop error:%s", err.Error())
		return "", err
	}
	return sp.UserUuid, nil
}

func (os *OdSrv) CreateStock(ctx context.Context, req *v1.CreateStockReq) (*v1.CreateStockReply, error) {
	if len(req.Stock.UserUuid) == 0 {
		return &v1.CreateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("user uuid is nil")
	}
	if req.Stock.Storage < 1 {
		return &v1.CreateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("storage must be greater than 1")
	}
	if req.Stock.ProductId < 1 {
		return &v1.CreateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("product is wrong")
	}
	userUuid, err := os.getUserUuidByPdId(ctx, int64(req.Stock.ProductId))
	if err != nil {
		return &v1.CreateStockReply{Stock: &v1.StockInfo{}}, err
	}
	if userUuid != req.Stock.UserUuid {
		return &v1.CreateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("user uuid is not right")
	}
	res, err := os.scCase.CreateSc(ctx, req)
	if err != nil {
		os.logger.Errorf("create stock error:%s", err.Error())
		return &v1.CreateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("failed to create stock ,internal server error")
	}
	return res, nil
}

func (os *OdSrv) UpdateStock(ctx context.Context, req *v1.UpdateStockReq) (*v1.UpdateStockReply, error) {
	if len(req.Stock.UserUuid) == 0 {
		//return &v1.CreateStockReply{Stock: &v1.StockInfo{}},fmt.Errorf()
	}
	if req.Stock.Storage < 1 {
		return &v1.UpdateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("storage must be greater than 1")
	}
	if req.Stock.ProductId < 1 {
		return &v1.UpdateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("product is wrong")
	}
	userUuid, err := os.getUserUuidByPdId(ctx, int64(req.Stock.ProductId))
	if err != nil {
		return &v1.UpdateStockReply{Stock: &v1.StockInfo{}}, err
	}
	if userUuid != req.Stock.UserUuid {
		return &v1.UpdateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("user uuid is not right")
	}
	res, err := os.scCase.UpdateSc(ctx, req)
	if err != nil {
		os.logger.Errorf("update stock error:%s", err.Error())
		return &v1.UpdateStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("failed to updte stock ,internal server error")
	}
	return res, nil
}

func (os *OdSrv) DeleteStock(ctx context.Context, req *v1.DeleteStockReq) (*v1.DeleteStockReply, error) {
	if req.ProductId < 1 {
		return &v1.DeleteStockReply{}, fmt.Errorf("product is wrong")
	}

	userUuid, err := os.getUserUuidByPdId(ctx, int64(req.ProductId))
	if err != nil {
		return &v1.DeleteStockReply{}, err
	}
	if userUuid != req.UserUuid {
		return &v1.DeleteStockReply{}, fmt.Errorf("user uuid is not right")
	}
	_, err = os.scCase.DeleteSc(ctx, req)
	if err != nil {
		os.logger.Errorf("delete stock error:%s", err.Error())
		return &v1.DeleteStockReply{}, fmt.Errorf("failed to delete stock ,internal server error")
	}
	return &v1.DeleteStockReply{}, nil
}

func (os *OdSrv) GetStock(ctx context.Context, req *v1.GetStockReq) (*v1.GetStockReply, error) {
	if req.ProductId < 1 {
		return &v1.GetStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("product id is nil")
	}
	res, err := os.scCase.GetSc(ctx, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("so such a stock")
		}
		os.logger.Errorf("get stock error:%s", err.Error())
		return &v1.GetStockReply{Stock: &v1.StockInfo{}}, fmt.Errorf("failed to get stock ,internal server error")
	}
	return res, nil
}
