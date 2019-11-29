package app

import (
	"gin-items/lib/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(data map[string]interface{}) {
	httpCode := http.StatusOK
	g.C.JSON(httpCode, gin.H{
		"code": data["code"],
		"msg": e.GetMsg()
	})
	//g.C.JSON(httpCode, Response{
	//	Code: data["code"],
	//	Msg:   e.GetMsg(errCode),
	//	Data:  data,
	//})
}