package service

import (
	"gin-items/lib/app"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"gin-items/dao"
	"gin-items/lib/setting"
)

type ItemService struct {

}

func (itemService *ItemService) GetItemList (c *gin.Context) (map[string]interface{}, error) {
	// 获取参数 && 校验参数
	fields := c.Query("fields")
	itemState := c.DefaultQuery("item_state", "1")
	skuState := c.DefaultQuery("sku_state", "1")
	like := c.QueryMap("like")
	order := c.QueryMap("order")
	page := com.StrTo(c.DefaultQuery("page", "1")).MustInt()
	pageSize := com.StrTo(c.DefaultQuery("limit", com.ToStr(setting.PageSize))).MustInt()
	offset := (page - 1) * pageSize

	where := make(map[string]interface{})
	where["item_id"] = com.StrTo(c.Query("item_id")).MustInt()
	where["sku_id"] = com.StrTo(c.Query("sku_id")).MustInt()
	where["bar_code"] = com.ToStr(c.Query("bar_code"))
	where["sku_code"] = com.ToStr(c.Query("sku_code"))
	for k := range where {
		if where[k] == "" || where[k] == 0 {
			delete(where, k)
		}
	}
	where["item_state"] = com.StrTo(itemState).MustInt()
	where["sku_state"] = com.StrTo(skuState).MustInt()

	for k, v := range like {
		// 对外仅暴露name
		if k == "name" {
			like["sku_name"] = v
			delete(like, k)
		}
	}

	items, err := dao.GetItemList(fields, offset, pageSize, where, like, order)
	if err != nil {
		return nil, err
	}
	total, err := dao.GetItemTotal(where, like)

	data := make(map[string]interface{})
	data["list"] = items
	data["total"] = total

	return data, nil
}

func (itemService *ItemService) GetItem(c *gin.Context) (map[string]interface{}, error) {
	itemId := com.StrTo(c.Query("item_id")).MustInt()
	fields := c.Query("fields")

	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	valid.Required(fields, "fields")
	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
	}
	data := make(map[string]interface{})

	return data, nil
}
