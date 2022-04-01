package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	cmv1 "micro-shop/api/cm/v1"
	"micro-shop/internal/model"
	"micro-shop/internal/service"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"strconv"
)

type CmHl struct {
	cmSrv  *service.CmSrv
	logger *log.Logger
	ctx    context.Context
}

func NewCmHl(logger *log.Logger, cmSrv *service.CmSrv) *CmHl {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("handler/comment:")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	return &CmHl{
		logger: logger,
		ctx:    context.Background(),
		cmSrv:  cmSrv,
	}
}

func (ch *CmHl) CreateCm(c *gin.Context) {
	var cm = new(model.Comment)
	if err := c.ShouldBind(cm); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		ch.logger.Printf("bind json error:%s", err.Error())
		return
	}
	if cm.ProductId < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request product is nil", nil))
		return
	}
	if len(cm.Content) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request content is nil", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &cmv1.CreateCmReq{Cm: &cmv1.Comment{
		ProductId:  int64(cm.ProductId),
		Content:    cm.Content,
		UserUuid:   userUuid,
		ToUserUuid: cm.ToUserUuid,
	}}
	res, err := ch.cmSrv.CreateCm(ch.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ch.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCm(res.Cm)))
}

func (ch *CmHl) DeleteCm(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &cmv1.DeleteCmReq{Id: int64(id), UserUuid: userUuid}
	_, err = ch.cmSrv.DeleteCm(ch.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ch.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "OK", ""))
}

func (ch *CmHl) UpdateCm(c *gin.Context) {
	var cm = new(model.Comment)
	if err := c.ShouldBindJSON(cm); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		ch.logger.Printf("bind json error:%s", err.Error())
		return
	}
	if cm.ProductId < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request product is nil", nil))
		return
	}
	if cm.Id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is nil", nil))
		return
	}

	if len(cm.Content) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request content is nil", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &cmv1.UpdateCmReq{Cm: &cmv1.Comment{Id: int64(cm.Id), UserUuid: userUuid, Content: cm.Content, ProductId: int64(cm.ProductId)}}
	res, err := ch.cmSrv.UpdateCm(ch.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ch.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCm(res.Cm)))
}

func (ch *CmHl) ListCm(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("productid"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 1
	}
	var req = &cmv1.ListCmReq{ProductId: int64(id), Page: int64(page), Limit: int64(limit)}
	res, err := ch.cmSrv.ListCm(ch.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ch.logger.Println(err)
		return
	}
	if len(res.CmList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "OK", ""))
		return
	}
	var cmList = make([]*model.Comment, len(res.CmList))
	for i, item := range res.CmList {
		cmList[i] = transformCm(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", cmList))
}

func (ch *CmHl) GetCm(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("commentid"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(401, "requet error", nil))
		return
	}
	var req = &cmv1.GetCmReq{Id: int64(id)}
	res, err := ch.cmSrv.GetCm(ch.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		ch.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformCm(res.Cm)))
}

func transformCm(res *cmv1.Comment) *model.Comment {
	return &model.Comment{
		Id:         uint64(res.Id),
		UserUuid:   res.UserUuid,
		ToUserUuid: res.ToUserUuid,
		Content:    res.Content,
		CreateAt:   res.CreateAt,
		ProductId:  uint64(res.ProductId),
		DeleteAt:   res.DeleteAt,
		IsDeleted:  uint64(res.IsDeleted),
	}
}
