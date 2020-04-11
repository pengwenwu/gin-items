package service

import (
	"fmt"
	"gin-items/library/define"
	"github.com/astaxie/beego/validation"

	"gin-items/helper"
	"gin-items/library/ecode"
	"gin-items/model"
)

func (serv *Service) GetItemList(params model.ArgItemSearch) (itemList []*model.Item, total int, err error) {
	fields := params.Fields
	offset := (params.Page - 1) * params.PageSize
	whereMap := params.GetWhereMap()
	serv.dao.GetSearchItemIds(params, fields, whereMap, offset, params.PageSize)
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

func (serv *Service) GetItemById(params model.ArgGetItemById, itemId int) (item map[string]interface{}, err error) {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	valid.Required(params.Fields, "fields")
	if valid.HasErrors() {
		err = ecode.ItemIllegalItemId
		return
	}
	item = make(map[string]interface{})
	itemFields := helper.GetItemFields()
	for k, v := range itemFields {
		getField := helper.GetVerifyField(v, params.Fields)
		if len(getField) == 0 {
			continue
		}

		var list []map[string]string
		where := make(map[string]interface{})
		where["item_id"] = itemId
		switch k {
		case "base":
			itemData := make(map[string]string)
			itemData, err = serv.dao.GetItem(itemId, getField)
			for k,v := range itemData {
				item[k] = v
			}
		case "photos":
			where["state"] = define.ITEM_PHOTOS_STATE_NORMAL
			list, err = serv.dao.GetItemPhotos(getField, where, "sort asc", 1, 20)
			item["photos"] = list
		case "parameters":
			where["state"] = define.ITEM_PARAMETERS_STATE_NORMAL
			list, err = serv.dao.GetItemParameters(getField, where, "sort asc", 1, 300)
			item["parameters"] = list
		case "skus":
			where["state"] = define.ITEM_SKU_STATE_NORMAL
			list, err = serv.dao.GetItemSkus(getField, where, "", 1, 20)
			item["skus"] = list
		case "props":
			var propData []model.ItemProps
			propData, err = serv.getPropsData(itemId)
			item["props"] = propData
		}
	}

	return
}

func (serv *Service) getPropsData(itemId int) (propsData []model.ItemProps, err error) {
	where := make(map[string]interface{})
	where["item_id"] = itemId
	where["state"] = define.ITEM_PROPS_STATE_NORMAL
	propsData, err = serv.dao.GetItemProps(where, "sort asc", 1, 20)
	if err != nil {
		return
	}
	where["state"] = define.ITEM_PROPS_VALUES_STATE_NORMAL
	for _,v := range propsData {
		where["prop_name"] = v.PropName
		var tmp []model.ItemPropValues
		tmp, err = serv.dao.GetItemPropValues(where, "sort asc", 1, 20)
		v.Values = tmp
		fmt.Printf("%+v", tmp)
	}
	return
}
