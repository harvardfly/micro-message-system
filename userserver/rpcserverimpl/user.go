package rpcserverimpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"micro-message-system/userserver/models"
	userpb "micro-message-system/userserver/protos"
)

/*
实现rpc的方法：
	FindByToken(context.Context, *FindByTokenRequest, *UserResponse) error
	FindById(context.Context, *FindByIdRequest, *UserResponse) error
*/

type (
	UserRpcServer struct {
		useModel *models.MembersModel
	}
)

var (
	ErrNotFound = errors.New("用户不存在")
)

func NewUserRpcServer(useModel *models.MembersModel) *UserRpcServer {
	return &UserRpcServer{useModel: useModel}
}
func (s *UserRpcServer) FindByToken(ctx context.Context, req *userpb.FindByTokenRequest, rsp *userpb.UserResponse) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FindByToken")
	defer span.Finish()
	member, err := s.useModel.FindByToken(req.Token)
	span.LogFields(
		log.String("event", "checktoken"),
		log.String("member", fmt.Sprintf("%v", member)),
	)
	if err != nil {
		return ErrNotFound
	}
	rsp.Token = member.Token
	rsp.Id = member.Id
	rsp.Username = member.Username
	rsp.Password = member.Password
	return nil
}
func (s *UserRpcServer) FindById(ctx context.Context, req *userpb.FindByIdRequest, rsp *userpb.UserResponse) error {
	member, err := s.useModel.FindById(req.Id)
	if err != nil {
		return ErrNotFound
	}
	rsp.Token = member.Token
	rsp.Id = member.Id
	rsp.Password = member.Password
	return nil
}
