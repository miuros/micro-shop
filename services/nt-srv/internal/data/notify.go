package data

import (
	"context"
	"nt-srv/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ntRepo struct {
	data   *Data
	logger *log.Helper
}

func (nr *ntRepo) CreateNt(ctx context.Context, nt *biz.Notice) (*biz.Notice, error) {
	err := nr.data.maria.Create(nt).Error
	if err != nil {
		return nil, err
	}
	return nt, nil
}

func (nr *ntRepo) UpdateStatus(ctx context.Context, nt *biz.Notice) error {
	var err error
	switch nt.Type {
	case "chat":
		err = nr.data.maria.Model(&biz.Notice{}).Where("status = 1 and id < ? and toUseUuid=?", nt.Id, nt.UserUuid).Update(&biz.Notice{Status: 2}).Error
	case "notice":
		err = nr.data.maria.Model(&biz.Notice{}).Where("status =1 and id=? and toUserUuid=?", nt.Id, nt.UserUuid).Update(&biz.Notice{Status: 2}).Error

	}
	return err
}

func (nr *ntRepo) DeleteNt(ctx context.Context, nt *biz.Notice) (*biz.Notice, error) {
	err := nr.data.maria.Model(&biz.Notice{}).Where("id=? and to_user_uuid=? or user_uuid=?", nt.Id, nt.UserUuid, nt.UserUuid).Update("is_deleted=1").Error
	if err != nil {
		return nil, err
	}
	return nt, nil
}

func (nr *ntRepo) GetNt(ctx context.Context, nt *biz.Notice) (*biz.Notice, error) {
	err := nr.data.maria.Model(&biz.Notice{}).Where("id=? and to_user_uuid=? or user_uuid=? and type=?", nt.Id, nt.UserUuid, nt.UserUuid, nt.Type).First(nt).Error
	if err != nil {
		return nil, err
	}
	return nt, nil
}

func (nr *ntRepo) ListNt(ctx context.Context, qn *biz.QueryNotice) ([]biz.Notice, int, error) {
	var num int
	var err error
	var ntList []biz.Notice
	switch qn.Type {
	case "chat":
		nr.data.maria.Model(&biz.Notice{}).Where("(user_uuid=? or to_user_uuid=?) and type=?  and status=1", qn.UserUuid, qn.UserUuid, qn.Type).Count(&num)
		err = nr.data.maria.Model(&biz.Notice{}).Where("(user_uuid=? or to_user_uuid=?) and type=?", qn.UserUuid, qn.UserUuid, qn.Type).Order("id desc").Offset(qn.Offset).Limit(qn.Limit).Find(&ntList).Error
	case "notice":
		err = nr.data.maria.Model(&biz.Notice{}).Where("type=?", qn.Type).Count(&num).Error
		if err != nil {
			return nil, 0, err
		}
		if qn.Status != 0 {
			err = nr.data.maria.Model(&biz.Notice{}).Where("to_user_uuid=? and type=? and status=? and is_deleted=0", qn.UserUuid, qn.Type, qn.Status).Offset(qn.Offset).Limit(qn.Limit).Find(&ntList).Error
		} else {
			err = nr.data.maria.Model(&biz.Notice{}).Where("to_user_uuid=? and type=? and is_deleted=0", qn.UserUuid, qn.Type).Offset(qn.Offset).Limit(qn.Limit).Find(&ntList).Error
		}

	}
	if err != nil {
		return nil, 0, err
	}
	return ntList, num, nil
}

var _ biz.NtRepo = (*ntRepo)(nil)

func NewNtRepo(data *Data, logger log.Logger) biz.NtRepo {
	return &ntRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/notice")),
	}
}
