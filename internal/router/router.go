package router

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"micro-shop/internal/handler"
	"micro-shop/internal/middleware"
	"micro-shop/internal/role"
	"micro-shop/internal/service"
	"micro-shop/internal/ws"
	"net/http"
	"os"
)

type Router struct {
	router *gin.Engine

	logger *log.Logger
}

func InitRouter(enf *casbin.Enforcer) (*Router, func(), error) {
	rt := &Router{
		router: gin.Default(),
		logger: log.New(os.Stdout, "router:", log.Lshortfile|log.Ltime),
	}
	rt.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
		return
	})
	logger := log.New(os.Stdout, "router:", log.Lshortfile|log.Lshortfile)
	rl := NewLimiter()
	rt.router.Use(middleware.Jaeger(logger))
	rt.router.Use(middleware.Limiter(rl))
	var funcList = make([]func() error, 0, 5)
	var cancel = func(fList []func() error) func() {
		return func() {
			for _, item := range fList {
				_ = item()
			}
		}
	}

	rt.router.Use(middleware.CORS())
	client, err := consulapi.NewClient(&consulapi.Config{
		Address: viper.GetString("consul.addr"),
	})
	if err != nil {
		return nil, nil, err
	}
	userConf, _, err := client.Agent().Service("service/user", nil)
	if err != nil {
		return nil, nil, err
	}
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userConf.Address, userConf.Port), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	funcList = append(funcList, userConn.Close)
	us := service.NewUserSrv(logger, userConn)
	uh := handler.NewUserHandler(logger, us)
	rt.router.GET("/api/user/v1/mail", uh.Mail)
	rt.router.POST("/api/user/v1/login", uh.Login)
	rt.router.POST("/api/user/v1/register", uh.Register)
	userv1 := rt.router.Group("/api/user/v1")
	userv1.Use(middleware.Jwt(logger, enf))
	userv1.GET("/get/", uh.GetUser)
	userv1.POST("/create", uh.CreateUser)
	userv1.DELETE("/delete/", uh.DeleteUser)
	userv1.PUT("/update", uh.UpdateUser)
	userv1.GET("/list", uh.ListUser)

	arv1 := rt.router.Group("/api/ar/v1")
	arv1.Use(middleware.Jwt(logger, enf))
	arv1.GET("/get", uh.GetAddress)
	arv1.POST("/create", uh.CreateAr)
	arv1.PUT("/update", uh.UpdateAr)
	arv1.DELETE("/delete", uh.DeleteAr)
	arv1.GET("/list", uh.ListAr)

	cmConf, _, err := client.Agent().Service("service/comment", nil)
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}
	cmConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cmConf.Address, cmConf.Port), grpc.WithInsecure())
	if err != nil {
		return nil, cancel(funcList), err
	}
	cmSrv := service.NewCmSrv(logger, cmConn)
	cmHl := handler.NewCmHl(logger, cmSrv)
	rt.router.GET("/api/cm/v1/get", cmHl.GetCm)
	cmv1 := rt.router.Group("/api/cm/v1")
	cmv1.Use(middleware.Jwt(logger, enf))
	cmv1.POST("/create", cmHl.CreateCm)
	cmv1.PUT("/update", cmHl.UpdateCm)
	cmv1.DELETE("/delete", cmHl.DeleteCm)
	cmv1.GET("/list", cmHl.ListCm)

	ntConf, _, err := client.Agent().Service("service/notice", nil)
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}

	ntConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ntConf.Address, ntConf.Port), grpc.WithInsecure())
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}
	funcList = append(funcList, ntConn.Close)
	ns := service.NewNtSrv(logger, ntConn)
	nh := handler.NewNtHl(logger, ns)
	ntv1 := rt.router.Group("/api/nt/v1")
	ntv1.Use(middleware.Jwt(logger, enf))
	ntv1.GET("/get/", nh.GetNt)
	ntv1.POST("/create", nh.CreateNt)
	ntv1.PUT("/update", nh.UpdateNt)
	ntv1.DELETE("/delete/", nh.DeleteNt)
	ntv1.GET("/list", nh.ListNt)

	pdConf, _, err := client.Agent().Service("service/product", nil)
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}
	pdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", pdConf.Address, pdConf.Port), grpc.WithInsecure())
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}
	funcList = append(funcList, pdConn.Close)
	pdSrv := service.NewPdSrv(logger, pdConn)
	ph := handler.NewPdHl(logger, pdSrv)

	rt.router.GET("/api/pd/v1/get", ph.GetPd)
	rt.router.GET("/api/pd/v1/list", ph.ListPd)
	rt.router.GET("/api/pd/v1/listbycg", ph.ListPdByCgId)
	rt.router.GET("/api/pd/v1/listforshop", ph.ListPdForSp)
	pdv1 := rt.router.Group("/api/pd/v1")
	pdv1.Use(middleware.Jwt(logger, enf))
	pdv1.POST("/create", ph.CreatePd)
	pdv1.PUT("/update", ph.UpdatePd)
	pdv1.DELETE("/delete", ph.DeletePd)

	rt.router.GET("/api/sp/v1/get", ph.GetSp)
	rt.router.GET("/api/sp/v1/getbyuid", ph.GetSpByUuid)
	rt.router.GET("/api/sp/v1/list", ph.ListSp)
	spv1 := rt.router.Group("/api/sp/v1")
	spv1.Use(middleware.Jwt(logger, enf))
	spv1.POST("/create", ph.CreateSp)
	spv1.PUT("/update", ph.UpdateSp)
	spv1.DELETE("/delete", ph.DeleteSp)

	rt.router.GET("/api/bn/v1/get", ph.GetBn)
	rt.router.GET("/api/bn/v1/list", ph.ListBn)
	bnv1 := rt.router.Group("/api/bn/v1")
	bnv1.Use(middleware.Jwt(logger, enf))
	bnv1.POST("/create", ph.CreateBn)
	bnv1.PUT("/update", ph.UpdateBn)
	bnv1.DELETE("/delete", ph.DeleteBn)

	rt.router.GET("/api/cg/v1/get", ph.GetCg)
	rt.router.GET("/api/cg/v1/list", ph.ListCg)
	cgv1 := rt.router.Group("/api/cg/v1")
	cgv1.Use(middleware.Jwt(logger, enf))
	cgv1.POST("/create", ph.CreateCg)
	cgv1.PUT("/update", ph.UpdateCg)
	cgv1.DELETE("/delete", ph.DeleteCg)

	ctv1 := rt.router.Group("/api/ct/v1")
	ctv1.Use(middleware.Jwt(logger, enf))
	ctv1.GET("/get", ph.GetCt)
	ctv1.GET("/list", ph.ListCt)
	ctv1.POST("/create", ph.CreateCt)
	ctv1.PUT("/update", ph.UpdateCt)
	ctv1.DELETE("/delete", ph.DeleteCt)

	odConf, _, err := client.Agent().Service("service/order", nil)
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}
	odConn, err := grpc.Dial(fmt.Sprintf("%s:%d", odConf.Address, odConf.Port), grpc.WithInsecure())
	if err != nil {
		cancel(funcList)()
		return nil, nil, err
	}
	funcList = append(funcList, odConn.Close)
	odSrv := service.NewOdSrv(logger, odConn)
	oh := handler.NewOdHl(odSrv, pdSrv, us, ns, logger)
	rt.router.GET("/api/sc/v1/get", oh.GetSc)
	odv1 := rt.router.Group("/api/od/v1")
	odv1.Use(middleware.Jwt(logger, enf))
	odv1.GET("/get/", oh.GetOd)
	odv1.POST("/create", oh.CreateOd)
	odv1.PUT("/update", oh.UpdateOd)
	odv1.DELETE("/delete", oh.DeleteOd)
	odv1.GET("/list", oh.ListOd)
	odv1.GET("/listforsp", oh.ListOdForSp)
	odv1.GET("/listbycate", oh.ListOdByCateId)
	odv1.POST("/pay", oh.PayOd)

	scv1 := rt.router.Group("/api/sc/v1")
	scv1.Use(middleware.Jwt(logger, enf))
	scv1.POST("/create", oh.CreateSc)
	scv1.PUT("/update", oh.UpdateSc)
	scv1.DELETE("/delete", oh.DeleteSc)

	wh := ws.NewWsClient(logger, ns, pdSrv)
	whv1 := rt.router.Group("/api/chat/v1")
	whv1.GET("chat", wh.Upgrade)

	uphl := handler.NewUpHl(
		viper.GetString("oss.ak"),
		viper.GetString("oss.sk"),
		viper.GetString("oss.bucket"),
		&storage.ZoneHuanan,
		logger)
	uhv1 := rt.router.Group("/api/file/v1")
	uhv1.Use(middleware.Jwt(logger, enf))
	uhv1.POST("upload", uphl.UploadFile)

	rolehl := role.NewCb(enf, logger)
	rolev1 := rt.router.Group("/api/role/v1")
	rolev1.GET("/list", rolehl.ListRole)
	return rt, cancel(funcList), nil
}

func (rt *Router) Run(port string) error {
	rt.logger.Println("starting server...")
	return rt.router.Run(port)
}

func (rt *Router) Stop() {
	rt.logger.Println("stopping server")
}
