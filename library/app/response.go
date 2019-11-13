package app

import (
	"github.com/gin-gonic/gin"

	"gin-items/library/e"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	State int
	Msg string
	Data interface{}
}

func (g *Gin) Response(httpCode, state int, data interface{}) {
	g.C.JSON(httpCode, Response{
		State: state,
		Msg:   e.GetMsg(state),
		Data:  data,
	})
}