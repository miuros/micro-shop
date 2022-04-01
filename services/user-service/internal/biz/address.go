package biz

import (
	"context"
	"fmt"
	v1 "user-srv/api/user/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/gorm"
)

type AddressInfo struct {
	Id       int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserUuid string `json:"userUuid" gorm:"column:user_uuid;not null"`
	Address  string `json:"address" gorm:"column:address;not null"`
	Mobile   string `json:"mobile" gorm:"column:mobile;not null"`
	Alias    string `json:"alias" gorm:"column:alias;omitempty"`
}

type QueryAddress struct {
	UserUuid string
	Page     int
	Limit    int
}
type AddressRepo interface {
	CreateAddress(context.Context, *AddressInfo) (*AddressInfo, error)
	UpdateAddress(context.Context, *AddressInfo) (*AddressInfo, error)
	GetAddress(context.Context, *AddressInfo) (*AddressInfo, error)
	DeleteAddress(context.Context, *AddressInfo) error
	ListAddress(context.Context, *QueryAddress) ([]AddressInfo, error)
}
type AddressUseCase struct {
	repo   AddressRepo
	logger *log.Helper
}

func NewAddressUseCase(repo AddressRepo, logger log.Logger) *AddressUseCase {
	return &AddressUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz")),
	}
}

func (a *AddressUseCase) CreateAddress(ctx context.Context, in *v1.AddressInfo) (*v1.AddressInfo, error) {
	var ad = &AddressInfo{
		UserUuid: in.UserUuid,
		Address:  in.Address,
		Mobile:   in.Mobile,
		Alias:    in.Alias,
	}
	ad, err := a.repo.CreateAddress(ctx, ad)
	if err != nil {
		return nil, err
	}
	return createAddress(ad), nil
}

func (a *AddressUseCase) UpdateAddress(ctx context.Context, in *v1.AddressInfo) (*v1.AddressInfo, error) {
	info := &AddressInfo{
		Id: int(in.Id),
	}
	info, err := a.repo.GetAddress(ctx, info)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no such a address info")
		} else {
			a.logger.Errorf("failed to get address info:%s", err.Error())
			return nil, err
		}
	}
	if info.UserUuid != in.UserUuid {
		return nil, fmt.Errorf("uuid is wrong")
	}
	info = &AddressInfo{
		Id:      int(in.Id),
		Address: in.Address,
		Mobile:  in.Mobile,
		Alias:   in.Alias,
	}
	info, err = a.repo.UpdateAddress(ctx, info)
	if err != nil {
		return nil, err
	}
	return createAddress(info), nil
}

func (a *AddressUseCase) GetAddress(ctx context.Context, in *v1.GetAddressRequest) (*v1.GetAddressReply, error) {
	var info = &AddressInfo{
		Id:       int(in.Id),
		UserUuid: in.UserUuid,
	}
	info, err := a.repo.GetAddress(ctx, info)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no such a address info")
		} else {
			a.logger.Errorf("failed to get address info:%s", err.Error())
			return nil, err
		}
	}
	return &v1.GetAddressReply{Address: createAddress(info)}, nil
}

func (a *AddressUseCase) ListAddress(ctx context.Context, in *v1.ListAddressRequest) (*v1.ListAddressReply, error) {
	var query = &QueryAddress{
		Limit:    int(in.Limit),
		Page:     int(in.Page),
		UserUuid: in.UserUuid,
	}
	list, err := a.repo.ListAddress(ctx, query)
	if err != nil {
		a.logger.Errorf("failed to get address list:%s", err.Error())
		return nil, err
	}
	var replyList = new(v1.ListAddressReply)
	replyList.AddressList = make([]*v1.AddressInfo, len(list))
	for i, item := range list {
		replyList.AddressList[i] = createAddress(&item)
	}
	return replyList, nil
}

func (a *AddressUseCase) DeleteAddress(ctx context.Context, in *v1.DeleteAddressRequest) (*v1.DeleteAddressRely, error) {
	var ad = &AddressInfo{
		Id: int(in.Id),
	}
	ad, err := a.repo.GetAddress(ctx, ad)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no such a address info: %s", err.Error())
		} else {
			a.logger.Errorf("failed to get address info:%s", err.Error())
			return nil, err
		}
	}
	if ad.UserUuid != in.UserUuid {
		return nil, fmt.Errorf("uuid is wrong")
	}
	err = a.repo.DeleteAddress(ctx, ad)
	if err != nil {
		a.logger.Errorf("failed to delete address info: %s", err.Error())
		return nil, err
	}
	return nil, nil
}
func createAddress(src *AddressInfo) *v1.AddressInfo {
	return &v1.AddressInfo{
		Id:       uint64(src.Id),
		UserUuid: src.UserUuid,
		Address:  src.Address,
		Mobile:   src.Mobile,
		Alias:    src.Alias,
	}
}
