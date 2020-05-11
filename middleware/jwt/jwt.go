package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gin-items/library/ecode"
)

type result struct {
	State ecode.Code `json:"state"`
	Msg string `json:"msg"`
}

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		result := checkToken(context)
		if result.State != ecode.OK {
			context.Abort()
			context.JSON(http.StatusUnauthorized, result)
		}
	}
}

func checkToken(context *gin.Context) (res result) {
	authorization := context.Request.Header.Get("authorization")
	res = result{
		State: ecode.OK,
		Msg: ecode.OK.Message(),
	}
	if len(authorization) == 0 {
		res.State = ecode.NoAuthorization
		res.Msg = ecode.NoAuthorization.Message()
		return
	}
	if strings.Contains(authorization, "Bearer") == false {
		res.State = ecode.AuthorizationErr
		res.Msg = ecode.AuthorizationErr.Message()
		return
	}

	token := strings.Replace(authorization, "Bearer ", "", -1)
	fmt.Println(token)
	return
}