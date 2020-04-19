package http

import (
	"gin-items/library/setting"
	"gin-items/service"
	//"gin-items/middleware/jwt"
	"github.com/gin-gonic/gin"
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
		// 新增单个item
		itemGroup.POST("/item", AddItem)
	}

	return r
}

func initService()  {
	serv = service.New()
}