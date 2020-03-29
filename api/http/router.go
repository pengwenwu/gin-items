package http

import (
	//"gin-items/middleware/jwt"
	"github.com/gin-gonic/gin"

	"gin-items/lib/setting"
)

func InitRouter() *gin.Engine {
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

	itemGroup := r.Group("/item")
	//apiv1.Use(jwt.Jwt())
	{
		// 获取item列表
		itemGroup.GET("/item", GetItemList)
		// 获取单个item
		itemGroup.GET("/item/:item_id", GetItem)
	}

	return r
}