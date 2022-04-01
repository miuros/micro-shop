package model

type Message struct {
	Id         uint   `json:"id"`
	UserName   string `json:"userName"`
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
