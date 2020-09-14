package http

import (
	"gin-items/helper"
	"gin-items/library/app"
	"gin-items/library/constant"
	"gin-items/library/token"
	"gin-items/model"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

//获取商品列表
func GetItemList(c *gin.Context) {
	tokenData := getTokenData(c)
	paramItemSearch := model.NewParamItemSearch()
	paramItemSearch.GroupBy = "item_id"

	if bindErr := c.BindJSON(&paramItemSearch); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	list, total, err := serv.GetItemList(paramItemSearch, tokenData)
	resp := &app.ResponseList{}
	resp.Data = list
	resp.Total = total
	app.Response(c, resp, err)
	return
}

// 获取item基础数据
func GetItemBaseByItemId(c *gin.Context) {
	tokenData := getTokenData(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()

	item, err := serv.GetItemBaseByItemId(itemId, tokenData)
	if err != nil {
		app.Response(c, nil, err)
		return
	}
	resp := &app.ResponseData{Data: item}
	app.Response(c, resp, nil)
	return
}

func GetItemByItemId(c *gin.Context) {
	tokenData := getTokenData(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()

	item, err := serv.GetItemByItemId(itemId, tokenData)

	if err != nil {
		app.Response(c, nil, err)
		return
	}

	resp := &app.ResponseData{Data: item}
	app.Response(c, resp, nil)
	return
}

func AddItem(c *gin.Context) {
	tokenData := getTokenData(c)
	item := &model.Item{
		Items: &model.Items{
			State:   constant.ItemStateNormal,
			Appkey:  tokenData.AppKey,
			Channel: tokenData.Channel,
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
	tokenData := getTokenData(c)
	params := new(struct {
		ItemIds []int `json:"item_ids"`
	})
	if bindErr := c.BindJSON(&params); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	itemList, err := serv.GetItemByItemIds(params.ItemIds, tokenData)
	if err != nil {
		app.Response(c, nil, err)
		return
	}

	resp := &app.ResponseData{Data: itemList}
	app.Response(c, resp, nil)
}

func getTokenData(c *gin.Context) *token.MyCustomClaims {
	tokenData, _ := c.Keys["token_data"].(*token.MyCustomClaims)
	return tokenData
}

func UpdateItem(c *gin.Context) {
	tokenData := getTokenData(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()
	item := &model.Item{
		Items: &model.Items{
			ItemId:  itemId,
			Appkey:  tokenData.AppKey,
			Channel: tokenData.Channel,
		},
	}

	if bindErr := c.BindJSON(&item); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	err := serv.UpdateItem(item, tokenData)
	app.Response(c, nil, err)
}

func DeleteItem(c *gin.Context)  {
	tokenData := getTokenData(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()
	params := new(struct {
		IsFinalDelete bool `json:"is_final_delete"`
	})
	if bindErr := c.BindJSON(&params); bindErr != nil {
		err := helper.GetEcodeBindJson(bindErr)
		app.Response(c, nil, err)
		return
	}

	err := serv.DeleteItem(itemId, params.IsFinalDelete, tokenData)
	app.Response(c, nil, err)
	return
}

func RecoverItem(c *gin.Context)  {
	tokenData := getTokenData(c)
	itemId := com.StrTo(c.Param("item_id")).MustInt()

	err := serv.RecoverItem(itemId, tokenData)
	app.Response(c, nil, err)
	return
}
