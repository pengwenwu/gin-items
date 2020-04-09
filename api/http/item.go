package http

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"gin-items/library/app"
	"gin-items/library/setting"
	"gin-items/model"
)

//获取商品列表
func GetItemList(c *gin.Context) {
	appGin := app.Gin{C: c}

	argItemSearch := model.ArgItemSearch{
		Page:     setting.Page,
		PageSize: setting.PageSize,
		Order:    "item_id",
		Desc:     "desc",
	}

	if err := c.BindJSON(&argItemSearch); err != nil {
		return
	}

	list, total, err := serv.GetItemList(argItemSearch)
	type pageData struct {
		Data  []*model.Item
		Total int
	}
	data := pageData{list, total}
	appGin.Response(data, err)
	return
}

func GetItemById(c *gin.Context) {
	appGin := app.Gin{C: c}

	itemId := com.StrTo(c.Param("item_id")).MustInt()
	argGetItemById := model.ArgGetItemById{}
	if err := c.BindJSON(&argGetItemById);err != nil{
		//appGin.Response(nil, err)
		return
	}

	item, err := serv.GetItemById(argGetItemById, itemId)
	if err != nil {
		appGin.Response(nil, err)
		return
	}

	appGin.Response(item, nil)
	return
}
