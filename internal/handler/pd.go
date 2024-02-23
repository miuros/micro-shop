package handler

import (
	"context"
	"log"
	pdv1 "micro-shop/api/pd/v1"
	"micro-shop/internal/model"
	"micro-shop/internal/service"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PdHl struct {
	pdSrv  *service.PdSrv
	logger *log.Logger
	ctx    context.Context
}

func NewPdHl(logger *log.Logger, pdSrv *service.PdSrv) *PdHl {
	logger.SetOutput(os.Stdout)

	logger.SetPrefix("handler/product")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	return &PdHl{
		logger: logger,
		pdSrv:  pdSrv,
		ctx:    context.Background(),
	}
}

func (ph *PdHl) GetPd(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
		return
	}
	var req = &pdv1.GetPdReq{Id: int64(id)}
	res, err := ph.pdSrv.GetPd(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformPd(res.Pd)))
}

func (ph *PdHl) CreatePd(c *gin.Context) {
	var pd = new(model.Product)
	if err := c.ShouldBindJSON(pd); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}

	if len(pd.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "product name is nil", nil))
		return
	}
	if pd.SellPrice < 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "sell price is wrong", nil))
		return
	}
	if pd.CategoryId < 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "category id is wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var reqSp = &pdv1.GetSpByUuidReq{}
	var err error
	role := c.GetString("roleName")
	var shopUuid string
	if role == "shopper" {
		reqSp.UserUuid = userUuid
		res, err := ph.pdSrv.GetSpByUuid(ph.ctx, reqSp)
		if err != nil {
			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			ph.logger.Println(err)
			return
		}
		pd.ShopId = uint64(res.Sp.Id)
		shopUuid = userUuid
	} else {
		if pd.ShopId < 1 {
			c.JSON(http.StatusOK, util.GetResponse(402, "shop id error", nil))
			return
		}
		res, err := ph.pdSrv.GetSp(ph.ctx, &pdv1.GetShopReq{Id: int64(pd.ShopId)})
		if err != nil {
			ph.logger.Println(err)
			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			return
		}
		shopUuid = res.Sp.UserUuid
	}
	var req = &pdv1.Product{
		Name:        pd.Name,
		ImageUrl:    pd.ImageUrl,
		OriginPrice: float32(pd.OriginPrice),
		SellPrice:   float32(pd.SellPrice),
		Desc:        pd.Desc,
		Tags:        pd.Tags,
		ShopId:      int64(pd.ShopId),
		CategoryId:  int64(pd.CategoryId),
		Extra:       pd.Extra,
	}
	reply, err := ph.pdSrv.CreatePd(ph.ctx, &pdv1.CreatePdReq{UserUuid: shopUuid, Pd: req})
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ph.logger.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformPd(reply.Pd)))
}

func (ph *PdHl) UpdatePd(c *gin.Context) {
	var pd = new(model.Product)
	if err := c.ShouldBindJSON(pd); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if pd.Id < 1 {

		c.JSON(http.StatusOK, util.GetResponse(402, "id is nil", nil))
	}

	if len(pd.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "product name is nil", nil))
		return
	}
	if pd.SellPrice < 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "sell price is wrong", nil))
		return
	}
	if pd.CategoryId < 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "category id is wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var res = new(pdv1.GetShopReply)
	var err error
    res, err = ph.pdSrv.GetSp(ph.ctx, &pdv1.GetShopReq{Id: int64(pd.ShopId)})
		if err != nil {
			ph.logger.Println(err)
			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			return
		}
	if role == "shopper" {
		
		if userUuid != res.Sp.UserUuid {
			c.JSON(http.StatusOK, util.GetResponse(403, "user uuid is wrong", nil))
			return
		}
	}


	var req = &pdv1.Product{
		Id:          int64(pd.Id),
		Name:        pd.Name,
		ImageUrl:    pd.ImageUrl,
		OriginPrice: float32(pd.OriginPrice),
		SellPrice:   float32(pd.SellPrice),
		ShopId:      int64(pd.ShopId),
		Desc:        pd.Desc,
		Tags:        pd.Tags,
		CategoryId:  int64(pd.CategoryId),
		Extra:       pd.Extra,
	}
	reply, err := ph.pdSrv.UpdatePd(ph.ctx, &pdv1.UpdatePdReq{UserUuid: res.Sp.UserUuid, Pd: req})
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformPd(reply.Pd)))

}

func (ph *PdHl) DeletePd(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id error", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &pdv1.DeletePdReq{}
	req.Id = int64(id)
	if role == "shopper" {
		req.UserUuid = userUuid
	} else {
		req.UserUuid = reqUuid
	}
	_, err = ph.pdSrv.DeletePd(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (ph *PdHl) ListPd(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	name := c.Query("name")
	var req = &pdv1.ListPdReq{Limit: int64(limit), Page: int64(page)}
	if len(name) != 0 {
		req.Name = name
	}
	res, err := ph.pdSrv.ListPd(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.PdList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var pdList = make([]*model.Product, len(res.PdList))
	for i, item := range res.PdList {
		pdList[i] = transformPd(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", pdList))
}

func (ph *PdHl) ListPdForSp(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	shopId, err := strconv.Atoi(c.Query("shopId"))
	if err != nil || shopId < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request shop id error", nil))
		return
	}
	var req = &pdv1.ListForSpReq{ShopId: uint64(shopId), Page: uint64(page), Limit: uint64(limit)}
	res, err := ph.pdSrv.ListPdForSp(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.PdList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var pdList = make([]*model.Product, len(res.PdList))
	for i, item := range res.PdList {
		pdList[i] = transformPd(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", pdList))
}

func (ph *PdHl) ListPdByCgId(c *gin.Context) {
	cgId, err := strconv.Atoi(c.Query("categoryId"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request category id is wrong", nil))
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	var req = &pdv1.ListPdByCiReq{Page: uint64(page), Limit: uint64(limit), CategoryId: uint64(cgId)}
	res, err := ph.pdSrv.ListPdByCg(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.PdList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var pdList = make([]*model.Product, len(res.PdList))
	for i, item := range res.PdList {
		pdList[i] = transformPd(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", pdList))
}

func (ph *PdHl) GetSp(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	var req = &pdv1.GetShopReq{Id: int64(id)}
	res, err := ph.pdSrv.GetSp(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSp(res.Sp)))
}

func (ph *PdHl) CreateSp(c *gin.Context) {
	var sp = new(model.Shop)
	if err := c.ShouldBindJSON(sp); err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	if len(sp.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request name is nil", nil))
		return
	}
	if len(sp.UserUuid) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
		return
	}
	if len(sp.Address) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request address  is nil", nil))
		return
	}
	var req = &pdv1.CreateShopReq{Sp: &pdv1.Shop{
		Name:     sp.Name,
		Address:  sp.Address,
		ImageUrl: sp.ImageUrl,
	}}
	var err error
	if role == "shopper" {
		if userUuid != sp.UserUuid {
			c.JSON(http.StatusOK, util.GetResponse(402, "user uuid is wrong", nil))
			return
		}
		req.Sp.UserUuid = userUuid
	} else {
		req.Sp.UserUuid = sp.UserUuid
	}
	res, err := ph.pdSrv.CreateSp(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSp(res.Sp)))
}

func (ph *PdHl) UpdateSp(c *gin.Context) {
	var sp = new(model.Shop)
	if err := c.ShouldBindJSON(sp); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if sp.Id < 1 {

	}
	if len(sp.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request name is nil", nil))
		return
	}
	if len(sp.UserUuid) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
		return
	}
	if len(sp.Address) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request address  is nil", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	if role == "shopper" {
		sp.UserUuid = userUuid
	}
	var req = &pdv1.UpdateShopReq{Sp: &pdv1.Shop{
		Id:       int64(sp.Id),
		Name:     sp.Name,
		ImageUrl: sp.ImageUrl,
		Address:  sp.Address,
		UserUuid: sp.UserUuid,
	}}
	res, err := ph.pdSrv.UpdateSp(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSp(res.Sp)))
}

func (ph *PdHl) GetSpByUuid(c *gin.Context) {
	userUuid := c.Query("userUuid")
	var req = &pdv1.GetSpByUuidReq{UserUuid: userUuid}
	res, err := ph.pdSrv.GetSpByUuid(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ph.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformSp(res.Sp)))
}

func (ph *PdHl) DeleteSp(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &pdv1.DeleteShopReq{Id: int64(id)}
	if role == "shopper" {
		req.UserUuid = userUuid
	} else {
		if len(reqUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = reqUuid
	}
	_, err = ph.pdSrv.DeleteSp(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (ph *PdHl) ListSp(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	name := c.Query("name")
	var req = &pdv1.ListShopReq{Page: int64(page), Limit: int64(limit)}
	if len(name) != 0 {
		req.Name = name
	}
	res, err := ph.pdSrv.ListSp(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.SpList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var spList = make([]*model.Shop, len(res.SpList))
	for i, item := range res.SpList {
		spList[i] = transformSp(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", spList))
}

func (ph *PdHl) CreateCt(c *gin.Context) {
	var ct = new(model.Cart)
	if err := c.ShouldBindJSON(ct); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if ct.ProductId < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request product is wrong", nil))
		return
	}
	if ct.Num < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request number is wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	if role == "purchaser" || role == "shopper" {
		ct.UserUuid = userUuid
		res, err := ph.pdSrv.GetPd(ph.ctx, &pdv1.GetPdReq{Id: int64(ct.ProductId)})
		if err != nil {
			ph.logger.Println(err)
			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			return
		}
		ct.ShopId = uint64(res.Pd.ShopId)
		ct.ProductName = res.Pd.Name
		ct.ImageUrl = res.Pd.ImageUrl
		ct.Price = float64(float32(ct.Num) * res.Pd.SellPrice)
		sp, err := ph.pdSrv.GetSp(ph.ctx, &pdv1.GetShopReq{Id: int64(res.Pd.ShopId)})
		if err != nil {
			ph.logger.Println(err)

			c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
			return
		}
		ct.ShopName = sp.Sp.Name
	} else {
		if len(ct.UserUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
	}
	var req = &pdv1.CreateCartReq{C: &pdv1.Cart{
		ProductId:   int64(ct.ProductId),
		ProductName: ct.ProductName,
		ImageUrl:    ct.ImageUrl,
		UserUuid:    ct.UserUuid,
		ShopId:      int64(ct.ShopId),
		ShopName:    ct.ShopName,
		Num:         int64(ct.Num),
		Price:       float32(ct.Price),
	}}
	reply, err := ph.pdSrv.CreateCt(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCt(reply.C)))
}

func (ph *PdHl) UpdateCt(c *gin.Context) {
	var ct = new(model.Cart)
	if err := c.ShouldBindJSON(ct); err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if ct.Id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id wrong", nil))
		return
	}
	if ct.Num < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request number is wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &pdv1.UpdateCartReq{C: &pdv1.Cart{
		Id:  int64(ct.Id),
		Num: int64(ct.Num),
	}}
	if role == "purchaser" || role == "shopper" {
		req.C.UserUuid = userUuid
	} else {
		req.C.UserUuid = ct.UserUuid
	}
	reply, err := ph.pdSrv.UpdateCt(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCt(reply.C)))
}

func (ph *PdHl) DeleteCt(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &pdv1.DeleteCartReq{Id: int64(id)}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	} else {
		if len(reqUuid) == 0 {
			c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
			return
		}
		req.UserUuid = reqUuid
	}
	_, err = ph.pdSrv.DeleteCt(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ph.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (ph *PdHl) GetCt(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	roleName := c.GetString("roleName")
	var req = &pdv1.GetCartReq{Id: int64(id)}
	if roleName == "shopper" || roleName == "purchaser" {
		req.UserUuid = userUuid
	} else {
		req.UserUuid = reqUuid
	}
	res, err := ph.pdSrv.GetCart(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCt(res.C)))
}

func (ph *PdHl) ListCt(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	reqUuid := c.Query("userUuid")
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &pdv1.ListCartReq{Page: int64(page), Limit: int64(limit)}
	if role == "purchaser" || role == "shopper" {
		req.UserUuid = userUuid
	} else {
		req.UserUuid = reqUuid
	}
	res, err := ph.pdSrv.ListCt(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.CartList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var ctList = make([]*model.Cart, len(res.CartList))
	for i, item := range res.CartList {
		ctList[i] = transformCt(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ctList))
}

func (ph *PdHl) CreateBn(c *gin.Context) {
	var bn = new(model.Banner)
	if err := c.ShouldBindJSON(bn); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if len(bn.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request name is nil", nil))
		return
	}
	if len(bn.ImageUrl) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request image url is nil", nil))
		return
	}
	if len(bn.RedirectUrl) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request redirect url is nil", nil))
		return
	}
	var req = &pdv1.CreateBnReq{Bn: &pdv1.Banner{Name: bn.Name, ImageUrl: bn.ImageUrl, RedirectUrl: bn.RedirectUrl}}
	res, err := ph.pdSrv.CreateBn(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformBn(res.Bn)))
}

func (ph *PdHl) UpdateBn(c *gin.Context) {
	var bn = new(model.Banner)
	if err := c.ShouldBindJSON(bn); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if bn.Id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	var req = &pdv1.UpdateBnReq{Bn: &pdv1.Banner{Id: int64(bn.Id), Name: bn.Name, ImageUrl: bn.ImageUrl, RedirectUrl: bn.RedirectUrl}}
	res, err := ph.pdSrv.UpdateBn(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformBn(res.Bn)))
}

func (ph *PdHl) DeleteBn(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	var req = &pdv1.DeleteBnReq{Id: int64(id)}
	_, err = ph.pdSrv.DeleteBn(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (ph *PdHl) GetBn(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	var req = &pdv1.GetBnReq{Id: int64(id)}
	res, err := ph.pdSrv.GetBn(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformBn(res.Bn)))
}

func (ph *PdHl) ListBn(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	var req = &pdv1.ListBnReq{Page: int64(page), Limit: int64(limit)}
	res, err := ph.pdSrv.ListBn(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.BnList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var bnList = make([]*model.Banner, len(res.BnList))
	for i, item := range res.BnList {
		bnList[i] = transformBn(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", bnList))
}

func (ph *PdHl) CreateCg(c *gin.Context) {
	var cg = new(model.Category)
	if err := c.ShouldBindJSON(cg); err != nil {
		ph.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if len(cg.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request name is nil", nil))
		return
	}
	var req = &pdv1.CreateCgReq{Cg: &pdv1.Category{Name: cg.Name}}
	res, err := ph.pdSrv.CreateCg(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ph.logger.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCg(res.Cg)))
}

func (ph *PdHl) UpdateCg(c *gin.Context) {
	var cg = new(model.Category)
	if err := c.ShouldBindJSON(cg); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		ph.logger.Printf("bind json error:%s", err.Error())
		return
	}
	if cg.Id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	if len(cg.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request name is nil", nil))
		return
	}
	var req = &pdv1.UpdateCgReq{Cg: &pdv1.Category{Id: cg.Id, Name: cg.Name}}
	res, err := ph.pdSrv.UpdateCg(ph.ctx, req)
	if err != nil {
		ph.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCg(res.Cg)))
}

func (ph *PdHl) DeleteCg(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	var req = &pdv1.DeleteCgReq{Id: uint64(id)}
	_, err = ph.pdSrv.DeleteCg(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (ph *PdHl) GetCg(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	var req = &pdv1.GetCgReq{Id: uint64(id)}
	res, err := ph.pdSrv.GetCg(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ph.logger.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCg(res.Cg)))
}

func (ph *PdHl) ListCg(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	var req = &pdv1.ListCgReq{Limit: uint64(limit), Page: uint64(page)}
	res, err := ph.pdSrv.ListCg(ph.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ph.logger.Println(err.Error())
		return
	}
	if len(res.CgList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
		return
	}
	var cgList = make([]*model.Category, len(res.CgList))
	for i, item := range res.CgList {
		cgList[i] = transformCg(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", cgList))
}

func transformCg(res *pdv1.Category) *model.Category {
	return &model.Category{
		Id:   res.Id,
		Name: res.Name,
	}
}

func transformBn(res *pdv1.Banner) *model.Banner {
	return &model.Banner{
		Id:          uint64(res.Id),
		Name:        res.Name,
		ImageUrl:    res.ImageUrl,
		RedirectUrl: res.RedirectUrl,
	}
}

func transformCt(res *pdv1.Cart) *model.Cart {
	return &model.Cart{
		Id:          uint64(res.Id),
		ProductId:   uint64(res.ProductId),
		ProductName: res.ProductName,
		ShopId:      uint64(res.ShopId),
		ShopName:    res.ShopName,
		Price:       float64(res.Price),
		Num:         uint64(res.Num),
		UserUuid:    res.UserUuid,
		ImageUrl:    res.ImageUrl,
	}
}

func transformSp(res *pdv1.Shop) *model.Shop {
	return &model.Shop{
		Id:       uint64(res.Id),
		Name:     res.Name,
		ImageUrl: res.ImageUrl,
		UserUuid: res.UserUuid,
		Address:  res.Address,
		CreateAt: res.CreateAt,
		DeleteAt: res.DeleteAt,
		IsDelete: uint64(res.IsDeleted),
	}
}

func transformPd(res *pdv1.Product) *model.Product {
	return &model.Product{
		Id:          uint64(res.Id),
		Name:        res.Name,
		ShopId:      uint64(res.ShopId),
		CategoryId:  uint64(res.CategoryId),
		SellPrice:   float64(res.SellPrice),
		OriginPrice: float64(res.OriginPrice),
		Tags:        res.Tags,
		Desc:        res.Desc,
		Extra:       res.Extra,

		ImageUrl:  res.ImageUrl,
		CreateAt:  res.CreateAt,
		DeleteAt:  res.DeleteAt,
		IsDeleted: uint64(res.IsDeleted),
	}
}
