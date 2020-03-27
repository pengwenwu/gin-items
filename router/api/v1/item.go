package v1

import (
	"fmt"
	"gin-items/lib/app"
	"gin-items/lib/e"
	"gin-items/service"
	"github.com/gin-gonic/gin"
)

type itemListParams struct {
	WhereIn map[string][]interface{} `json:"where_in"`
	BarCode string `json:"bar_code"`
	ItemName string `json:"item_name"`
}

//获取商品列表
func GetItemList(c *gin.Context) {
	fmt.Println(1, c.ContentType())
	appGin := app.Gin{C: c}

	var itemListParams itemListParams
	if c.BindJSON(&itemListParams) != nil {
		appGin.Response(e.ErrorContentType, nil)
		return
	}
	data := make(map[string]interface{})
	data["params"] = itemListParams
	appGin.Response(e.Success, data)
	fmt.Println(itemListParams, c)
	return

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