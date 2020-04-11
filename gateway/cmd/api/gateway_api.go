package main

/*
gateway的作用：
网关没有rpc 他只有api 网关是调用其它rpc服务的，所以相当于client
1.用户调用网关api发送消息  网关调用kafka producer rpc服务（/send）
2.获取应该绑定哪个im服务的ip 如gateway记录没有则新增 用户与im服务器绑定关系
	token-imAddress  GetServerAddress（/address）
*/

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"

	"micro-message-system/common/middleware"
	gateWayConfig "micro-message-system/gateway/cmd/config"
	"micro-message-system/gateway/controller"
	"micro-message-system/gateway/logic"
	"micro-message-system/gateway/models"
	imProto "micro-message-system/imserver/protos"
	userProto "micro-message-system/userserver/protos"
)

func main() {
	userRpcFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_api.json",
		Usage: "please use xxx -f config_rpc.json",
	}
	configFile := flag.String(userRpcFlag.Name, userRpcFlag.Value, userRpcFlag.Usage)
	flag.Parse()
	conf := new(gateWayConfig.ApiConfig)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}
	engineGateWay, err := gorm.Open(conf.Engine.Name, conf.Engine.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Address
		});

	// Create a new service. Optionally include some options here.
	rpcService := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Transport(grpc.NewTransport()),
		micro.WrapClient(hystrix.NewClientWrapper()), // 客户端熔断
		micro.Flags(userRpcFlag),
	)
	rpcService.Init()
	// 创建用户服务客户端 直接可以通过它调用user prc的服务
	userRpcModel := userProto.NewUserService(conf.UserRpcServer.ServerName, rpcService.Client())
	// 创建IM服务客户端 直接可以通过它调用im prc的服务
	imRpcModel := imProto.NewImService(conf.ImRpcServer.ServerName, rpcService.Client())
	gateWayModel := models.NewGateWayModel(engineGateWay)
	// 把用户服务的客户端注册到网关
	gateLogic := logic.NewGateWayLogic(userRpcModel, gateWayModel, conf.ImRpcServer.ImServerList, imRpcModel)
	gateWayController := controller.NewGateController(gateLogic)
	service := web.NewService(
		web.Name(conf.Server.Name),
		web.Registry(etcdRegisty),
		web.Version(conf.Version),
		web.Flags(userRpcFlag),
		web.Address(conf.Port),
		web.Flags(userRpcFlag),
	)
	router := gin.Default()

	userRouterGroup := router.Group("/gateway")
	// 中间件验证
	userRouterGroup.Use(middleware.ValidAccessToken)
	{
		userRouterGroup.POST("/send", gateWayController.SendHandle)
		userRouterGroup.POST("/address", gateWayController.GetServerAddressHandle)
	}
	service.Handle("/", router)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
