package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	odv1 "user-srv/api/od/v1"
	"user-srv/internal/biz"
	"user-srv/internal/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/streadway/amqp"
)

type userRepo struct {
	data   *Data
	logger *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) *userRepo {
	ur := &userRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data")),
	}
	if err := ur.Notice(); err != nil {
		panic(err)
	}
	return ur
}

func (ur *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {

	err := ur.data.maria.Create(u).Error
	if err != nil {
		return nil, err
	}
	ur.cache(ctx, createRedisKey(u.Uuid), u)
	return u, nil

}

func (u *userRepo) cache(ctx context.Context, key string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err == nil {
		u.data.redis.Set(ctx, key, jsonData, time.Duration(60*time.Second))
	}
}

func (ur *userRepo) UpdateUser(ctx context.Context, u *biz.User) (*biz.User, error) {

	err := ur.data.maria.Model(&biz.User{}).Where("uuid=?", u.Uuid).Update(u).Error
	if err != nil {
		return nil, err
	}

	ur.cache(ctx, createRedisKey(u.Uuid), u)

	return u, nil
}

func (ur *userRepo) DeleteUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	err := ur.data.maria.Model(&biz.User{}).Where("uuid=?", u.Uuid).Update(&biz.User{DeleteAt: time.Now().String(), IsDeleted: 1}).Error
	if err != nil {
		return nil, err
	}
	ur.data.redis.Del(ctx, createRedisKey(u.Uuid))
	return u, err

}

func (ur *userRepo) GetUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	byteData, err := ur.data.redis.Get(ctx, createRedisKey(u.Uuid)).Bytes()
	var data = new(biz.User)
	if err == nil {
		err = unCache(byteData, data)
		if err == nil {
			return data, nil
		}
	}
	if err != nil {
		err = ur.data.maria.Model(&biz.User{}).Where("uuid=?", u.Uuid).First(data).Error
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (u *userRepo) ListUser(ctx context.Context, qu *biz.QueryUser) ([]biz.User, error) {

	var userList []biz.User
	if len(qu.Name) != 0 {
		u.data.maria.Model(&biz.User{}).Where("name like ?", fmt.Sprintf("%%%s%%", qu.Name))
	}
	err := u.data.maria.Offset(qu.Offset).Limit(qu.Limit).Find(&userList).Error
	if err != nil {
		return nil, err
	}
	return userList, nil

}

func (ur *userRepo) FindUserByName(ctx context.Context, u *biz.User) (*biz.User, error) {
	err := ur.data.maria.Model(&biz.User{}).Where("name=?", u.Name).First(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

var _ biz.UserRepo = (*userRepo)(nil)

func createRedisKey(key string) string {
	return "uuid_" + key
}

func unCache(byteData []byte, data *biz.User) error {
	return json.Unmarshal(byteData, data)
}

func (ur *userRepo) Notice() error {
	c, err := ur.data.rabbit.Channel()
	if err != nil {
		return err
	}
	msg, err := c.Consume("mail_queue", "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go ur.mail(msg)
	return nil
}

type Order struct {
	Id        int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductId int64   `protobuf:"varint,2,opt,name=productId,proto3" json:"productId,omitempty"`
	Number    int64   `protobuf:"varint,3,opt,name=number,proto3" json:"number,omitempty"`
	PayType   int64   `protobuf:"varint,4,opt,name=payType,proto3" json:"payType,omitempty"`
	Status    int64   `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
	AddressId int64   `protobuf:"varint,6,opt,name=addressId,proto3" json:"addressId,omitempty"`
	IsDeleted int64   `protobuf:"varint,7,opt,name=isDeleted,proto3" json:"isDeleted,omitempty"`
	Price     float32 `protobuf:"fixed32,8,opt,name=price,proto3" json:"price,omitempty"`
	UserUuid  string  `protobuf:"bytes,9,opt,name=userUuid,proto3" json:"userUuid,omitempty"`
	CateId    int64   `protobuf:"varint,10,opt,name=cateId,proto3" json:"cateId,omitempty"`
	PayTime   string  `protobuf:"bytes,11,opt,name=payTime,proto3" json:"payTime,omitempty"`
	CreateAt  string  `protobuf:"bytes,12,opt,name=createAt,proto3" json:"createAt,omitempty"`
	UpdateAt  string  `protobuf:"bytes,13,opt,name=updateAt,proto3" json:"updateAt,omitempty"`
	DeletedAt string  `protobuf:"bytes,14,opt,name=deletedAt,proto3" json:"deletedAt,omitempty"`
}

func (or *userRepo) GetOd(ctx context.Context, id int64, uuid string) (*Order, error) {
	var req = &odv1.GetOrderReq{Id: id, UserUuid: uuid}
	res, err := or.data.odSrv.GetOd(context.Background(), req)
	if err != nil {
		return nil, err

	}
	return &Order{
		Id:        res.Item.Id,
		UserUuid:  res.Item.UserUuid,
		Status:    res.Item.Status,
		IsDeleted: res.Item.IsDeleted,
	}, nil
}

func (ur *userRepo) mail(msg <-chan amqp.Delivery) {
	ctx := context.Background()
	for m := range msg {
		var od = new(Order)
		err := json.Unmarshal(m.Body, od)
		if err != nil {
			ur.logger.Errorf("json unmarshal error:%s", err.Error())
			continue
		}
		od, err = ur.GetOd(ctx, od.Id, od.UserUuid)
		if err != nil {
			ur.logger.Errorf("get order error:%s", err.Error())
			continue
		}
		if od.Status == 2 || od.IsDeleted == 1 {
			continue
		}
		user, err := ur.GetUser(ctx, &biz.User{Uuid: od.UserUuid})
		if err != nil {
			ur.logger.Errorf("get user error:%s", err.Error())
			continue
		}
		err = util.SendMail(user.Mail, user.Name, ur.data.mail, ur.data.host, ur.data.passwd)
		if err != nil {
			ur.logger.Errorf("send mail error:%s", err.Error())
		}
	}
}
