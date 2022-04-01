package model

type Notice struct {
	Id         uint64 `json:"id"`
	UserUuid   string `json:"userUuid"`
	UserName   string `json:"userName"`
	ShopId     uint64 `json:"shopId"`
	ShopName   string `json:"shopName"`
	ToUserUuid string `json:"toUserUuid"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	CreateAt   string `json:"createAt"`
	Status     uint64 `json:"status"`
	IsDeleted  uint64 `json:"isDeleted"`
}
