package http

import (
	"gin-items/lib/ecode"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Total int `json:"total,omitempty"`
}

func (g *Gin) httpResponse(httpCode, errCode int, data interface{}, total int) {
	response := Response{
		Code:httpCode,
		Msg:ecode.GetMsg(errCode),
		Data:data,
		Total:total,
	}
	g.C.JSON(httpCode, response)
	return
}

