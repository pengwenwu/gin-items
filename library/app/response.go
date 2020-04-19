package app

import (

	//"encoding/json"
	"net/http"

	//"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"gin-items/library/ecode"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	State int         `json:"state"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(data interface{}, err error) {
	code := http.StatusOK
	bcode := ecode.Cause(err)
	state := bcode.Code()
	msg := bcode.Message()

	g.C.JSON(code, Response{
		State: state,
		Msg: msg,
		Data: data,
	})
	return
}