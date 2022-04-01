package model

type UserForRegister struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Passwd    string `json:"passwd"`
	Mobile    string `json:"mobile"`
	Code      string `json:"code"`
	Mail      string `json:"mail"`
	RoleName  string `json:"roleName"`
	CreateAt  string `json:"createAt"`
	DeleteAt  string `json:"deleteAt"`
	IsDeleted int64  `json:"isDeleted"`
}

type User struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Mail      string `json:"mail"`
	RoleName  string `json:"roleName"`
	CreateAt  string `json:"createAt"`
	DeleteAt  string `json:"deleteAt"`
	IsDeleted int64  `json:"isDeleted"`
}

type AddressInfo struct {
	Id       uint64 `json:"id"`
	UserUuid string `json:"userUuid"`
	Mobile   string `json:"mobile"`
	Address  string `json:"address"`
	Alias    string `json:"alias"`
}
