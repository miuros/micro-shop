package main

import (
	"encoding/csv"
	"flag"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"io"
	"micro-shop/internal/router"
	"os"
	"os/signal"
	"syscall"
)

var (
	conf   *string = flag.String("conf", "config/config.yaml", "config file position")
	cas    *string = flag.String("casbin", "config/model.conf", "casbin config file position")
	policy *string = flag.String("policy", "config/policy.csv", "")
)

func init() {
	flag.Parse()

	viper.SetConfigFile(*conf)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.ReadInConfig()
	})
}

func main() {

	enf, err := NewEnforcer(*cas, viper.GetString("maria.driver"), viper.GetString("maria.endpoint"))
	if err != nil {
		panic(err)
	}
	rt, cancel, err := router.InitRouter(enf)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		rt.Run(viper.GetString("http.addr"))
	}()
	select {
	case <-sigChan:
		cancel()
		rt.Stop()
	}
}

func NewEnforcer(conf string, driver, endpoint string) (*casbin.Enforcer, error) {
	db, err := gorm.Open(driver, endpoint)
	if err != nil {
		return nil, err
	}
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	enf, err := casbin.NewEnforcer(conf, adapter)
	if err != nil {
		return nil, err
	}
	enf.EnableLog(true)
	//ok, err := enf.AddPolicies(readCsv(*policy))
	//if err != nil || !ok {
	//	panic(err)
	//}
	err = enf.LoadPolicy()

	if err != nil {
		return nil, err
	}
	return enf, nil
}

func readCsv(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	read := csv.NewReader(f)
	var data [][]string
	for {
		line, err := read.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if len(line) == 0 {
			continue
		}
		data = append(data, line)
	}
	return data
}
