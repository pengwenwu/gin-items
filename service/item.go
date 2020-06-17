package service

import (
	"gin-items/helper"
	"gin-items/library/define"
	"gin-items/model"
	"github.com/astaxie/beego/validation"
	"strings"
)

func (serv *Service) GetItemList(params model.ArgItemSearch) (itemList map[int]interface{}, total int, err error) {
	valid := validation.Validation{}
	valid.Required(params.Fields, "fields")
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}
	fields := params.Fields
	whereMap := params.GetWhereMap()
	like := params.Like
	for k,v := range like {
		if k == "name" {
			like["sku_name"] = v
			delete(like, k)
		}
	}
	itemIds, err := serv.dao.GetSearchItemIds("item_id", whereMap, params.WhereIn, like, params.Order, params.GroupBy, params.Page, params.PageSize)
	if err != nil {
		return
	}
	total, err = serv.dao.GetSearchItemTotal("count(distinct(item_id))", whereMap, params.WhereIn, like, params.Order)
	if total <= 0 {
		return
	}
	// 查询对应的商品详情
	argGetItem := model.ArgGetItemById{Fields:fields}
	itemList = make(map[int]interface{})
	for _,itemId := range itemIds {
		item, err  := serv.GetItemById(argGetItem, itemId, define.ItemSkuStateNormal)
		if err != nil {
			continue
		}
		itemList[itemId] = item
	}
	return
}

func (serv *Service) GetItemById(params model.ArgGetItemById, itemId int, skuState interface{}) (item map[string]interface{}, err error) {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	valid.Required(params.Fields, "fields")
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
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
			where["state"] = define.ItemPhotosStateNormal
			list, err = serv.dao.GetItemPhotos(getField, where, "sort asc", 1, 20)
			item["photos"] = list
		case "parameters":
			where["state"] = define.ItemParametersStateNormal
			list, err = serv.dao.GetItemParameters(getField, where, "sort asc", 1, 300)
			item["parameters"] = list
		case "skus":
			where["state"] = skuState
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
	where["state"] = define.ItemPropsStateNormal
	propsData, err = serv.dao.GetItemProps(where, "sort asc", 1, 20)
	if err != nil {
		return
	}
	where["state"] = define.ItemPropsValuesStateNormal
	for k,v := range propsData {
		where["prop_name"] = v.PropName
		var tmp []model.ItemPropValues
		tmp, err = serv.dao.GetItemPropValues(where, "sort asc", 1, 20)
		propsData[k].Values = append(v.Values, tmp...)
	}
	return
}

func (serv *Service) Add(item model.Item) (itemId int, err error) {
	valid := validation.Validation{}
	valid.Required(item.Appkey, "appkey")
	valid.Required(item.Name, "name")
	// todo: 参数验证
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}

	var propValues []model.ItemPropValues
	for _, prop := range item.Props {
		for _, propValue := range prop.Values {
			propValue.PropName = prop.PropName
			propValues = append(propValues, propValue)
		}
	}

	for k, sku := range item.Skus {
		skuProps := strings.Split(sku.Properties, ";")
		skuName := item.Name
		for _, v := range skuProps {
			skuProp := strings.Split(v, ":")
			skuName += " " + skuProp[1]

			for _, propValue := range propValues {
				if propValue.PropName == skuProp[0] && propValue.PropValueName == skuProp[1] && propValue.PropPhoto != "" {
					sku.SkuPhoto = propValue.PropPhoto
				}
			}
		}
		if sku.SkuName == "" {
			sku.SkuName = skuName
		}
		sku.Appkey = item.Appkey
		sku.Channel = item.Channel
		sku.ItemName = item.Name
		sku.State = define.ItemSkuStateNormal

		item.Skus[k] = sku
	}
	if len(item.Skus) == 0 {
		item.Skus[0] = model.ItemSkus{
			SkuName:    item.Name,
			SkuPhoto:   item.Photo,
		}
	}
	// 当主图没有时，轮播图第一张图设置为主图
	if len(item.Photos) > 0 && item.Photo == "" {
		item.Photo = item.Photos[0].Photo
	}
	// 当没有轮播图、主图的时候，设置默认图
	if len(item.Photos) == 0 && item.Photo == "" {
		item.Photo = define.ItemDefaultPhoto
	}
	// 当没有轮播图的时候，选设置的第一张默认图
	if len(item.Photos) == 0 {
		item.Photos[0] = model.ItemPhotos{
			Photo:  item.Photo,
		}
	}

	baseItems := model.Items{
		ItemId:  0,
		Appkey:  item.Appkey,
		Channel: item.Channel,
		Name:    item.Name,
		Photo:   item.Photo,
		Detail:  item.Detail,
		Model:   model.Model{},
	}
	itemId, err = serv.dao.InsertItem(baseItems)
	if err != nil {
		return
	}
	serv.addSkus(itemId, item.Skus)
	serv.addProps(itemId, item.Props)

	return
}

// 添加sku
func (serv *Service) addSkus(itemId int, skus []model.ItemSkus)  {
	for _, sku :=range skus {
		sku.ItemId = itemId
		serv.dao.InsertSku(sku)
		// todo: 报警校验失败
	}
}

// 添加规格属性
func (serv *Service) addProps(itemId int, props []model.ItemProps) {
	for _, prop := range props {
		prop.ItemId = itemId
		prop.State = define.ItemPropsStateNormal
		serv.dao.InsertProp(prop)
		// todo: 报警校验
		for _, propValue := range prop.Values {
			propValue.ItemId = itemId
			propValue.PropName = prop.PropName
			propValue.State = define.ItemPropsValuesStateNormal
			serv.dao.InsertPropValue(propValue)
			// todo：报警校验
		}
	}
}
