package http

import (
	"gin-items/helper"
	"gin-items/library/app"
	"gin-items/model"

	"github.com/gin-gonic/gin"
)

func GetSkuList(c *gin.Context)  {
	//tokenData := getTokenData(c)
	argSkuList := model.NewArgSkuList()

	if bindErr := c.BindJSON(&argSkuList); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}
	//
	//list, total, err := serv.GetItemList(argItemSearch, tokenData)
	//resp := &app.ResponseList{}
	//resp.Data = list
	//resp.Total = total
	//app.Response(c, resp, err)
	return
}
