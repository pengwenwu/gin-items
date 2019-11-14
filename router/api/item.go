package api

import (
	"gin-items/lib/app"
	"gin-items/lib/e"
	"gin-items/lib/setting"
	"gin-items/lib/util"
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/unknwon/com"

	"gin-items/service"
)

/**
获取商品列表
 */
func GetItemList(c *gin.Context) {
	appGin := app.Gin{C: c}

	itemsService := service.Items{
		OffSet: util.GetOffset(c),
		PageSize: setting.PageSize,
	}

	total, err := itemsService.Count()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ErrorGetItemCount, nil)
		return
	}

	items, err := itemsService.GetItemList()
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ErrorGetItemListFail, nil)
		return
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["data"] = items

	appGin.Response(http.StatusOK, e.Success, data)
}