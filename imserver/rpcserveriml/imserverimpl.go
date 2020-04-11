package rpcserveriml

/*
im rpc 服务方法的实现 向kafka发送数据 producer
*/

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/micro/go-micro/broker"

	"micro-message-system/common/baseerror"
	"micro-message-system/imserver/protos"
	"micro-message-system/imserver/util"
)

type (
	ImRpcServerIml struct {
		sync.Mutex
		publisherServerMap map[string]*util.KafkaBroker
	}
)

var (
	PublishMessageErr = baseerror.NewBaseError("发送消息失败")
)

func NewImRpcServerIml(publisherServerMap map[string]*util.KafkaBroker) *ImRpcServerIml {

	return &ImRpcServerIml{publisherServerMap: publisherServerMap}
}

func (s *ImRpcServerIml) PublishMessage(ctx context.Context, req *im.PublishMessageRequest, rsp *im.PublishMessageResponse) error {
	body, err := json.Marshal(req)
	if err != nil {
		return PublishMessageErr
	}
	key := req.ServerName + req.Topic
	publisher := s.publisherServerMap[key]
	publisher.Publisher(&broker.Message{
		Body: body,
	})
	return nil
}
