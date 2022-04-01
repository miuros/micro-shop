package service

import (
	"context"
	"google.golang.org/grpc"
	"log"
	userv1 "micro-shop/api/user/v1"
	"os"
)

type UserSrv struct {
	uc     userv1.UserClient
	logger *log.Logger
}

func NewUserSrv(logger *log.Logger, conn grpc.ClientConnInterface) *UserSrv {
	uc := userv1.NewUserClient(conn)
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("data/user:")
	return &UserSrv{
		uc:     uc,
		logger: logger,
	}
}

func (us *UserSrv) Register(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserReply, error) {
	res, err := us.uc.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginReply, error) {
	res, err := us.uc.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserReply, error) {
	res, err := us.uc.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserReply, error) {
	res, err := us.uc.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserReply, error) {
	res, err := us.uc.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) ListUser(ctx context.Context, req *userv1.ListUserRequest) (*userv1.ListUserReply, error) {
	res, err := us.uc.ListUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) CreateAddress(ctx context.Context, req *userv1.CreateAddressRequest) (*userv1.CreateAddressReply, error) {
	res, err := us.uc.CreateAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) UpdateAddress(ctx context.Context, req *userv1.UpdateAddressRequest) (*userv1.UpdateAddressReply, error) {
	res, err := us.uc.UpdateAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv) DeleteAddress(ctx context.Context, req *userv1.DeleteAddressRequest) (*userv1.DeleteAddressRely, error) {
	res, err := us.uc.DeleteAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserSrv)GetAddress(ctx context.Context,req *userv1.GetAddressRequest)(*userv1.GetAddressReply,error){
	res,err:=us.uc.GetAddress(ctx,req)
	if err !=nil{
		return nil, err
	}
	return res,nil
}

func (us *UserSrv) ListAddress(ctx context.Context, req *userv1.ListAddressRequest) (*userv1.ListAddressReply, error) {
	res, err := us.uc.ListAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
