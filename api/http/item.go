package http

import (
	"gin-items/helper"
	"gin-items/library/app"
	"gin-items/library/define"
	"gin-items/library/rabbitmq"
	"gin-items/library/setting"
	"gin-items/library/token"
	"gin-items/model"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

//获取商品列表
func GetItemList(c *gin.Context) {
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

	list, total, err := serv.GetItemList(argItemSearch)
	resp := &app.ResponseList{}
	resp.Data = list
	resp.Total = total
	app.Response(c, resp, err)
	return
}

// 获取item基础数据
func GetItemBaseByItemId(c *gin.Context) {
	pubMsg()

	itemId := com.StrTo(c.Param("item_id")).MustInt()

	item, err := serv.GetItemBaseByItemId(itemId)
	if err != nil {
		app.Response(c, nil, err)
		return
	}
	resp := &app.ResponseData{Data: item}
	app.Response(c, resp, nil)
	return
}

func GetItemByItemId(c *gin.Context) {
	itemId := com.StrTo(c.Param("item_id")).MustInt()
	//argGetItemById := model.ArgGetItemById{}
	//if bindErr := c.BindJSON(&argGetItemById);bindErr != nil{
	//	err := helper.GetEcodeBindJson(bindErr)
	//	appGin.Response(nil, err)
	//	return
	//}

	item, err := serv.GetItemByItemId(itemId)

	if err != nil {
		app.Response(c, nil, err)
		return
	}

	resp := &app.ResponseData{Data:item}
	app.Response(c, resp, nil)
	return
}

func AddItem(c *gin.Context) {
	tokenData, _ := c.Keys["token_data"].(*token.MyCustomClaims)
	item := &model.Item{
		Items: &model.Items{
			State:   define.ItemStateNormal,
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

	resp := &app.ResponseData{Data:itemId}
	app.Response(c, resp, nil)
}

func pubMsg()  {
	producer, _ := rabbitmq.NewProducer()
	producer.Send(rabbitmq.TradeCreate, "tradeCreate")
	producer.Send(rabbitmq.TradeChange, "tradeChange")
	producer.Close()
	return
}
