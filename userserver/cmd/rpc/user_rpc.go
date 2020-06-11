package main

import (
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"io"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	rl "github.com/juju/ratelimit"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	userRpcConfig "micro-message-system/userserver/cmd/config"
	"micro-message-system/userserver/models"
	userpb "micro-message-system/userserver/protos"
	"micro-message-system/userserver/rpcserverimpl"
)

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 注意：填下地址不能加http:
			LocalAgentHostPort: "192.168.33.16:6831",
		},
	}
	tracer, closer, err := cfg.New(service, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {
	userRpcFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_rpc.json",
		Usage: "please use xxx -f config_rpc.json",
	}
	configFile := flag.String(userRpcFlag.Name, userRpcFlag.Value, userRpcFlag.Usage)
	flag.Parse()
	conf := new(userRpcConfig.RpcConfig)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}
	engineUser, err := gorm.Open(conf.Engine.Name, conf.Engine.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	tracer, _ := initJaeger("micro-message-system.user")
	opentracing.SetGlobalTracer(tracer)
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Address
		});
	b := rl.NewBucketWithRate(float64(conf.Server.RateLimit), int64(conf.Server.RateLimit))
	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Transport(grpc.NewTransport()),
		micro.WrapHandler(
			ratelimit.NewHandlerWrapper(b, false),
			wrapperTrace.NewHandlerWrapper(tracer),
		),
		micro.Flags(userRpcFlag),
	)
	service.Init()
	userModel := models.NewMembersModel(engineUser)
	userRpcServer := rpcserverimpl.NewUserRpcServer(userModel)
	if err := userpb.RegisterUserHandler(service.Server(), userRpcServer); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
