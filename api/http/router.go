package http

import (
	//"gin-items/middleware/jwt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"

	"gin-items/library/setting"
	"gin-items/service"
	"gin-items/model"
	"gin-items/library/ecode"
)


var (
	serv *service.Service
)

func InitRouter() *gin.Engine {
	initService()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	//r.GET("/auth", api.GetAuth)

	itemGroup := r.Group("")
	//apiv1.Use(jwt.Jwt())
	{
		// 获取item列表
		itemGroup.GET("/item", GetItemList)
		// 获取单个item
		itemGroup.GET("/item/:item_id", GetItemById)
	}

	return r
}

func initService()  {
	serv = service.New()
}

func bind(c *gin.Context, v model.ParamValidator) (err error) {
	if err = c.Bind(&v); err != nil {
		err = errors.WithStack(err)
		return
	}
	if !v.Validate() {
		err = ecode.RequestErr
		c.JSON(http.StatusBadRequest, gin.H{
			"code": ecode.InvalidParams,
		})
		return
	}
	return
}