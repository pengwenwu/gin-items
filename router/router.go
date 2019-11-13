package router

import (
	//"gin-items/middleware/jwt"
	"github.com/gin-gonic/gin"

	"gin-items/library/setting"
	"gin-items/router/api"
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

	apiGroup := r.Group("")
	//apiv1.Use(jwt.Jwt())

	{
		// 获取item列表
		apiGroup.GET("/item", api.GetItemList)
	}

	{
		// 获取item列表
		//api.GET("/item", api.GetItems)
		////新建标签
		//apiv1.POST("/tags", v1.AddTag)
		////更新指定标签
		//apiv1.PUT("/tags/:id", v1.EditTag)
		////删除指定标签
		//apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//
		////获取文章列表
		//apiv1.GET("/articles", v1.GetArticles)
		////获取指定文章
		//apiv1.GET("/articles/:id", v1.GetArticle)
		////新建文章
		//apiv1.POST("/articles", v1.AddArticle)
		////更新指定文章
		//apiv1.PUT("/articles/:id", v1.EditArticle)
		////删除指定文章
		//apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}