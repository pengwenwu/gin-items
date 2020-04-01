package app

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"gin-items/library/ecode"
)

type Gin struct {
	C *gin.Context
	Error error
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(data interface{}, err error) {
	code := http.StatusOK
	g.Error = err
	bcode := ecode.Cause(err)
	g.C.JSON(code, Response{
		Code: bcode.Code(),
		Msg:  bcode.Message(),
		Data: data,
	})
	return
}