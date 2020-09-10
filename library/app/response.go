package app

import (
	"gin-items/library/ecode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	State int `json:"state"`
	Msg string `json:"msg"`
}

type Responser interface {
	SetBaseInfo(error)
}

type ResponseData struct {
	BaseResponse
	Data interface{} `json:"data"`
}

type ResponseList struct {
	ResponseData
	Total int64 `json:"total"'`
}

func (resp *BaseResponse) SetBaseInfo(err error) {
	bcode := ecode.Cause(err)
	resp.State = bcode.Code()
	resp.Msg = bcode.Message()
}

func (resp *ResponseData) SetBaseInfo(err error) {
	bcode := ecode.Cause(err)
	resp.State = bcode.Code()
	resp.Msg = bcode.Message()
}

func (resp *ResponseList) SetBaseInfo(err error) {
	bcode := ecode.Cause(err)
	resp.State = bcode.Code()
	resp.Msg = bcode.Message()
}

func Response(c *gin.Context, resp Responser, err error)  {
	if resp == nil {
		resp = &BaseResponse{}
	}
	resp.SetBaseInfo(err)
	c.JSON(http.StatusOK, resp)
	return
}