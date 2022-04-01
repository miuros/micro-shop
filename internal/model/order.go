package model

type Cate struct {
	Id        uint64  `json:"id"`
	UserUuid  string  `json:"userUuid"`
	Status    uint64  `json:"status"`
	Price     float64 `json:"price"`
	AddressId uint64  `json:"addressID"`
	DeleteAt  string  `json:"deleteAt"`
	IsDeleted uint64  `json:"isDeleted"`
}

type Order struct {
	Id        uint64  `json:"id"`
	ProductId uint64  `json:"productId"`
	Number    uint64  `json:"number"`
	PayType   uint64  `json:"payType"`
	Status    uint64  `json:"status"`
	AddressId uint64  `json:"addressId"`
	IsDeleted uint64  `json:"isDeleted"`
	Price     float64 `json:"price"`
	UserUuid  string  `json:"userUuid"`
	CateId    uint64  `json:"cateId"`
	PayTime   string  `json:"payTime"`
	CreateAt  string  `json:"createAt"`
	UpdateAt  string  `json:"updateAt"`
	DeleteAt  string  `json:"deleteAt"`
}

type Stock struct {
	Id        uint64 `json:"id"`
	ProductId uint64 `json:"productId"`
	Storage   uint64 `json:"storage"`
	Sale      uint64 `json:"sale"`
	IsDeleted uint64 `json:"isDeleted"`
	UserUuid  string `json:"userUuid"`
}
