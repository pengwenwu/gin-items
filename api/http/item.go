package http

import (
	"gin-items/lib/app"
	"gin-items/lib/define"
	"gin-items/lib/ecode"
	"gin-items/service"
	"gin-items/model/item"

	"github.com/gin-gonic/gin"
	"net/http"
)

//获取商品列表
func GetItemList(c *gin.Context) {
	appGin := app.Gin{C: c}

	params := new(item.ArgItemSearch)
	if err := c.Bind(params); err != nil {
		httpResponse()
	}
	params := itemListParams {
		ItemState: define.ITEM_STATE_NORMAL,
		SkuState: define.ITEM_SKU_STATE_NORMAL,
	}
	if c.BindJSON(&params) != nil {
		appGin.Response(http.StatusUnsupportedMediaType, ecode.UnsupportedMediaType, nil)
		return
	}
	//data := make(map[string]interface{})
	//data["params"] = params;
	//appGin.Response(http.StatusOK, ecode.Success, data)
	//return

	itemService := service.ItemService{}
	data, err := itemService.GetItemList(params)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, ecode.ErrorGetItemListFail, nil)
		return
	}

	appGin.Response(http.StatusOK, ecode.Success, data)
	return
}

func GetItem(c *gin.Context) {
	appGin := app.Gin{C: c}

	itemService := service.ItemService{}
	data, err := itemService.GetItem(c)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, ecode.ErrorGetItemFail, nil)
		return
	}

	appGin.Response(http.StatusOK, ecode.Success, data)
	return
}
