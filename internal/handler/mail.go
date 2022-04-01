package handler

import (
	"github.com/spf13/viper"
	"math/rand"
	"micro-shop/internal/util"
	"strconv"
	"sync"
	"time"
)

type Val struct {
	Code     string
	DeadLine int64
}

type MlHl struct {
	Ticker *time.Ticker
	Info   map[string]*Val
	locker *sync.Mutex
}

func NewMlHl() *MlHl {
	var mh = &MlHl{
		Info:   make(map[string]*Val),
		Ticker: time.NewTicker(15 * time.Second),
		locker: &sync.Mutex{},
	}
	go mh.countDown()
	return mh
}

func (mh *MlHl) GetCode(mail string) (string, bool) {
	val, ok := mh.Info[mail]
	if !ok {
		return "", false
	}
	return val.Code, true
}

func (mh *MlHl) countDown() {
	for {
		select {
		case <-mh.Ticker.C:
			now := time.Now().Unix()
			for key, val := range mh.Info {
				if val.DeadLine < now {
					mh.locker.Lock()
					delete(mh.Info, key)
					mh.locker.Unlock()
				}
			}
		}
	}
}

func (mh *MlHl) Mail(mail string) error {
	code := generateRandomKey()

	var val = &Val{
		Code:     code,
		DeadLine: time.Now().Add(time.Minute).Unix(),
	}
	mh.locker.Lock()
	mh.Info[mail] = val
	mh.locker.Unlock()
	err := util.SendMail(mail, code, viper.GetString("mail.addr"), viper.GetString("mail.host"), viper.GetString("mail.passwd"))
	return err
}

func generateRandomKey() string {
	var list string = ""
	rand.Seed(time.Now().Unix())
	for i := 0; i < 6; i++ {
		num := rand.Intn(10)
		list += strconv.Itoa(num)
	}
	return list
}
