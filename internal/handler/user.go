package handler

import (
	"context"
	"log"
	userv1 "micro-shop/api/user/v1"
	"micro-shop/internal/model"
	"micro-shop/internal/service"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	logger  *log.Logger
	userSrv *service.UserSrv
	ctx     context.Context
	mh      *MlHl
}

func NewUserHandler(logger *log.Logger, userSrv *service.UserSrv) *UserHandler {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("handler/user")
	logger.SetFlags(log.Ltime | log.Lshortfile)
	uh := &UserHandler{
		logger:  logger,
		userSrv: userSrv,
		ctx:     context.Background(),
		mh:      NewMlHl(),
	}
	return uh
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var u = new(model.UserForRegister)
	if err := c.ShouldBindJSON(u); err != nil {
		uh.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	if u.RoleName != "shopper" && u.RoleName != "purchaser" {
		c.JSON(http.StatusOK, util.GetResponse(402, "request role name error", nil))
		return
	}
	if len(u.Mail) == 0 || len(u.Name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "name or mail is nil", nil))
		return
	}
	role := c.GetString("roleName")
	if role == "shopper" || role == "purchaser" {
		c.JSON(http.StatusOK, util.GetResponse(403, "request error", nil))
		return
	}
	var req = &userv1.CreateUserRequest{User: &userv1.UserInfo{
		Name:     u.Name,
		Mail:     u.Mail,
		Password: u.Passwd,
		Mobile:   u.Mobile,
		RoleName: u.RoleName,
	}}
	res, err := uh.userSrv.Register(uh.ctx, req)
	if err != nil {
		uh.logger.Printf("register error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", map[string]interface{}{
		"user": transformUser(res.User),
	}))
}

func (uh *UserHandler) Login(c *gin.Context) {

	name := c.PostForm("name")
	passwd := c.PostForm("passwd")
	if len(name) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(401, "name is nil", nil))
		return
	}
	if len(passwd) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(401, "password is nil", nil))
		return
	}
	var req = &userv1.LoginRequest{Name: name, Password: passwd}
	res, err := uh.userSrv.Login(uh.ctx, req)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, err.Error(), nil))
		return
	}
	tokenStr, err := util.ReleaseToken(res.User.Name, res.User.Uuid, res.User.RoleName)
	if err != nil {
		uh.logger.Printf("release token error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, "internal server error", nil))
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", map[string]interface{}{
		"token": tokenStr,
		"user":  transformUser(res.User),
	}))
}

func (uh *UserHandler) Mail(c *gin.Context) {
	mail := c.Query("mail")
	if len(mail) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request mail is nil", nil))
		return
	}
	err := uh.mh.Mail(mail)
	if err != nil {
		uh.logger.Printf("mail error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, "internal server error", nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
}

func (uh *UserHandler) Register(c *gin.Context) {
	var u = new(model.UserForRegister)
	if err := c.ShouldBindJSON(u); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request error", nil))
		uh.logger.Printf("bind json error:%s", err.Error())
		return
	}
	if len(u.Name) < 1 {
		c.JSON(http.StatusOK, util.GetResponse(402, "name is nil", nil))
		return
	}
	if len(u.Passwd) < 6 {
		c.JSON(http.StatusOK, util.GetResponse(402, "length of password must greater than 6", nil))
		return
	}
	if len(u.Mobile) != 11 {
		c.JSON(http.StatusOK, util.GetResponse(402, "length of mobile must be 11", nil))
		return
	}
	if len(u.Mail) < 5 {
		c.JSON(http.StatusOK, util.GetResponse(402, "mail error", nil))
		return
	}
	if len(u.Code) != 6 {
		c.JSON(http.StatusOK, util.GetResponse(402, "verify code error", nil))
		return
	}
	code, ok := uh.mh.GetCode(u.Mail)
	if !ok {
		c.JSON(http.StatusOK, util.GetResponse(402, "verify code has out of date", nil))
		return
	}
	if code != u.Code {
		c.JSON(http.StatusOK, util.GetResponse(402, "code is not right", nil))
		return
	}
	if u.RoleName != "shopper" && u.RoleName != "purchaser" {
		c.JSON(http.StatusOK, util.GetResponse(402, "role name is wrong,nil", nil))
	}
	var req = &userv1.CreateUserRequest{User: &userv1.UserInfo{
		Name:     u.Name,
		Mail:     u.Mail,
		Password: u.Passwd,
		Mobile:   u.Mobile,
		RoleName: u.RoleName,
	}}
	res, err := uh.userSrv.Register(uh.ctx, req)
	if err != nil {
		uh.logger.Printf("register error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", map[string]interface{}{
		"user": transformUser(res.User),
	}))
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var u = new(model.User)
	if err := c.ShouldBindJSON(u); err != nil {
		uh.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	roleName := c.GetString("roleName")
	var res = new(userv1.UpdateUserReply)
	var err error
	if roleName == "purchaser" || roleName == "shopper" {
		var req = &userv1.UpdateUserRequest{User: &userv1.UserInfo{
			Uuid:   userUuid,
			Name:   u.Name,
			Mobile: u.Mobile,
			Mail:   u.Mail,
		}}

		res, err = uh.userSrv.UpdateUser(uh.ctx, req)
	} else {
		var req = &userv1.UpdateUserRequest{User: &userv1.UserInfo{
			Uuid:   u.Uuid,
			Mail:   u.Mail,
			Mobile: u.Mobile,
			Name:   u.Name,
		}}
		res, err = uh.userSrv.UpdateUser(uh.ctx, req)
	}

	if err != nil {
		uh.logger.Printf(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformUser(res.User)))

}

func (uh *UserHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	uuid := c.GetString("userUuid")
	reqUuid := c.Query("userUuid")
	roleName := c.GetString("roleName")
	var err error
	if roleName == "purchaser" || roleName == "shopper" {
		var req = &userv1.DeleteUserRequest{Uuid: uuid}
		_, err = uh.userSrv.DeleteUser(uh.ctx, req)
	} else {
		var req = &userv1.DeleteUserRequest{Uuid: reqUuid}
		_, err = uh.userSrv.DeleteUser(uh.ctx, req)
	}

	if err != nil {
		uh.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	reqUuid := c.Query("uuid")
	role := c.GetString("roleName")
	userUUid := c.GetString("userUuid")
	var res *userv1.GetUserReply
	var err error
	var req = &userv1.GetUserRequest{Uuid: reqUuid}
	res, err = uh.userSrv.GetUser(uh.ctx, req)

	if err != nil {
		uh.logger.Println(err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
	}
	if reqUuid == userUUid {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformUser(res.User)))
		return
	}
	if role == "shopper" || role == "purchaser" {
		res.User.Password = ""
		res.User.Mail = ""
		res.User.Mobile = ""
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformUser(res.User)))

}

func (uh *UserHandler) ListUser(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	name := c.Query("name")
	role := c.GetString("role")
	if role == "shopper" || role == "purchaser" {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", ""))
		return
	}
	var page int
	var limit int
	var err error
	if page, err = strconv.Atoi(pageStr); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request page error", nil))
		return
	}
	if limit, err = strconv.Atoi(limitStr); err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request limit error", nil))
		return
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}
	var req = &userv1.ListUserRequest{Limit: int64(limit), Page: int64(page)}
	if len(name) != 0 {
		req.Name = name
	}
	res, err := uh.userSrv.ListUser(uh.ctx, req)
	if err != nil {
		uh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.UserList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
		return
	}
	var uList = make([]*model.User, len(res.UserList))
	for i, item := range res.UserList {
		uList[i] = transformUser(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", uList))
}

func (uh *UserHandler) GetAddress(c *gin.Context) {
	idStr := c.Query("id")
	reqUuid := c.Query("userUuid")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(402, "request id is wrong", nil))
		return
	}
	if len(reqUuid) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(402, "request user id is nil", nil))
		return
	}
	var req = new(userv1.GetAddressRequest)
	req.Id = uint64(id)
	req.UserUuid = reqUuid
	res, err := uh.userSrv.GetAddress(uh.ctx, req)
	if err != nil {
		uh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformAddress(res.Address)))
}

func (uh *UserHandler) CreateAr(c *gin.Context) {
	var ad = new(model.AddressInfo)
	if err := c.ShouldBindJSON(ad); err != nil {
		uh.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &userv1.CreateAddressRequest{Address: &userv1.AddressInfo{
		Mobile:  ad.Mobile,
		Alias:   ad.Alias,
		Address: ad.Address,
	}}
	if role == "purchaser" || role == "shopper" {
		req.Address.UserUuid = userUuid
	} else {
		req.Address.UserUuid = ad.UserUuid
	}
	res, err := uh.userSrv.CreateAddress(uh.ctx, req)
	if err != nil {
		uh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformAddress(res.Address)))
}

func (uh *UserHandler) UpdateAr(c *gin.Context) {
	var ar = new(model.AddressInfo)
	if err := c.ShouldBindJSON(ar); err != nil {
		uh.logger.Printf("bind json error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")
	var req = &userv1.UpdateAddressRequest{Address: &userv1.AddressInfo{
		Id:       ar.Id,
		UserUuid: ar.UserUuid,
		Address:  ar.Address,
		Alias:    ar.Alias,
		Mobile:   ar.Mobile,
	}}
	if role == "shopper" || role == "purchaser" {
		req.Address.UserUuid = userUuid
	}
	res, err := uh.userSrv.UpdateAddress(uh.ctx, req)
	if err != nil {
		uh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", transformAddress(res.Address)))
}

func (uh *UserHandler) DeleteAr(c *gin.Context) {
	var idStr = c.Query("id")
	var reqUuid = c.Query("userUuid")
	id, err := strconv.Atoi(idStr)
	role := c.GetString("roleName")
	userUuid := c.GetString("userUuid")
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(401, "request error", nil))
		return
	}
	var req = &userv1.DeleteAddressRequest{Id: uint64(id), UserUuid: reqUuid}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	}
	_, err = uh.userSrv.DeleteAddress(uh.ctx, req)
	if err != nil {
		uh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
}

func (uh *UserHandler) ListAr(c *gin.Context) {
	reqUuid := c.Query("userUuid")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 5
	}
	userUuid := c.GetString("userUuid")
	role := c.GetString("roleName")

	var req = &userv1.ListAddressRequest{UserUuid: reqUuid, Page: uint64(page), Limit: uint64(limit)}
	if role == "shopper" || role == "purchaser" {
		req.UserUuid = userUuid
	}
	res, err := uh.userSrv.ListAddress(uh.ctx, req)
	if err != nil {
		uh.logger.Println(err)
		c.JSON(http.StatusOK, util.GetResponse(501, err.Error(), nil))
		return
	}
	if len(res.AddressList) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", nil))
		return
	}
	var arList = make([]*model.AddressInfo, len(res.AddressList))
	for i, item := range res.AddressList {
		arList[i] = transformAddress(item)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", arList))
}

func transformAddress(req *userv1.AddressInfo) *model.AddressInfo {
	return &model.AddressInfo{
		Id:       req.Id,
		UserUuid: req.UserUuid,
		Address:  req.Address,
		Mobile:   req.Mobile,
		Alias:    req.Alias,
	}
}

func transformUser(req *userv1.UserInfo) *model.User {
	return &model.User{
		Uuid:      req.Uuid,
		Name:      req.Name,
		Mobile:    req.Mobile,
		Mail:      req.Mail,
		RoleName:  req.RoleName,
		CreateAt:  req.CreateAt,
		DeleteAt:  req.DeleteAt,
		IsDeleted: int64(req.IsDeleted),
	}
}
