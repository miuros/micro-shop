package model

type Comment struct {
	Id         uint64 `json:"id"`
	ProductId  uint64 `json:"productId"`
	UserUuid   string `json:"userUuid"`
	ToUserUuid string `json:"toUserUuid"`
	Content    string `json:"content"`
	CreateAt   string `json:"createAt"`
	UpdateAt   string `json:"UpdateAt"`
	DeleteAt   string `json:"deleteAt"`
	IsDeleted  uint64 `json:"isDeleted"`
}
