package data

import (
	"context"
	"user-srv/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

func NewAddressRepo(data *Data, logger log.Logger) biz.AddressRepo {
	return &addressRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data")),
	}
}

type addressRepo struct {
	data   *Data
	logger *log.Helper
}

func (a *addressRepo) CreateAddress(ctx context.Context, info *biz.AddressInfo) (*biz.AddressInfo, error) {
	err := a.data.maria.Create(info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (a *addressRepo) UpdateAddress(ctx context.Context, info *biz.AddressInfo) (*biz.AddressInfo, error) {
	err := a.data.maria.Model(&biz.AddressInfo{}).Where("id=?", info.Id).Update(info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (a *addressRepo) GetAddress(ctx context.Context, info *biz.AddressInfo) (*biz.AddressInfo, error) {
	err := a.data.maria.Model(&biz.AddressInfo{}).Where("id=? and user_uuid=?", info.Id, info.UserUuid).First(info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (a *addressRepo) DeleteAddress(ctx context.Context, info *biz.AddressInfo) error {
	err := a.data.maria.Model(&biz.AddressInfo{}).Where("id=?", info.Id).Delete(info).Error
	return err
}

func (a *addressRepo) ListAddress(ctx context.Context, query *biz.QueryAddress) ([]biz.AddressInfo, error) {
	var addressList []biz.AddressInfo
	offset := query.Limit * (query.Page - 1)
	err := a.data.maria.Model(&biz.AddressInfo{}).Where("user_uuid=?", query.UserUuid).Offset(offset).Limit(query.Limit).Find(&addressList).Error
	if err != nil {
		return nil, err
	}
	return addressList, nil
}

var _ biz.AddressRepo = (*addressRepo)(nil)
