package service

import (
	"github.com/astaxie/beego/validation"

	"gin-items/library/ecode"
	"gin-items/model"
)

func (serv *Service) GetItemList(params model.ArgItemSearch) (itemList []*model.Item, total int, err error) {
	fields := params.Fields
	offset := (params.Page - 1) * params.PageSize
	whereMap := params.GetWhereMap()
	serv.dao.GetSearchItemIds(params, fields, offset, params.PageSize, whereMap)
	if err != nil {
		return
	}
	total, err = serv.dao.GetSearchItemTotal(whereMap)
	if total <= 0 {
		return
	}

	// 查询对应的商品详情
	return
}

func (serv *Service) GetItemById(params model.ArgGetItemById, itemId int) (item model.Item, err error) {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	valid.Required(params, "fields")
	if valid.HasErrors() {
		err = ecode.ItemIllegalItemId
		return
	}

	item, err = serv.dao.GetItemById(itemId, params.Fields)
	return
}
