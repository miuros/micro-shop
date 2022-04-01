package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
	Id         uint   `json:"id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	UserUuid   string `json:"userUuid"`
	UserType   string `json:"UserType"`
	ShopId     uint64 `json:"shopId"`
	ShopName   string `json:"shopName"`
	ToUserUuid string `json:"toUserUuid"`
	Status     uint   `json:"status"`
	CreateAt   string `json:"createAt"`
	IsDeleted  uint   `json:"isDeleted"`
}

func main() {
	var header = http.Header{}
	header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDg0NzU5ODIsImlzcyI6Im1pY3JvLXNob3AiLCJVc2VybmFtZSI6Im5hbmEtcHVyY2hhc2VyIiwiUm9sZU5hbWUiOiJwdXJjaGFzZXIiLCJVc2VyVXVpZCI6Ijc3NGYzY2ZmLThiNDYtNGY1MC1iMGM1LTM0MzlhNDM4MGM0NSJ9.zAnPTVX3d9Io7HbeUYMsxLogrRVQctnXAfqRmxcn72c")
	conn, _, err := websocket.DefaultDialer.Dial("ws://miuros.fun:10000/api/chat/v1/chat", header)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.WriteJSON(&Message{Content: "test"})
}
