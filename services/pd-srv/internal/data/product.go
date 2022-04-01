package data

import (
	"context"
	"encoding/json"
	"fmt"
	"pd-srv/internal/biz"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type productRepo struct {
	data   *Data
	logger *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data")),
	}
}

func (pr *productRepo) CreatePd(ctx context.Context, pd *biz.Product) (*biz.Product, error) {
	err := pr.data.maria.Create(pd).Error
	if err != nil {
		return nil, err
	}
	pr.cache(ctx, strconv.Itoa(pd.Id), pd)
	return pd, nil
}

func (pr *productRepo) UpdatePd(ctx context.Context, pd *biz.Product) (*biz.Product, error) {
	err := pr.data.maria.Model(&biz.Product{}).Where("id=?", pd.Id).Update(pd).Error
	if err != nil {
		return nil, err
	}
	pr.cache(ctx, strconv.Itoa(pd.Id), pd)
	return pd, nil
}

func (pr *productRepo) ListForSp(ctx context.Context, qfs *biz.QueryForSp) ([]biz.Product, error) {
	var pdList []biz.Product
	err := pr.data.maria.Model(&biz.Product{}).Where("shop_id=?", qfs.ShopId).Offset(qfs.Offset).Limit(qfs.Limit).Find(&pdList).Error
	if err != nil {
		return nil, err
	}
	return pdList, nil
}

func (pr *productRepo) DeletePd(ctx context.Context, id int) (*biz.Product, error) {
	var pd = &biz.Product{
		Id:        id,
		DeleteAt:  time.Now().Format("2006-01-02:15-04"),
		IsDeleted: 1,
	}
	err := pr.data.maria.Model(&biz.Product{}).Where("id=?", id).Update(pd).Error
	if err != nil {
		return nil, err
	}
	pr.delCache(ctx, strconv.Itoa(id))
	return pd, nil
}

func (pr *productRepo) GetPd(ctx context.Context, id int) (*biz.Product, error) {
	var pd = new(biz.Product)
	res, err := pr.getCache(ctx, strconv.Itoa(id))
	if err == nil {
		err = json.Unmarshal([]byte(res), pd)
		if err == nil {
			return pd, nil
		}
	}
	if err != nil {

		err = pr.data.maria.Model(pd).Where("id=?", id).First(pd).Error
		if err != nil {
			return nil, err
		}

	}
	return pd, nil
}

func (pr *productRepo) ListPd(ctx context.Context, q *biz.Query) ([]biz.Product, error) {
	var pdList []biz.Product
	if len(q.Name) != 0 {
		pr.data.maria = pr.data.maria.Model(&biz.Product{}).Where("name like ?", fmt.Sprintf("%%%s%%", q.Name))
	}
	err := pr.data.maria.Model(&biz.Product{}).Offset(q.Offset).Limit(q.Limit).Find(&pdList).Error
	if err != nil {
		return nil, err
	}
	return pdList, nil
}

func (pr *productRepo) ListPdByCgId(ctx context.Context, q *biz.QueryByCgId) ([]biz.Product, error) {
	var pdList []biz.Product
	err := pr.data.maria.Model(&biz.Product{}).Where("category_id=?", q.CgId).Offset(q.Offset).Limit(q.Limit).Find(&pdList).Error
	if err != nil {
		return nil, err
	}
	return pdList, nil
}

func (pr *productRepo) getCache(ctx context.Context, key string) (string, error) {
	res := pr.data.redis.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.String(), nil
}

func (pr *productRepo) cache(ctx context.Context, key string, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		pr.logger.Errorf("marshal json data error:%s", err.Error())
		return
	}
	pr.data.redis.HSet(ctx, "product", key, string(data))
}

func (pr *productRepo) delCache(ctx context.Context, key string) {
	pr.data.redis.HDel(ctx, "product", key)
}
