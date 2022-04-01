package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	ntv1 "micro-shop/api/notice/v1"
	pdv1 "micro-shop/api/pd/v1"
	"micro-shop/internal/model"
	"micro-shop/internal/service"
	"micro-shop/internal/util"
	"net/http"
	"os"
	"sync"
	"time"
)

type WsClient struct {
	upgrader *websocket.Upgrader
	ntSrv    *service.NtSrv
	pdSrv    *service.PdSrv
	logger   *log.Logger
	locker   *sync.Mutex
	msgChan  chan *model.Message
	clients  map[string]*websocket.Conn
}

func NewWsClient(logger *log.Logger, ntSrv *service.NtSrv, pdSrv *service.PdSrv) *WsClient {
	logger.SetOutput(os.Stdout)
	logger.SetFlags(log.Ltime | log.Lshortfile)
	log.SetPrefix("handler/websocket:")
	cli := &WsClient{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		msgChan: make(chan *model.Message, 50),
		logger:  logger,
		ntSrv:   ntSrv,
		pdSrv:   pdSrv,
		clients: make(map[string]*websocket.Conn),
		locker:  &sync.Mutex{},
	}
	cli.logger.SetPrefix("websocket:")
	go cli.Record()
	return cli
}

func (wc *WsClient) Upgrade(c *gin.Context) {
	tokenStr := c.Query("token")
	if len(tokenStr) < 7 {
		c.JSON(http.StatusOK, util.GetResponse(402, "login first plz", nil))
		return
	}
	claim, err := util.ParseToken("Bearer " + tokenStr)
	if err != nil {
		c.JSON(http.StatusOK, util.GetResponse(501, "internal server error", nil))
		wc.logger.Printf("parse token error:%s", err.Error())
		return
	}
	if claim.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusOK, util.GetResponse(201, "login again plz", nil))
		return
	}
	userUuid := claim.UserUuid
	roleName := claim.RoleName
	userName := claim.Username
	var sp *model.Shop
	if roleName == "shopper" {
		res, err := wc.pdSrv.GetSpByUuid(context.Background(), &pdv1.GetSpByUuidReq{UserUuid: userUuid})
		if err != nil {

			c.JSON(http.StatusOK, util.GetResponse(501, "internal server error", nil))
			return
		}
		sp = &model.Shop{
			Id:       uint64(res.Sp.Id),
			UserUuid: res.Sp.UserUuid,
			CreateAt: res.Sp.CreateAt,
			Name:     res.Sp.Name,
		}
	}
	conn, err := wc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wc.logger.Printf("upgrade protocol http to websocket error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(500, "failed to connect", nil))
		return
	}
	defer conn.Close()
	wc.locker.Lock()
	wc.clients[userUuid] = conn
	wc.locker.Unlock()
	defer func() {
		if err := recover(); err != nil {
			wc.logger.Printf("panic:%v", err)
		}
		wc.locker.Lock()
		delete(wc.clients, userUuid)
		wc.locker.Unlock()
	}()
	for {
		var msg = new(model.Message)
		err := conn.ReadJSON(msg)

		if err != nil {
			if err == websocket.ErrCloseSent || err == websocket.ErrReadLimit || err == websocket.ErrBadHandshake {
				wc.locker.Lock()
				delete(wc.clients, userUuid)
				wc.locker.Unlock()
				wc.logger.Printf("error :%s", err.Error())
				return
			}
			continue
		}
		if len(msg.ToUserUuid) == 0 {
			continue
		}
		msg.Type = "chat"
		if len(msg.Content) == 0 {
			continue
		}
		msg.CreateAt = time.Now().Format("2006-01-02:15-04")
		msg.UserUuid = userUuid
		if roleName == "shopper" {
			msg.UserType = "shop"
			msg.ShopId = sp.Id
			msg.UserName = sp.Name
			msg.ShopName = sp.Name
		} else {
			msg.UserType = "user"
			msg.UserName = userName
		}

		msg.Status = 1
		err = conn.WriteJSON(msg)
		if err != nil {
			if err == websocket.ErrCloseSent || err == websocket.ErrReadLimit || err == websocket.ErrBadHandshake {
				return
			}
			wc.logger.Printf("write json error:%s", err.Error())
			continue
		}
		err = wc.SendMessage(msg)
		if err != nil {
			wc.logger.Printf("send message error:%s", err.Error())
			continue
		}
		wc.msgChan <- msg

	}
}

func (wc *WsClient) SendMessage(msg *model.Message) error {
	if toConn, ok := wc.clients[msg.ToUserUuid]; ok {
		err := toConn.WriteJSON(msg)
		if err != nil {
			if err == websocket.ErrCloseSent || err == websocket.ErrReadLimit || err == websocket.ErrBadHandshake {
				wc.locker.Lock()
				delete(wc.clients, msg.ToUserUuid)
				wc.locker.Unlock()
				return nil
			}
			return err
		}
	}
	return nil
}

func (wc *WsClient) Record() {
	ctx := context.Background()
	for item := range wc.msgChan {
		var req = &ntv1.CreateNtReq{N: &ntv1.Notice{
			Content:    item.Content,
			Type:       item.Type,
			UserUuid:   item.UserUuid,
			UserName:   item.UserName,
			ToUserUuid: item.ToUserUuid,
		}}
		_, err := wc.ntSrv.CreateNt(ctx, req)
		if err != nil {
			wc.logger.Println(err.Error())
		}

	}
}
