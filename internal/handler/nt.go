package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	ntv1 "micro-shop/api/notice/v1"
	"micro-shop/internal/model"
	"micro-shop/internal/service"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"strconv"
)

type NtHl struct {
	ntSrv  *service.NtSrv
	logger *log.Logger
	ctx    context.Context
}

func NewNtHl(logger *log.Logger, ntSrv *service.NtSrv) *NtHl {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("handler/notice")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	return &NtHl{
		logger: logger,
		ntSrv:  ntSrv,
		ctx:    context.Background(),
	}
}

func (nh *NtHl) CreateNt(c *gin.Context) {
	var nt = new(model.Notice)
	if err := c.ShouldBindJSON(nt); err != nil {
		nh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if nt.Type != "chat" && nt.Type != "notice" {
		c.JSON(http.StatusOK, util.GetResponse(402, "type is wrong", nil))
		return
	}
	if len(nt.Content) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "content is nil", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	userName := c.GetString("userName")
	role := c.GetString("roleName")
	var req = &ntv1.CreateNtReq{N: &ntv1.Notice{Content: nt.Content, Type: nt.Type, ToUserUuid: nt.UserUuid, Status: int64(nt.Status)}}
	if role == "shopper" || role == "purchaser" {
		req.N.UserUuid = userUuid
		req.N.UserName = userName
	} else {
		req.N.UserUuid = nt.UserUuid
		req.N.UserName = nt.UserName
	}
	_, err := nh.ntSrv.CreateNt(nh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		nh.logger.Println(err)
		return
	}

	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (nh *NtHl) UpdateNt(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id wrong", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &ntv1.UpdateStatusReq{Id: int64(id), UserUuid: userUuid}
	_, err = nh.ntSrv.UpdateNtStatus(nh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		nh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (nh *NtHl) DeleteNt(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "wrong id ", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &ntv1.DeleteNtReq{Id: int64(id), UserUuid: userUuid}
	_, err = nh.ntSrv.DeleteNt(nh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		nh.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
}

func (nh *NtHl) GetNt(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
		return
	}
	typ := c.Query("type")
	if len(typ) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "type is nil", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &ntv1.GetNtReq{Id: uint64(id), UserUuid: userUuid, Type: typ}
	res, err := nh.ntSrv.GetNt(nh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		nh.logger.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformNt(res.Nt)))
}

func (nh *NtHl) ListNt(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	typ := c.Query("type")
	if len(typ) == 0 || (typ != "notice" && typ != "chat") {
		c.JSON(http.StatusOK, util.GetResponse(402, "type error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	var req = &ntv1.ListNtReq{Page: uint64(page), Limit: uint64(limit), UserUuid: userUuid, Type: typ}
	res, err := nh.ntSrv.ListNt(nh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		nh.logger.Println(err)
		return
	}
	if res.Num == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var ntList = make([]*model.Notice, len(res.NtList))
	for i, item := range res.NtList {
		ntList[i] = transformNt(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", map[string]interface{}{
		"page":  page,
		"limit": limit,
		"data":  ntList,
		"num":   res.Num,
	}))

}

func transformNt(res *ntv1.Notice) *model.Notice {
	return &model.Notice{
		Id:         uint64(res.Id),
		UserUuid:   res.UserUuid,
		ToUserUuid: res.ToUserUuid,
		UserName:   res.UserName,
		Type:       res.Type,
		Content:    res.Content,
		CreateAt:   res.CreateAt,
		Status:     uint64(res.Status),
		IsDeleted:  uint64(res.IsDeleted),
	}
}
