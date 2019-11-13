package api

import (
	"net/http"

	"github.com/gin-gonic/gin"


)

func GetItemList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code" : 400,
		"msg" : "aaa",
		"data" : "",
	})
}