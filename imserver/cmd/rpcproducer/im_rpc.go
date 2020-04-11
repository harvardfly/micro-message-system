package main

/*
im rpc服务 向kafka发送消息 相当于kafka的producer生产者rpc服务 供网关调用
*/

import (
	"flag"
	"github.com/micro/go-plugins/broker/kafka"
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"

	imConfig "micro-message-system/imserver/cmd/config"
	proto "micro-message-system/imserver/protos"
	"micro-message-system/imserver/rpcserveriml"
	"micro-message-system/imserver/util"
)

func main() {
	imFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_rpc.json",
		Usage: "please use xxx -f config_rpc.json",
	}
	configFile := flag.String(imFlag.Name, imFlag.Value, imFlag.Usage)
	flag.Parse()
	conf := new(imConfig.ImRpcConfig)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Address
		})
	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Transport(grpc.NewTransport()),
		micro.Flags(imFlag),
	)
	publisherServerMap := make(map[string]*util.KafkaBroker)
	for _, item := range conf.ImServerList {
		amqbAddress := item.KafkaAddress
		p, err := util.NewKafkaBroker(
			item.Topic,
			kafka.NewBroker(func(options *broker.Options) {
				options.Addrs = amqbAddress
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		publisherServerMap[item.ServerName+item.Topic] = p
	}
	imRpcServer := rpcserveriml.NewImRpcServerIml(publisherServerMap)
	if err := proto.RegisterImHandler(service.Server(), imRpcServer); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
