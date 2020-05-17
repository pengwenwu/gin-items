package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gin-items/library/ecode"
	"gin-items/library/token"
)

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		result := checkToken(context)
		if result.State != ecode.OK.Code() {
			context.Abort()
			context.JSON(http.StatusUnauthorized, result)
		}
		context.Set("token_data", result.Data)
	}
}

func checkToken(context *gin.Context) (result token.DecodeResult) {
	authorization := context.Request.Header.Get("authorization")
	if len(authorization) == 0 {
		result.State = ecode.NoAuthorization.Code()
		result.Msg = ecode.NoAuthorization.Message()
		return
	}
	if strings.Contains(authorization, "Bearer") == false {
		result.State = ecode.AuthorizationErr.Code()
		result.Msg = ecode.AuthorizationErr.Message()
		return
	}

	tokenStr := strings.Replace(authorization, "Bearer ", "", -1)
	tokenDecodeRes := token.UnSafeDecode(tokenStr)
	result = tokenDecodeRes
	return
}