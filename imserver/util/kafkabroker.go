package util

/*
kafka broker 发送/订阅消息
*/

import (
	"log"

	"github.com/micro/go-micro/broker"
)

type (
	KafkaBroker struct {
		topic       string
		kafkaBroker broker.Broker
	}
)

// 初始化kafka连接
func NewKafkaBroker(topic string, kafkaBroker broker.Broker) (*KafkaBroker, error) {
	// 初始化
	if err := kafkaBroker.Init(); err != nil {
		return nil, err
	}
	if err := kafkaBroker.Connect(); err != nil {
		log.Printf("kafka连接失败：s%", err)
		return nil, err
	}
	return &KafkaBroker{topic: topic, kafkaBroker: kafkaBroker}, nil
}

//发送者要写在rpc里，被网关调用
func (p *KafkaBroker) Publisher(msg *broker.Message) {
	if err := p.kafkaBroker.Publish(p.topic, msg); err != nil {
		log.Printf("[publisher %s err] : %+v", p.topic, err)
	}
	log.Printf("[publisher %s] : %s", p.topic, string(msg.Body))
}

//从kafka订阅消息
func (p *KafkaBroker) Subscribe(handlerFunc func(msg []byte) error) {
	log.Println("订阅消息开始")
	sub, err := p.kafkaBroker.Subscribe(p.topic, func(publication broker.Event) error {
		if err := handlerFunc(publication.Message().Body); err != nil {
			log.Println("handlerFunc msg err %+v", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[Subscribe %s err] : %+v", p.topic, err)
	}
	log.Printf("sub:s%", sub)
	log.Printf("[publisher err]")
}
