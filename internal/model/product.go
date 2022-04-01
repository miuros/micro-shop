package model


type Category struct{
	Id uint64 `json:"id"`
	Name string `json:"name"`
}

type Shop struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	UserUuid string `json:"userUuid"`
	Address  string `json:"address"`
	CreateAt string `json:"createAt"`
	DeleteAt string `json:"deleteAt"`
	IsDelete uint64 `json:"isDelete"`
}

type Cart struct {
	Id          uint64  `json:"id"`
	ImageUrl    string  `json:"imageUrl"`
	UserUuid    string  `json:"userUuid"`
	ProductId   uint64  `json:"productId"`
	ProductName string  `json:"productName"`
	ShopId      uint64  `json:"shopId"`
	ShopName    string  `json:"shopName"`
	Num         uint64  `json:"num"`
	Price       float64 `json:"price"`
}

type Product struct {
	Id          uint64  `json:"id"`
	Name        string  `json:"name"`
	CategoryId  uint64  `json:"categoryId"`
	OriginPrice float64 `json:"originPrice"`
	SellPrice   float64 `json:"sellPrice"`
	ImageUrl    string  `json:"imageUrl"`
	Desc        string  `json:"desc"`
	Tags        string  `json:"tags"`
	ShopId      uint64  `json:"shopId"`
	Extra       string  `json:"extra"`
	CreateAt    string  `json:"createAt"`
	DeleteAt    string  `json:"deleteAt"`
	IsDeleted   uint64  `json:"isDeleted"`
}

type Banner struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	RedirectUrl string `json:"redirectUrl"`
}
