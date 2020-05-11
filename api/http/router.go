package http

import (
	"github.com/gin-gonic/gin"

	"gin-items/library/setting"
	"gin-items/middleware/jwt"
	"gin-items/middleware/log"
	"gin-items/service"
)

var (
	serv *service.Service
)

func InitRouter() *gin.Engine {
	initService()

	r := gin.New()
	r.Use(
		gin.Logger(), // Logger:控制台输出（线上环境可取消）
		gin.Recovery(), // panic异常500处理
		log.LoggerToFile(), // logrus日志
	)
	gin.SetMode(setting.Config().RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	//r.GET("/auth", api.GetAuth)

	itemGroup := r.Group("")
	itemGroup.Use(jwt.Jwt())
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

func initService() {
	serv = service.New()
}
