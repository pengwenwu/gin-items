package http

import (
	"gin-items/helper"
	"gin-items/library/app"
	"gin-items/library/define"
	"gin-items/library/setting"
	"gin-items/library/token"
	"gin-items/model"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

//获取商品列表
func GetItemList(c *gin.Context) {
	appGin := app.Gin{C: c}

	argItemSearch := model.ArgItemSearch{
		ItemState: define.ItemStateNormal,
		SkuState: define.ItemSkuStateNormal,
		Page:     setting.Config().APP.Page,
		PageSize: setting.Config().APP.PageSize,
		Order:    "item_id desc",
		GroupBy: "item_id",
	}

	if bindErr := c.BindJSON(&argItemSearch); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		appGin.Response(nil, err)
		return
	}

	list, total, err := serv.GetItemList(argItemSearch)
	type pageData struct {
		Data  map[int]interface{} `json:"data"`
		Total int `json:"total"`
	}
	data := pageData{list, total}
	appGin.Response(data, err)
	return
}

func GetItemById(c *gin.Context) {
	appGin := app.Gin{C: c}

	itemId := com.StrTo(c.Param("item_id")).MustInt()
	argGetItemById := model.ArgGetItemById{}
	if bindErr := c.BindJSON(&argGetItemById);bindErr != nil{
		err := helper.GetEcodeBindJson(bindErr)
		appGin.Response(nil, err)
		return
	}
	// 先获取item的状态
	argGetItemState := model.ArgGetItemById{Fields:"item_id,state"}
	itemInfo, err := serv.GetItemById(argGetItemState, itemId, define.ItemSkuStateNormal)

	item, err := serv.GetItemById(argGetItemById, itemId, itemInfo["state"])
	if err != nil {
		appGin.Response(nil, err)
		return
	}

	appGin.Response(item, nil)
	return
}

func AddItem(c *gin.Context)  {
	appGin := app.Gin{C: c}

	tokenData, _ := c.Keys["token_data"].(*token.MyCustomClaims)
	item := model.Item{
		Items: model.Items{
			State:define.ItemStateNormal,
			Appkey: tokenData.AppKey,
			Channel: tokenData.Channel,
		},
	}
	if bindErr := c.BindJSON(&item);bindErr != nil{
		err := helper.GetEcodeBindJson(bindErr)
		appGin.Response(nil, err)
		return
	}

	itemId, err := serv.Add(item)
	if err != nil {
		appGin.Response(nil, err)
		return
	}

	appGin.Response(itemId, nil)
}
