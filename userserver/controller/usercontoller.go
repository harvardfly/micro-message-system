package controller

import (
	"github.com/gin-gonic/gin"

	"micro-message-system/common/baseresponse"
	"micro-message-system/userserver/logic"
)

type (
	UserController struct {
		userLogic *logic.UserLogic
	}
)

func NewUserController(userLogic *logic.UserLogic) *UserController {

	return &UserController{userLogic: userLogic}
}

func (c *UserController) Login(context *gin.Context) {
	r := new(logic.LoginRequest)
	if err := context.ShouldBindJSON(r); err != nil {
		// 通用参数验证方法
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.userLogic.Login(r)
	baseresponse.HttpResponse(context, res, err)
	return
}

func (c *UserController) Register(context *gin.Context) {
	r := new(logic.RegisterRequest)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.userLogic.Register(r)
	baseresponse.HttpResponse(context, res, err)
	return
}
