package service

import (
	"context"
	"fmt"
	"user-srv/api/user/v1"
)

func (us *UserService) CreateAddress(ctx context.Context, in *v1.CreateAddressRequest) (*v1.CreateAddressReply, error) {
	if in.Address.UserUuid == "" {
		return &v1.CreateAddressReply{}, fmt.Errorf("user uuid is nil")
	}
	if in.Address.Address == "" {
		return &v1.CreateAddressReply{}, fmt.Errorf("address is nil")
	}
	if in.Address.Mobile == "" {

	}
	reply, err := us.addressCase.CreateAddress(ctx, in.Address)
	if err != nil {
		us.logger.Errorf("create address error:%s", err.Error())
		return &v1.CreateAddressReply{}, err
	}
	return &v1.CreateAddressReply{Address: reply}, nil
}

func (us *UserService) UpdateAddress(ctx context.Context, in *v1.UpdateAddressRequest) (*v1.UpdateAddressReply, error) {
	if in.Address.UserUuid == "" {
		return &v1.UpdateAddressReply{}, fmt.Errorf("user uuid is nil")
	}

	reply, err := us.addressCase.UpdateAddress(ctx, in.Address)
	if err != nil {
		us.logger.Errorf("update address error:%s", err.Error())
		return &v1.UpdateAddressReply{}, err
	}
	return &v1.UpdateAddressReply{Address: reply}, nil
}

func (us *UserService) GetAddress(ctx context.Context, in *v1.GetAddressRequest) (*v1.GetAddressReply, error) {
	reply, err := us.addressCase.GetAddress(ctx, in)
	if err != nil {
		us.logger.Errorf("get address error:%s", err.Error())
		return &v1.GetAddressReply{}, fmt.Errorf("internal server error")
	}
	return reply, nil
}

func (us *UserService) ListAddress(ctx context.Context, in *v1.ListAddressRequest) (*v1.ListAddressReply, error) {
	reply, err := us.addressCase.ListAddress(ctx, in)
	if err != nil {
		us.logger.Errorf("list address error:%s", err.Error())
		return &v1.ListAddressReply{AddressList: make([]*v1.AddressInfo, 0)}, fmt.Errorf("internal server error")
	}
	return reply, nil
}

func (us *UserService) DeleteAddress(ctx context.Context, in *v1.DeleteAddressRequest) (*v1.DeleteAddressRely, error) {
	reply, err := us.addressCase.DeleteAddress(ctx, in)
	if err != nil {
		us.logger.Errorf("delete address error:%s", err.Error())
		return &v1.DeleteAddressRely{}, err
	}
	return reply, nil
}
