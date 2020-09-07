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
	token := getToken(c)
	argItemSearch := &model.ArgItemSearch{
		ItemState: define.ItemStateNormal,
		SkuState:  define.ItemSkuStateNormal,
		Page:      setting.Config().APP.Page,
		PageSize:  setting.Config().APP.PageSize,
		Order:     "item_id desc",
		GroupBy:   "item_id",
	}

	if bindErr := c.BindJSON(&argItemSearch); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	list, total, err := serv.GetItemList(argItemSearch, token)
	resp := &app.ResponseList{}
	resp.Data = list
	resp.Total = total
	app.Response(c, resp, err)
	return
}

// 获取item基础数据
func GetItemBaseByItemId(c *gin.Context) {
	token := getToken(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()

	item, err := serv.GetItemBaseByItemId(itemId, token)
	if err != nil {
		app.Response(c, nil, err)
		return
	}
	resp := &app.ResponseData{Data: item}
	app.Response(c, resp, nil)
	return
}

func GetItemByItemId(c *gin.Context) {
	token := getToken(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()
	//argGetItemById := model.ArgGetItemById{}
	//if bindErr := c.BindJSON(&argGetItemById);bindErr != nil{
	//	err := helper.GetEcodeBindJson(bindErr)
	//	appGin.Response(nil, err)
	//	return
	//}

	item, err := serv.GetItemByItemId(itemId, token)

	if err != nil {
		app.Response(c, nil, err)
		return
	}

	resp := &app.ResponseData{Data: item}
	app.Response(c, resp, nil)
	return
}

func AddItem(c *gin.Context) {
	token := getToken(c)
	item := &model.Item{
		Items: &model.Items{
			State:   define.ItemStateNormal,
			Appkey:  token.AppKey,
			Channel: token.Channel,
		},
	}

	if bindErr := c.BindJSON(&item); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	itemId, err := serv.Add(item)
	if err != nil {
		app.Response(c, nil, err)
		return
	}

	resp := &app.ResponseData{Data: itemId}
	app.Response(c, resp, nil)
}

func GetItemByItemIds(c *gin.Context) {
	token := getToken(c)
	params := new(struct {
		ItemIds []int `json:"item_ids"`
	})
	if bindErr := c.BindJSON(&params); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	itemList, err := serv.GetItemByItemIds(params.ItemIds, token)
	if err != nil {
		app.Response(c, nil, err)
		return
	}

	resp := &app.ResponseData{Data: itemList}
	app.Response(c, resp, nil)
}

func getToken(c *gin.Context) *token.MyCustomClaims {
	tokenData, _ := c.Keys["token_data"].(*token.MyCustomClaims)
	return tokenData
}
