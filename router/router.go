package router

import (
	//"gin-items/middleware/jwt"
	"github.com/gin-gonic/gin"

	"gin-items/lib/setting"
	"gin-items/router/api/v1"
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

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.Jwt())

	{
		// 获取item列表
		apiv1.GET("/item", v1.GetItemList)
		apiv1.GET("/item/:item_id", v1.GetItem)
	}

	return r
}