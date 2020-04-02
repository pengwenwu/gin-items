package http

import (
	"gin-items/library/app"
	"gin-items/library/ecode"
	"gin-items/model"
	"gin-items/service"

	"github.com/gin-gonic/gin"
	"net/http"
)

//获取商品列表
func GetItemList(c *gin.Context) {
	appGin := app.Gin{C: c}

	var (
		err error
		argItemSearch = &model.ArgItemSearch{}
	)
	if err = bind(c, argItemSearch); err != nil {
		return
	}
	list, total, err := serv.GetItemList(argItemSearch)
	type pageData struct {
		Data []*model.Item
		Total int
	}
	data := pageData{list, total}
	appGin.Response(data, err)
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
