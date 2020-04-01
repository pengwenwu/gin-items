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
	appGin.Response(serv.GetItemList(argItemSearch))
	c.JSON(svc.Search(c, v.MIDs))


	params := new(model.ArgItemSearch)
	if err := c.Bind(params); err != nil {
		appGin.Response(http.StatusUnsupportedMediaType, ecode.UnsupportedMediaType, nil)
		return
	}
	data := make(map[string]interface{})
	data["params"] = params
	appGin.Response(http.StatusOK, ecode.Success, data)
	return


	itemService := service.ItemService{}
	data, err := service.GetItemList(params)
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
