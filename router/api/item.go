package api

import (
	"net/http"

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
		appGin.Response(http.StatusInternalServerError, e.ErrorGetItemListFail, nil)
		return
	}

	appGin.Response(http.StatusOK, e.Success, data)
	return
}