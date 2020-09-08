package service

import (
	"fmt"
	"gin-items/library/token"
	"github.com/astaxie/beego/validation"

	"gin-items/helper"
	"gin-items/library/define"
	"gin-items/library/rabbitmq"
	"gin-items/model"
)

func (serv *Service) GetItemList(params *model.ArgItemSearch, tokenData *token.MyCustomClaims) (itemList []*model.Item, total int, err error) {
	//valid := validation.Validation{}
	//valid.Required(params.Fields, "fields")
	//if valid.HasErrors() {
	//	err = helper.GetEcodeValidParam(valid.Errors)
	//	return
	//}
	//fields := params.Fields
	whereMap := params.GetWhereMap()
	whereMap["appkey"] = tokenData.AppKey
	whereMap["channel"] = tokenData.Channel
	like := params.Like
	for k, v := range like {
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
	for _, itemId := range itemIds {
		item, err := serv.GetItemByItemId(itemId, tokenData)
		if err != nil {
			continue
		}
		itemList = append(itemList, item)
	}
	return
}

func (serv *Service) GetItemBaseByItemId(itemId int, tokenData *token.MyCustomClaims) (item *model.Items, err error) {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}
	item, err = serv.dao.GetItem(map[string]interface{}{
		"item_id": itemId,
		"appkey":  tokenData.AppKey,
		"channel": tokenData.Channel,
	})
	return
}

func (serv *Service) GetItemByItemId(itemId int, tokenData *token.MyCustomClaims) (item *model.Item, err error) {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}
	item = &model.Item{
		Items:      &model.Items{},
		Photos:     nil,
		Parameters: nil,
		Skus:       nil,
		Props:      nil,
	}
	where := map[string]interface{}{
		"item_id": itemId,
		"appkey":  tokenData.AppKey,
		"channel": tokenData.Channel,
	}
	item.Items, err = serv.dao.GetItem(where)

	if err != nil {
		return
	}
	where["state"] = item.Items.State
	item.Skus, err = serv.dao.GetSkus(where)
	delete(where, "appkey")
	delete(where, "channel")
	where["state"] = define.ItemPhotosStateNormal
	item.Photos, err = serv.dao.GetPhotos(where)
	where["state"] = define.ItemParametersStateNormal
	item.Parameters, err = serv.dao.GetParameters(where)
	where["state"] = define.ItemPropsStateNormal
	item.Props, err = serv.dao.GetProps(where)
	if len(item.Props) > 0 {
		for k, prop := range item.Props {
			prop.Values, err = serv.dao.GetPropValues(where)
			item.Props[k] = prop
		}
	}
	return
}

func (serv *Service) Add(item *model.Item) (itemId int, err error) {
	valid := &validation.Validation{}
	item.Valid(valid)
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}

	if len(item.Skus) == 0 {
		item.Skus = append(item.Skus, &model.ItemSkus{
			Appkey:   item.Appkey,
			Channel:  item.Channel,
			ItemName: item.Name,
			SkuName:  item.Name,
			SkuPhoto: item.Photo,
			State:    define.ItemSkuStateNormal,
		})
	}

	baseItems := &model.Items{
		ItemId:  0,
		Appkey:  item.Appkey,
		Channel: item.Channel,
		Name:    item.Name,
		Photo:   item.Photo,
		Detail:  item.Detail,
		State:   define.ItemStateNormal,
	}
	itemId, err = serv.dao.InsertItem(baseItems)
	if err != nil {
		return
	}
	serv.addSkus(itemId, item.Skus)
	serv.addProps(itemId, item.Props)
	serv.addPhotos(itemId, item.Photos)
	serv.addParameters(itemId, item.Parameters)

	return
}

// 添加sku
func (serv *Service) addSkus(itemId int, skus []*model.ItemSkus) {
	for _, sku := range skus {
		sku.ItemId = itemId
		_ = serv.dao.InsertSku(sku)
		// todo: 报警校验失败
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncSkuInsertData{
			ItemId: itemId,
			SkuId:  sku.SkuId,
		})
		pub.Send(rabbitmq.SkuInsert, pubData)
	}
}

// 添加规格属性
func (serv *Service) addProps(itemId int, props []*model.ItemProps) {
	for _, prop := range props {
		prop.ItemId = itemId
		prop.State = define.ItemPropsStateNormal
		_ = serv.dao.InsertProp(prop)
		// todo: 报警校验
		for _, propValue := range prop.Values {
			propValue.ItemId = itemId
			propValue.PropName = prop.PropName
			propValue.State = define.ItemPropsValuesStateNormal
			_ = serv.dao.InsertPropValue(propValue)
			// todo：报警校验
		}
	}
}

func (serv *Service) addPhotos(itemId int, photos []*model.ItemPhotos) {
	for _, photo := range photos {
		photo.ItemId = itemId
		photo.State = define.ItemPhotosStateNormal
		_ = serv.dao.InsertPhoto(photo)
		// todo: 报警校验
	}
}

func (serv *Service) addParameters(itemId int, parameters []*model.ItemParameters) {
	for _, parameter := range parameters {
		parameter.ItemId = itemId
		parameter.State = define.ItemParametersStateNormal
		_ = serv.dao.InsertParameter(parameter)
		// todo: 报警校验
	}
}

func (serv *Service) SyncSkuInsert(recvData *rabbitmq.SyncSkuInsertData) {
	itemId := recvData.ItemId
	skuId := recvData.SkuId
	itemSearch := &model.ItemSearches{}
	itemBase, err := serv.dao.GetItem(map[string]interface{}{"item_id": itemId})
	if err != nil {
		// todo: log记录
		return
	}
	skuData, err := serv.dao.GetSku(map[string]interface{}{"sku_id": skuId})
	if err != nil {
		// todo: log记录
		return
	}

	itemSearch.ItemId = itemBase.ItemId
	itemSearch.Channel = itemBase.Channel
	itemSearch.Appkey = itemBase.Appkey
	itemSearch.ItemState = itemBase.State
	itemSearch.SkuId = skuData.SkuId
	itemSearch.SkuName = skuData.SkuName
	itemSearch.SkuCode = skuData.SkuCode
	itemSearch.BarCode = skuData.BarCode
	itemSearch.SkuState = skuData.State

	_ = serv.dao.InsertSearches(itemSearch)
	// todo 错误处理
	return
}

func (serv *Service) GetItemByItemIds(itemIds []int, tokenData *token.MyCustomClaims) ([]*model.Item, error) {
	valid := validation.Validation{}
	valid.MinSize(itemIds, 1, "item_ids")
	if valid.HasErrors() {
		err := helper.GetEcodeValidParam(valid.Errors)
		return nil, err
	}

	var itemList []*model.Item
	for _, itemId := range itemIds {
		itemDetail, err := serv.GetItemByItemId(itemId, tokenData)
		if err != nil {
			continue
		}
		itemList = append(itemList, itemDetail)
	}
	return itemList, nil
}

func (serv *Service) UpdateItem(item *model.Item, tokenData *token.MyCustomClaims) error {
	valid := &validation.Validation{}
	valid.Min(item.ItemId, 1, "item_id")
	item.Valid(valid)
	if valid.HasErrors() {
		err := helper.GetEcodeValidParam(valid.Errors)
		return err
	}

	fmt.Printf("%+v", item.Items)
	err := serv.dao.UpdateItem(item.Items, map[string]interface{}{"item_id": item.ItemId, "appkey": item.Appkey, "channel": item.Channel})
	fmt.Println(err)

	return nil
}
