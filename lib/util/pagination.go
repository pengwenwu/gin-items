package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"gin-items/lib/setting"
)

func GetOffset(c *gin.Context) (offset int) {
	offset = 0
	page := com.StrTo(c.PostForm("page")).MustInt()
	if page > 0 {
		offset = (page - 1) * setting.PageSize
 	}
	return
}
