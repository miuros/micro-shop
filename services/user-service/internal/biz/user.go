package biz

import (
	"context"
	"time"
	v1 "user-srv/api/user/v1"
	"user-srv/internal/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
	DeleteUser(context.Context, *User) (*User, error)
	GetUser(context.Context, *User) (*User, error)
	ListUser(context.Context, *QueryUser) ([]User, error)
	FindUserByName(context.Context, *User) (*User, error)
}

type User struct {
	Uuid      string `json:"uuid" gorm:"column:uuid;primaryKey;index"`
	Name      string `json:"name" gorm:"column:name;unique"`
	Password  string `json:"password" gorm:"column:password;not null"`
	Mobile    string `json:"mobile" gorm:"column:mobile"`
	Mail      string `json:"mail" gorm:"column:mail"`
	RoleName  string `json:"roleName" gorm:"column:role_name;not null"`
	CreateAt  string `json:"createAt" gorm:"column:create_at"`
	DeleteAt  string `json:"deleteAt" gorm:"column:delete_at"`
	IsDeleted int    `json:"isDeleted" gorm:"column:is_deleted"`
}

type QueryUser struct {
	Offset uint
	Limit  uint
	Name   string
}

type UserUseCase struct {
	v1.UnsafeUserServer
	repo   UserRepo
	logger *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz")),
	}
}

func (uu *UserUseCase) SearchUserByName(ctx context.Context, req *v1.SearchUserByNameRequest) ([]*v1.UserInfo, error) {
	var qu = &QueryUser{
		Name:   req.Name,
		Offset: uint(req.Limit * (req.Page - 1)),
		Limit:  uint(req.Limit),
	}
	res, err := uu.repo.ListUser(ctx, qu)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return make([]*v1.UserInfo, 0), nil
	}
	var uList = make([]*v1.UserInfo, len(res))
	for i, item := range res {
		uList[i] = CreateReply(&item)
	}
	return uList, nil
}

func (uu *UserUseCase) FindUserByName(ctx context.Context, req *v1.UserInfo) (*v1.UserInfo, error) {
	var u = &User{
		Name: req.Name,
	}
	res, err := uu.repo.FindUserByName(ctx, u)
	if err != nil {
		return nil, err
	}
	return CreateReply(res), nil
}

func (u *UserUseCase) CreateUser(ctx context.Context, request *v1.UserInfo) (*v1.UserInfo, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	encrypted, err := util.EncryptPassword(request.Password)
	if err != nil {
		return nil, err
	}
	var data = &User{
		Uuid:      id.String(),
		Name:      request.Name,
		Password:  encrypted,
		Mobile:    request.Mobile,
		Mail:      request.Mail,
		CreateAt:  time.Now().Format("2006-01-02:15-04"),
		RoleName:  request.RoleName,
		IsDeleted: 0,
	}
	res, err := u.repo.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}

	return CreateReply(res), nil
}

func (ur *UserUseCase) UpdateUser(ctx context.Context, req *v1.UserInfo) (*v1.UserInfo, error) {
	encodePw, err := util.EncryptPassword(req.Password)
	if err != nil {
		return &v1.UserInfo{}, err
	}
	var data = &User{
		Password: encodePw,
		Uuid:     req.Uuid,
		Name:     req.Name,
		Mobile:   req.Mobile,
		RoleName: req.RoleName,
		Mail:     req.Mail,
	}
	res, err := ur.repo.UpdateUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return CreateReply(res), nil
}

func (ur *UserUseCase) DeleteUser(ctx context.Context, req *v1.UserInfo) (*v1.UserInfo, error) {
	var u = &User{
		Uuid:     req.Uuid,
		DeleteAt: time.Now().Format("2006-01-02:15-04"),
	}
	res, err := ur.repo.DeleteUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return CreateReply(res), nil
}

func (ur *UserUseCase) GetUser(ctx context.Context, req *v1.UserInfo) (*v1.UserInfo, error) {
	var u = &User{Uuid: req.Uuid}
	res, err := ur.repo.GetUser(ctx, u)
	if err != nil {
		return nil, err
	}
	res.Password = ""
	return CreateReply(res), nil
}

func (u *UserUseCase) ListUser(ctx context.Context, request *v1.ListUserRequest) ([]*v1.UserInfo, error) {
	var qu = &QueryUser{
		Offset: uint(request.Limit * (request.Page - 1)),
		Limit:  uint(request.Limit),
		Name:   request.Name,
	}
	res, err := u.repo.ListUser(ctx, qu)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return make([]*v1.UserInfo, 0), nil
	}
	var uList = make([]*v1.UserInfo, len(res))
	for i, item := range res {
		uList[i] = CreateReply(&item)
	}
	return uList, nil
}
func CreateReply(data *User) *v1.UserInfo {
	return &v1.UserInfo{
		Uuid:      data.Uuid,
		Name:      data.Name,
		RoleName:  data.RoleName,
		Mobile:    data.Mobile,
		Mail:      data.Mail,
		Password:  data.Password,
		IsDeleted: int64(data.IsDeleted),
	}
}
