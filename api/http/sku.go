package http

import (
	"gin-items/helper"
	"gin-items/library/app"
	"gin-items/model"

	"github.com/gin-gonic/gin"
)

func GetSkuList(c *gin.Context)  {
	//tokenData := getTokenData(c)
	paramItemSearch := model.NewParamItemSearch()

	if bindErr := c.BindJSON(&paramItemSearch); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	//list, total, err := serv.GetSkuList(argItemSearch, tokenData)
	//resp := &app.ResponseList{}
	//resp.Data = list
	//resp.Total = total
	//app.Response(c, resp, err)
	return
}
