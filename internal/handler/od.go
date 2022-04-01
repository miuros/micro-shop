package handler

import (
	"context"
	"fmt"
	"log"
	ntv1 "micro-shop/api/notice/v1"
	odv1 "micro-shop/api/od/v1"
	pdv1 "micro-shop/api/pd/v1"
	userv1 "micro-shop/api/user/v1"
	"micro-shop/internal/model"
	"micro-shop/internal/service"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type OdHl struct {
	odSrv  *service.OdSrv
	pdSrv  *service.PdSrv
	urSrv  *service.UserSrv
	ntSrv  *service.NtSrv
	logger *log.Logger
	ctx    context.Context
}

func NewOdHl(odSrv *service.OdSrv, pdSrv *service.PdSrv, urSrv *service.UserSrv, ntSrv *service.NtSrv, logger *log.Logger) *OdHl {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("handler/order")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	return &OdHl{
		odSrv:  odSrv,
		urSrv:  urSrv,
		ntSrv:  ntSrv,
		pdSrv:  pdSrv,
		logger: logger,
		ctx:    context.Background(),
	}
}

func (oh *OdHl) CreateOd(c *gin.Context) {
	var odList []model.Order
	userUuid := c.GetString("userUuid")

	if err := c.ShouldBindJSON(&odList); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		oh.logger.Printf("bind json error:%s", err.Error())
		return
	}
	if len(odList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request is nil", nil))
		return
	}
	var total float32
	userName := c.GetString("userName")

	var shopIdList = make([]int64, len(odList))
	for i, item := range odList {
		if item.ProductId < 1 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request product is wrong", nil))
			return
		}
		if item.Number < 1 {
			c.JSON(http.StatusOK, util.GetResponse(402, fmt.Sprintf("request %d number is wrong", item.ProductId), nil))
			return
		}
		if item.AddressId < 1 {
			res, err := oh.urSrv.ListAddress(oh.ctx, &userv1.ListAddressRequest{UserUuid: userUuid, Limit: 10})
			if err != nil {
				c.JSON(http.StatusOK, util.GetResponse(501, "address error", nil))
				return
			}
			if len(res.AddressList) == 0 {
				c.JSON(http.StatusOK, util.GetResponse(402, "add address info plz", nil))
				return
			}
			odList[i].AddressId = uint64(res.AddressList[i].Id)
		} else {
			_, err := oh.urSrv.GetAddress(oh.ctx, &userv1.GetAddressRequest{UserUuid: userUuid, Id: uint64(item.AddressId)})
			if err != nil {
				c.JSON(http.StatusOK, util.GetResponse(402, "address not found", nil))
				return
			}
		}
		pd, err := oh.pdSrv.GetPd(oh.ctx, &pdv1.GetPdReq{Id: int64(item.ProductId)})
		if err != nil {
			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			return
		}
		odList[i].Price = float64(pd.Pd.SellPrice * float32(odList[i].Number))
		shopIdList[i] = pd.Pd.ShopId
		total += pd.Pd.SellPrice
	}
	role := c.GetString("roleName")
	var cate = &odv1.CreateCateReq{Cate: &odv1.Cate{}}
	var reqUuid string
	cate.Cate.AddressId = int64(odList[0].AddressId)
	cate.Cate.Price = total
	if role == "shopper" || role == "purchaser" {
		cate.Cate.UserUuid = userUuid

		reqUuid = userUuid
	} else {
		cate.Cate.UserUuid = odList[0].UserUuid
		reqUuid = odList[0].UserUuid
	}

	res, err := oh.odSrv.CreateCate(oh.ctx, cate)
	if err != nil {
		oh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	var resList = make([]*model.Order, len(odList))
	var errList = make([]string, len(odList))

	var i int = 0
	for idx, item := range odList {
		var req = &odv1.CreateOrderReq{Item: &odv1.Item{
			ProductId: int64(item.ProductId),
			Number:    int64(item.Number),
			UserUuid:  reqUuid,
			CreateAt:  time.Now().Format("2006-01-02:15-04"),

			CateId:    int64(res.Cate.Id),
			AddressId: int64(cate.Cate.AddressId),
			Price:     float32(item.Price),
		}}
		res, err := oh.odSrv.CreateOd(oh.ctx, req)
		if err != nil {
			oh.logger.Println(err)
			errList = append(errList, fmt.Sprintf("order %d :%s", item.ProductId, err.Error()))
			continue
		}
		i++

		resList = append(resList, transformOd(res.Item))

		sp, err := oh.pdSrv.GetSp(oh.ctx, &pdv1.GetShopReq{Id: shopIdList[idx]})
		if err != nil {
			continue
		}
		oh.ntSrv.CreateNt(oh.ctx, &ntv1.CreateNtReq{N: &ntv1.Notice{UserUuid: reqUuid, Content: "order created", Type: "notice", ToUserUuid: sp.Sp.UserUuid, UserName: userName}})
	}
	if i == 0 {
		c.JSON(http.StatusOK, util.GetResponse(501, "internal servere error", nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", map[string]interface{}{
		"err":     strings.Join(errList, ";"),
		"success": resList,
	}))

}

func (oh *OdHl) UpdateOd(c *gin.Context) {
	var od = new(model.Order)
	if err := c.ShouldBindJSON(od); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		oh.logger.Printf("bind json error: %s", err.Error())
		return
	}
	if od.Id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id wrong", nil))
		return
	}
	if od.AddressId < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request address is wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &odv1.UpdateOrderReq{}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {
		if len(od.UserUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = od.UserUuid
	}
	res, err := oh.odSrv.UpdateOd(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformOd(res.Item)))
}

func (oh *OdHl) DeleteOd(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is nil", nil))
		return
	}
	if id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id wrong", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	userName := c.GetString("userName")
	reqName := c.Query("userName")
	role := c.GetString("roleName")
	var req = &odv1.DeleteOrderReq{Id: int64(id)}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {
		if len(reqUuid) == 0 || len(reqName) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id or name is nil", nil))
			return
		}
		userName = reqName
		req.UserUuid = reqUuid
	}
	_, err = oh.odSrv.DeleteOd(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
	oh.Notice(id, req.UserUuid, userName, fmt.Sprintf(" deleted order %d", id))
}

func (oh *OdHl) GetOd(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id wrong", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &odv1.GetOrderReq{Id: int64(id)}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {
		if len(userUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = reqUuid
	}
	res, err := oh.odSrv.GetOd(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformOd(res.Item)))
}

func (oh *OdHl) ListOdByCateId(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("cateId"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	reqUuid := c.Query("userUuid")
	var req = &odv1.ListOrderByCateIdReq{CateId: int64(id)}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {
		if len(reqUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = reqUuid
	}

	res, err := oh.odSrv.ListOdByCateId(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.ItemList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var odList = make([]*model.Order, len(res.ItemList))
	for i, item := range res.ItemList {
		odList[i] = transformOd(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", odList))
}

func (oh *OdHl) ListOd(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &odv1.ListOrderReq{Limit: int64(limit), Page: int64(page)}
	reqUuid := c.Query("userUuid")
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {

		req.UserUuid = reqUuid
	}
	res, err := oh.odSrv.ListOd(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}

	if len(res.ItemList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var odList = make([]*model.Order, len(res.ItemList))
	for i, item := range res.ItemList {
		odList[i] = transformOd(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", odList))
}

func (oh *OdHl) ListOdForSp(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "status is wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")

	var req = &odv1.ListOdForSReq{Limit: uint64(limit), Page: uint64(page), Status: uint64(status)}
	role := c.GetString("roleName")
	if role == "shopper" {
		res, err := oh.pdSrv.GetSpByUuid(oh.ctx, &pdv1.GetSpByUuidReq{UserUuid: userUuid})
		if err != nil {
			oh.logger.Println(err)
			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			return
		}
		req.ShopId = uint64(res.Sp.Id)
	} else {
		shopId, err := strconv.Atoi(c.Query("shopId"))
		if err != nil {
			c.JSON(http.StatusOK, util.GetResponse(402, "shop id format error", nil))
			return
		}
		req.ShopId = uint64(shopId)
	}
	reply, err := oh.odSrv.ListOdForSp(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	if len(reply.OdList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var odList = make([]*model.Order, len(reply.OdList))
	for i, item := range reply.OdList {
		odList[i] = transformOd(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", odList))
}

func (oh *OdHl) GetSc(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
		return
	}
	var req = &odv1.GetStockReq{ProductId: int64(id)}
	res, err := oh.odSrv.GetSc(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSc(res.Stock)))
}

func (oh *OdHl) CreateSc(c *gin.Context) {
	var sc = new(model.Stock)
	if err := c.ShouldBindJSON(sc); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		oh.logger.Println(err.Error())
		return
	}
	if sc.ProductId < 1 || sc.Storage < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request params error", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	var req = &odv1.CreateStockReq{Stock: &odv1.StockInfo{ProductId: int64(sc.ProductId), Storage: int64(sc.Storage)}}

	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	if role == "shopper" {
		req.Stock.UserUuid = userUuid
	} else {
		req.Stock.UserUuid = reqUuid
	}
	res, err := oh.odSrv.CreateSc(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSc(res.Stock)))
}

func (oh *OdHl) UpdateSc(c *gin.Context) {
	var sc = new(model.Stock)
	if err := c.ShouldBindJSON(sc); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		oh.logger.Println(err.Error())
		return
	}
	if sc.ProductId < 1 || sc.Storage < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request params error", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	var req = &odv1.UpdateStockReq{Stock: &odv1.StockInfo{ProductId: int64(sc.ProductId), Storage: int64(sc.Storage)}}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	if role == "shopper" {
		req.Stock.UserUuid = userUuid
	} else {
		req.Stock.UserUuid = reqUuid
	}
	res, err := oh.odSrv.UpdateSc(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSc(res.Stock)))

}

func (oh *OdHl) DeleteSc(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	reqUuid := c.Query("userUuid")
	var req = &odv1.DeleteStockReq{ProductId: int64(id)}
	if role == "shopper" {
		req.UserUuid = userUuid
	} else {
		if len(reqUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = reqUuid
	}
	_, err = oh.odSrv.DeleteSc(oh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		oh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (oh *OdHl) Notice(id int, userUuid, userName, msg string) {
	od, err := oh.odSrv.GetOd(oh.ctx, &odv1.GetOrderReq{Id: int64(id), UserUuid: userUuid})
	if err != nil {
		oh.logger.Println(err)
		return
	}
	res, err := oh.pdSrv.GetPd(oh.ctx, &pdv1.GetPdReq{Id: int64(od.Item.ProductId)})
	if err != nil {
		oh.logger.Println(err)
		return
	}
	spRes, err := oh.pdSrv.GetSp(oh.ctx, &pdv1.GetShopReq{Id: res.Pd.ShopId})
	if err != nil {
		oh.logger.Println(err)
		return
	}
	_, err = oh.ntSrv.CreateNt(oh.ctx, &ntv1.CreateNtReq{N: &ntv1.Notice{UserUuid: userUuid, ToUserUuid: spRes.Sp.UserUuid, UserName: userName, Content: msg, Type: "notice"}})
	if err != nil {
		oh.logger.Println(err)
	}
}

func (oh *OdHl) PayOd(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	userName := c.GetString("userName")
	var req = &odv1.PayOdReq{Id: uint64(id)}
	reqUuid := c.PostForm("userUuid")
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {
		if len(reqUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = reqUuid
	}

	_, err = oh.odSrv.PayOd(oh.ctx, req)
	if err != nil {

		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
	oh.Notice(id, req.UserUuid, userName, fmt.Sprintf(" paid order %d", id))
}

func transformSc(res *odv1.StockInfo) *model.Stock {
	return &model.Stock{
		Id:        uint64(res.Id),
		ProductId: uint64(res.ProductId),
		Storage:   uint64(res.Storage),
		Sale:      uint64(res.Sale),
	}
}

func transformOd(res *odv1.Item) *model.Order {
	return &model.Order{
		Id:        uint64(res.Id),
		ProductId: uint64(res.ProductId),
		Number:    uint64(res.Number),
		PayType:   uint64(res.PayType),
		PayTime:   res.PayTime,
		Status:    uint64(res.Status),
		AddressId: uint64(res.AddressId),
		IsDeleted: uint64(res.IsDeleted),
		Price:     float64(res.Price),
		UserUuid:  res.UserUuid,
		CateId:    uint64(res.CateId),
		CreateAt:  res.CreateAt,
		UpdateAt:  res.UpdateAt,
		DeleteAt:  res.DeletedAt,
	}
}
