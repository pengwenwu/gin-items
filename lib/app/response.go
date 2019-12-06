package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-items/lib/e"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(errCode int, data map[string]interface{}) {
	// httpCode todo: 自动判断
	httpCode := http.StatusOK

	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:   e.GetMsg(errCode),
		Data:  data,
	})
}