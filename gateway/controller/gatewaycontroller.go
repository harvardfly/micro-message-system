package controller

import (
	"github.com/gin-gonic/gin"
	"micro-message-system/common/baseresponse"
	"micro-message-system/gateway/logic"
)

type (
	GateController struct {
		gateLogic *logic.GateWayLogic
	}
)

func NewGateController(gateLogic *logic.GateWayLogic) *GateController {

	return &GateController{gateLogic: gateLogic}
}

// 发送消息  用户之间  token - token
func (gat *GateController) SendHandle(c *gin.Context) {
	req := new(logic.SendRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		baseresponse.ParamError(c, err)
		return
	}
	res, err := gat.gateLogic.Send(req)
	baseresponse.HttpResponse(c, res, err)
	return
}

func (gat *GateController) GetServerAddressHandle(c *gin.Context) {
	req := new(logic.GetServerAddressRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		baseresponse.ParamError(c, err)
		return
	}
	res, err := gat.gateLogic.GetServerAddress(req)
	baseresponse.HttpResponse(c, res, err)
	return
}
