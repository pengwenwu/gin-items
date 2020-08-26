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

	itemGroup := r.Group("item")
	itemGroup.Use(jwt.Jwt())
	{
		// 获取item列表
		itemGroup.GET("", GetItemList)
		// 获取单个item基础信息
		itemGroup.GET("base/:item_id", GetItemBaseByItemId)
		// 获取item详情
		itemGroup.GET("detail/:item_id", GetItemByItemId)
		// 新增单个item
		itemGroup.POST("", AddItem)
	}

	return r
}

func initService() {
	serv = service.New()
}
