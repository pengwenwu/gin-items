package http

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gin-items/library/rabbitmq"
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

	// 启动mq消费者
	go initMqConsumer()

	r := gin.New()
	r.Use(
		gin.Logger(),       // Logger:控制台输出（线上环境可取消）
		gin.Recovery(),     // panic异常500处理
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
		// 获取多个item
		itemGroup.GET("/getByIds", GetItemByItemIds)
		// 更新item
		itemGroup.PUT("/:item_id", UpdateItem)
		// 删除item
		itemGroup.DELETE("/:item_id", DeleteItem)
	}

	return r
}

func initService() {
	serv = service.New()
}

func initMqConsumer() {
	// 启动多个消费者
	go func() {
		consumer, err := rabbitmq.NewConsumer()
		if err != nil {
			panic(fmt.Errorf("启动mq消费者失败 %s\n", err.Error()))
		}
		consumer.Received(rabbitmq.SyncSkuInsert, func(receivedData []byte) {
			data := &rabbitmq.SyncSkuInsertData{}
			_ = rabbitmq.MqUnpack(receivedData, data)
			serv.SyncSkuInsert(data)
		})
	}()

	go func() {
		consumer, err := rabbitmq.NewConsumer()
		if err != nil {
			panic(fmt.Errorf("启动mq消费者失败 %s\n", err.Error()))
		}
		consumer.Received(rabbitmq.SyncSkuUpdate, func(receivedData []byte) {
			data := &rabbitmq.SyncSkuUpdateData{}
			_ = rabbitmq.MqUnpack(receivedData, data)
			serv.SyncSkuUpdate(data)
		})
	}()

	go func() {
		consumer, err := rabbitmq.NewConsumer()
		if err != nil {
			panic(fmt.Errorf("启动mq消费者失败 %s\n", err.Error()))
		}
		consumer.Received(rabbitmq.SyncItemInsert, func(receivedData []byte) {
			data := &rabbitmq.SyncItemInsertData{}
			_ = rabbitmq.MqUnpack(receivedData, data)
			serv.SyncItemInsert(data)
		})
	}()

	go func() {
		consumer, err := rabbitmq.NewConsumer()
		if err != nil {
			panic(fmt.Errorf("启动mq消费者失败 %s\n", err.Error()))
		}
		consumer.Received(rabbitmq.SyncItemUpdate, func(receivedData []byte) {
			data := &rabbitmq.SyncItemUpdateData{}
			_ = rabbitmq.MqUnpack(receivedData, data)
			serv.SyncItemUpdate(data)
		})
	}()
}
