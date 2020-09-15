package http

import (
	"gin-items/helper"
	"gin-items/library/app"
	"gin-items/model"
	"github.com/unknwon/com"

	"github.com/gin-gonic/gin"
)

func GetSkuList(c *gin.Context)  {
	tokenData := getTokenData(c)
	paramItemSearch := model.NewParamItemSearch()

	if bindErr := c.BindJSON(&paramItemSearch); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	list, total, err := serv.GetSkuList(paramItemSearch, tokenData)
	resp := &app.ResponseList{
		ResponseData: app.ResponseData{Data:list},
		Total:        total,
	}
	app.Response(c, resp, err)
	return
}

func GetSkuBySkuId(c *gin.Context) {
	tokenData := getTokenData(c)
	skuId := com.StrTo(c.Param("sku_id")).MustInt()

	sku, err := serv.GetSkuBySkuId(skuId, tokenData)
	resp := &app.ResponseData{Data:sku}
	app.Response(c, resp, err)
	return
}
