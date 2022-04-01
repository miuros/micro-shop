package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"user-srv/api/user/v1"
	"user-srv/internal/pkg/util"
)

func (u *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	if len(in.User.Mobile) != 11 {
		return &v1.CreateUserReply{User: &v1.UserInfo{}}, fmt.Errorf("wrong format of mobile")
	}
	if len(in.User.Mail) == 0 {
		return &v1.CreateUserReply{User: &v1.UserInfo{}}, fmt.Errorf("mail is nil")
	}
	reply, err := u.userCase.CreateUser(ctx, in.User)
	if err != nil {
		u.logger.Errorf("create user error:%s", err.Error())
		return &v1.CreateUserReply{}, fmt.Errorf("failed to create user,internal server error")
	}
	return &v1.CreateUserReply{User: reply}, nil
}

func (us *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	if len(in.User.Uuid) == 0 {
		return &v1.UpdateUserReply{User: &v1.UserInfo{}}, fmt.Errorf("uuid is nil")
	}
	_, err := us.userCase.GetUser(ctx, in.User)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.UpdateUserReply{User: &v1.UserInfo{}}, fmt.Errorf("no such user")
		}
		us.logger.Errorf("get user error:%s", err.Error())
		return &v1.UpdateUserReply{User: &v1.UserInfo{}}, fmt.Errorf("failed to update user,internal server error")
	}
	reply, err := us.userCase.UpdateUser(ctx, in.User)
	if err != nil {
		us.logger.Errorf("update user error:%s", err.Error())
		return &v1.UpdateUserReply{}, fmt.Errorf("failed to update user,internal server error")
	}
	return &v1.UpdateUserReply{User: reply}, nil
}

func (us *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	_, err := us.userCase.DeleteUser(ctx, &v1.UserInfo{Uuid: in.Uuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.DeleteUserReply{}, fmt.Errorf("user not found")
		}
		us.logger.Errorf("delete user error:%s", err.Error())
		return &v1.DeleteUserReply{}, fmt.Errorf("failed to delete user,internal server error")
	}
	return &v1.DeleteUserReply{}, nil
}

func (u *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	reply, err := u.userCase.GetUser(ctx, &v1.UserInfo{Uuid: in.Uuid})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.GetUserReply{User: &v1.UserInfo{}}, fmt.Errorf("user not found")
		}
		u.logger.Errorf("get user error:%s", err.Error())
		return &v1.GetUserReply{}, fmt.Errorf("failed to get user,internal server error")
	}
	return &v1.GetUserReply{User: reply}, nil
}

func (u *UserService) ListUser(ctx context.Context, request *v1.ListUserRequest) (*v1.ListUserReply, error) {

	if request.Limit < 5 {
		request.Limit = 5
	}
	if request.Page < 1 {
		request.Page = 1
	}
	reply, err := u.userCase.ListUser(ctx, request)
	if err != nil {
		u.logger.Errorf("list user error:%s", err.Error())
		return &v1.ListUserReply{UserList: make([]*v1.UserInfo, 0)}, fmt.Errorf("failed to get user list,internal server error")
	}
	return &v1.ListUserReply{UserList: reply}, nil
}

func (u *UserService) SearchUserByName(ctx context.Context, in *v1.SearchUserByNameRequest) (*v1.SearchUserByNameReply, error) {
	if len(in.Name) == 0 {
		return &v1.SearchUserByNameReply{UserList: make([]*v1.UserInfo, 0)}, fmt.Errorf("name is nil")
	}
	reply, err := u.userCase.SearchUserByName(ctx, in)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &v1.SearchUserByNameReply{UserList: make([]*v1.UserInfo, 0)}, fmt.Errorf("user not found")
		}
		u.logger.Errorf("find by name error:%s", err.Error())
		return &v1.SearchUserByNameReply{UserList: make([]*v1.UserInfo, 0)}, fmt.Errorf("failed to find user ,internal server error")
	}
	return &v1.SearchUserByNameReply{UserList: reply}, nil
}

func (us *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	if len(in.Name) == 0 {
		return &v1.LoginReply{User: &v1.UserInfo{}}, fmt.Errorf("name is nil")
	}
	user, err := us.userCase.FindUserByName(ctx, &v1.UserInfo{Name: in.Name})
	if err != nil {
		if err==gorm.ErrRecordNotFound{
			return &v1.LoginReply{}, fmt.Errorf("no such a user")
		}
		us.logger.Errorf("login error:%s", err.Error())
		return &v1.LoginReply{}, fmt.Errorf("internal server error")
	}
	ok := util.VerifyPasswd(user.Password, in.Password)
	if !ok {
		return &v1.LoginReply{}, fmt.Errorf("password is wrong")
	}
	user.Password = ""
	return &v1.LoginReply{User: user}, nil
}

func (u *UserService) Logout(ctx context.Context, in *v1.LogoutRequest) (*v1.LogoutReply, error) {
	return &v1.LogoutReply{}, nil
}
