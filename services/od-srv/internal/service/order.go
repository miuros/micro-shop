package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	v1 "od-srv/api/order/v1"
	"od-srv/internal/data"
)

func (os *OdSrv) CreateOd(ctx context.Context, req *v1.CreateOrderReq) (*v1.CreateOrderReply, error) {
	if len(req.Item.UserUuid) == 0 {
		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("user uuid is nil")
	}

	if req.Item.AddressId < 1 {
		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("address is nil")
	}
	if req.Item.ProductId < 1 {
		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("product is nil")
	}
	if req.Item.CateId < 1 {
		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("cate  is nil")
	}
	if req.Item.Number < 1 {
		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("number must greater than 1")
	}
	if req.Item.Price < 1 {

		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("price error")
	}
	res, err := os.odCase.CreateOd(ctx, req.Item)
	if err != nil {
		if err == data.ERRStorageNotEnough {
			return &v1.CreateOrderReply{Item: &v1.Item{}}, err
		}
		os.logger.Errorf("create order error:%s", err.Error())
		return &v1.CreateOrderReply{Item: &v1.Item{}}, fmt.Errorf("failed to create order,internal server error")
	}
	return &v1.CreateOrderReply{Item: res}, nil

}

func (os *OdSrv) UpdateOd(ctx context.Context, req *v1.UpdateOrderReq) (*v1.UpdateOrderReply, error) {
	if req.AddressId < 1 {
		return &v1.UpdateOrderReply{Item: &v1.Item{}}, fmt.Errorf("address is nill")
	}
	if len(req.UserUuid) == 0 {
		return &v1.UpdateOrderReply{Item: &v1.Item{}}, fmt.Errorf("user uuid is nill")
	}
	od, err := os.odCase.GetOd(ctx, &v1.Item{Id: int64(req.Id), UserUuid: req.UserUuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.UpdateOrderReply{Item: &v1.Item{}}, fmt.Errorf("no such a order")
		}
		os.logger.Errorf("get order error:%s", err.Error())
		return &v1.UpdateOrderReply{Item: &v1.Item{}}, fmt.Errorf("internal server error")
	}
	if od.Status != 1 {
		return &v1.UpdateOrderReply{Item: &v1.Item{}}, fmt.Errorf("status not right")

	}
	res, err := os.odCase.UpdateOd(ctx, req)
	if err != nil {
		os.logger.Errorf("update order error:%s", err.Error())
		return &v1.UpdateOrderReply{Item: &v1.Item{}}, fmt.Errorf("failed to update order,internal server error")
	}
	return &v1.UpdateOrderReply{Item: res}, nil
}

func (os *OdSrv) DeleteOd(ctx context.Context, req *v1.DeleteOrderReq) (*v1.DeleteOrderReply, error) {
	if len(req.UserUuid) == 0 {
		return &v1.DeleteOrderReply{}, fmt.Errorf("user uuid is nil")
	}
	od, err := os.odCase.GetOd(ctx, &v1.Item{Id: req.Id, UserUuid: req.UserUuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeleteOrderReply{}, fmt.Errorf("no such a order")
		}
		os.logger.Errorf("get order error:%s", err.Error())
		return &v1.DeleteOrderReply{}, fmt.Errorf("failed to delete order,internal server error")
	}
	if od.Status != 1 {
		return &v1.DeleteOrderReply{}, fmt.Errorf("status not right")
	}
	_, err = os.odCase.DeleteOd(ctx, od)
	if err != nil {
		os.logger.Errorf("delete order error:%s", err.Error())
		return &v1.DeleteOrderReply{}, fmt.Errorf("failed to delete order,internal server error")
	}
	return &v1.DeleteOrderReply{}, nil

}

func (os *OdSrv) GetOd(ctx context.Context, req *v1.GetOrderReq) (*v1.GetOrderReply, error) {
	if req.Id < 1 {
		return &v1.GetOrderReply{Item: &v1.Item{}}, fmt.Errorf("order is wrong")
	}
	if len(req.UserUuid) == 0 {
		return &v1.GetOrderReply{Item: &v1.Item{}}, fmt.Errorf("user uuid is wrong")
	}
	item, err := os.odCase.GetOd(ctx, &v1.Item{Id: req.Id, UserUuid: req.UserUuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetOrderReply{Item: &v1.Item{}}, fmt.Errorf("no such a order")
		}
		os.logger.Errorf("get order error:%s", err.Error())
		return &v1.GetOrderReply{Item: &v1.Item{}}, fmt.Errorf("failed to get order,internal server error")
	}
	return &v1.GetOrderReply{Item: item}, nil

}

func (os *OdSrv) ListOd(ctx context.Context, req *v1.ListOrderReq) (*v1.ListOrderReply, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 5
	}
	odList, err := os.odCase.ListOd(ctx, req)
	if err != nil {
		os.logger.Errorf("list order error:%s", err.Error())
		return &v1.ListOrderReply{ItemList: make([]*v1.Item, 0)}, fmt.Errorf("failed to list order,internal server error")
	}
	return odList, nil
}

func (os *OdSrv) ListOdByCateId(ctx context.Context, req *v1.ListOrderByCateIdReq) (*v1.ListOrderByCateIdReply, error) {
	if req.CateId < 1 {
		return &v1.ListOrderByCateIdReply{ItemList: make([]*v1.Item, 0)}, fmt.Errorf("cate is is nil")
	}
	if len(req.UserUuid) == 0 {
		return &v1.ListOrderByCateIdReply{ItemList: make([]*v1.Item, 0)}, fmt.Errorf("user uuid is nil")
	}
	res, err := os.odCase.ListOdByCateId(ctx, req)
	if err != nil {
		if len(res.ItemList) == 0 {
			return &v1.ListOrderByCateIdReply{ItemList: make([]*v1.Item, 0), Cate: &v1.Cate{}}, fmt.Errorf("no such a cate")
		}
		os.logger.Errorf("get order by cate id error:%s", err.Error())
		return &v1.ListOrderByCateIdReply{ItemList: make([]*v1.Item, 0), Cate: &v1.Cate{}}, fmt.Errorf("no such a cate")
	}
	return res, nil
}

func (os *OdSrv) CreateCate(ctx context.Context, req *v1.CreateCateReq) (*v1.CreateCateReply, error) {
	if req.Cate.AddressId < 1 {
		return &v1.CreateCateReply{Cate: &v1.Cate{}}, fmt.Errorf("address is nil")
	}
	if len(req.Cate.UserUuid) == 0 {
		return &v1.CreateCateReply{Cate: &v1.Cate{}}, fmt.Errorf("user uuid is nil")
	}
	if req.Cate.Price < 1 {
		return &v1.CreateCateReply{Cate: &v1.Cate{}}, fmt.Errorf("price error")
	}
	res, err := os.caCase.CreateCate(ctx, &v1.Cate{AddressId: req.Cate.AddressId, UserUuid: req.Cate.UserUuid})
	if err != nil {
		os.logger.Errorf("create cate error:%s", err.Error())
		return &v1.CreateCateReply{Cate: &v1.Cate{}}, fmt.Errorf("failed to create cate,internal server error")
	}
	return &v1.CreateCateReply{Cate: res}, nil
}

func (os *OdSrv) ListOdForShopper(ctx context.Context, req *v1.ListOdForSReq) (*v1.ListOdForSpReply, error) {
	if req.ShopId < 1 {
		return &v1.ListOdForSpReply{OdList: make([]*v1.Item, 0)}, fmt.Errorf("shop id error")
	}
	if req.Status != 0 && req.Status != 1 && req.Status != 2 && req.Status != 3 {
		return &v1.ListOdForSpReply{OdList: make([]*v1.Item, 0)}, fmt.Errorf("status error")
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 5
	}
	res, err := os.odCase.ListOdForSp(ctx, req)
	if err != nil {
		os.logger.Errorf("list order for shop error:%s", err.Error())
		return &v1.ListOdForSpReply{OdList: make([]*v1.Item, 0)}, fmt.Errorf("internal server error")
	}
	return res, nil
}
func (os *OdSrv) PayOd(ctx context.Context, req *v1.PayOdReq) (*v1.PayOdReply, error) {
	if req.Id < 1 {
		return &v1.PayOdReply{}, fmt.Errorf("order id error")
	}
	if len(req.UserUuid) == 0 {
		return &v1.PayOdReply{}, fmt.Errorf("user id is nil")
	}
	od, err := os.GetOd(ctx, &v1.GetOrderReq{Id: int64(req.Id), UserUuid: req.UserUuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.PayOdReply{}, fmt.Errorf("no such a order")
		}
		os.logger.Errorf("get order error:%s", err.Error())
		return &v1.PayOdReply{}, fmt.Errorf("internal server error")
	}
	if od.Item.Status != 1 {
		return &v1.PayOdReply{}, fmt.Errorf("status not right")
	}
	res, err := os.odCase.PayOd(ctx, od.Item)
	if err != nil {
		os.logger.Errorf("pay order error:%s", err.Error())
		return &v1.PayOdReply{}, fmt.Errorf("internal server error")
	}
	return res, nil
}
