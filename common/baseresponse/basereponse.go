package baseresponse

/*
定义通用参数验证和返回值
*/

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"micro-message-system/common/baseerror"
	"micro-message-system/common/exception"
)

// 通用参数验证方法
func ParamError(ctx *gin.Context, err interface{}) {
	validErr, ok := err.(validator.ValidationErrors)
	if ok {
		errMap := map[string]string{}
		for _, ve := range validErr {
			key := ve.FieldNamespace + "." + ve.Tag
			errMap[key] = exception.ZhMessage[key]
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errMap})
		return
	}
	// gin.H 实际上就是 map[string]interface{}
	ctx.JSON(http.StatusBadRequest, gin.H{"message": exception.ErrParam})
	return
}

func HttpResponse(ctx *gin.Context, res, err interface{}) {
	baeError, ok := err.(*baseerror.BaseError)
	if ok {
		ctx.JSON(http.StatusOK, gin.H{"message": baeError.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": exception.ErrServer})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res})
	return
}
