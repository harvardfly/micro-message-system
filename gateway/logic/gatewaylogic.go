package logic

import (
	"context"
	"github.com/go-acme/lego/v3/log"
	"time"

	"micro-message-system/common/baseerror"
	"micro-message-system/common/config"
	"micro-message-system/gateway/models"
	"micro-message-system/imserver/protos"
	userProto "micro-message-system/userserver/protos"
)

type (
	GateWayLogic struct {
		//userRpcModel的返回值是userProto.NewUserService()  NewUserService()的返回值是UserService
		userRpcModel userProto.UserService
		gateWayModel *models.GateWayModel
		imRpcModel   im.ImService

		imAddressList []*config.ImRpc
	}
	SendRequest struct {
		FromToken string    `json:"fromToken"  binding:"required"`
		ToToken   string    `json:"toToken"  binding:"required"`
		Body      string    `json:"body"  binding:"required"`
		Timestamp time.Time `json:"timestamp"`
	}

	SendResponse struct {
		Message string
	}

	GetServerAddressRequest struct {
		Token string `json:"token" binding:"required"`
	}

	GetServerAddressResponse struct {
		Address string `json:"address"`
	}
)

var (
	SendMessageErr    = baseerror.NewBaseError("发送消息失败")
	UserNotFoundErr   = baseerror.NewBaseError("用户不存在")
	ImAddressErr      = baseerror.NewBaseError("请配置消息服务地址")
	AddDataErr        = baseerror.NewBaseError("维护关系错误")
	PublishMessageErr = baseerror.NewBaseError("发送消息到MQ失败")
	imRpcModelMapErr  = baseerror.NewBaseError("没有找到对应的RPC服务")
)

func NewGateWayLogic(userRpcModel userProto.UserService,
	gateWayModel *models.GateWayModel,
	imAddressList []*config.ImRpc,
	imRpcModel im.ImService,
) *GateWayLogic {

	return &GateWayLogic{
		userRpcModel:  userRpcModel,
		gateWayModel:  gateWayModel,
		imAddressList: imAddressList,
		imRpcModel:    imRpcModel,
	}
}

func (l *GateWayLogic) Send(r *SendRequest) (*SendResponse, error) {
	// 调用用户服务检查发送给的那个用户token是否存在
	if _, err := l.userRpcModel.FindByToken(context.TODO(), &userProto.FindByTokenRequest{Token: r.ToToken}); err != nil {
		log.Printf("调用用户rpc检查用户：s%", err)
		return nil, UserNotFoundErr
	}
	// 发送给的那个用户检查用户网关信息是否存在 用户需知道绑定哪一台消息中间件服务
	userGate, err := l.gateWayModel.FindByToken(r.ToToken)
	if err != nil {
		log.Printf("检查消息接收方token：s%", err)
		return nil, SendMessageErr
	}
	if userGate.Id < 0 {
		return nil, SendMessageErr
	}
	req := &im.PublishMessageRequest{
		FromToken:  r.FromToken,
		ToToken:    r.ToToken,
		Body:       r.Body,
		ServerName: userGate.ServerName,
		Topic:      userGate.Topic,
		Address:    userGate.ImAddress,
	}
	// 网关存在则调用imRpc服务发送消息
	log.Println(req)
	_, err = l.imRpcModel.PublishMessage(context.TODO(), req);
	// 发送消息逻辑
	if err != nil {
		log.Printf("gateway网关调用IM RPC 发送消息到kafka：s%", err)
		return nil, PublishMessageErr
	}
	// 发送消息逻辑结束
	return &SendResponse{Message: "成功发送消息到kafka"}, nil
}

// 获取用户应该绑定的im websocket服务地址 并绑定
func (l *GateWayLogic) GetServerAddress(r *GetServerAddressRequest) (*GetServerAddressResponse, error) {
	// 验证用户token是否存在
	u, err := l.userRpcModel.FindByToken(context.TODO(), &userProto.FindByTokenRequest{Token: r.Token})
	if err != nil {
		log.Printf("调用UserService FindByToken：s%", err)
		return nil, UserNotFoundErr
	}
	length := len(l.imAddressList)
	if length == 0 {
		return nil, ImAddressErr
	}
	// 用户id 除im服务个数 取余 分配im websocket服务
	i := u.Id % int64(length)
	imData := l.imAddressList[int(i)]

	if _, err := l.gateWayModel.Insert(&models.GateWay{
		Topic:      imData.Topic,
		Token:      r.Token,
		ImAddress:  imData.Address,
		ServerName: imData.ServerName,
	}); err != nil {
		return nil, AddDataErr
	}
	return &GetServerAddressResponse{
		Address: imData.Address,
	}, nil
}
