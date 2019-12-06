package v1

import (
	"github.com/gin-gonic/gin"

	"gin-items/lib/app"
	"gin-items/lib/e"
	"gin-items/service"
)

//获取商品列表
func GetItemList(c *gin.Context) {
	appGin := app.Gin{C: c}

	itemService := service.ItemService{}
	data, err := itemService.GetItemList(c)
	if err != nil {
		appGin.Response(e.ErrorGetItemListFail, nil)
		return
	}

	appGin.Response(e.Success, data)
	return
}

func GetItem(c *gin.Context)  {
	appGin := app.Gin{C: c}

	itemService := service.ItemService{}
	data, err := itemService.GetItem(c)
	if err != nil {
		appGin.Response(e.ErrorGetItemFail, nil)
		return
	}

	appGin.Response(e.Success, data)
	return
}