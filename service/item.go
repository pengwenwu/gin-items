package service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"go.uber.org/zap"

	"gin-items/helper"
	"gin-items/library/constant"
	"gin-items/library/ecode"
	"gin-items/library/rabbitmq"
	"gin-items/library/token"
	"gin-items/middleware/log"
	"gin-items/model"
)

func (serv *Service) GetItemList(param *model.ParamItemSearch, tokenData *token.MyCustomClaims) (itemList []*model.Item, total int64, err error) {
	whereMap := param.GetWhereMap(tokenData)
	like := param.Like
	for k, v := range like {
		if k == "name" {
			like["sku_name"] = v
			delete(like, k)
		}
	}
	itemSearchList, total, err := serv.dao.GetItemSearchList(whereMap, param.WhereIn, like, param.Order, param.GroupBy, param.Page, param.PageSize)
	if err != nil {
		return
	}

	// 查询对应的商品详情
	for _, itemSearch := range itemSearchList {
		item, err := serv.GetItemByItemId(itemSearch.ItemId, tokenData)
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
	item.Skus, err = serv.dao.GetSkuList(where)
	delete(where, "appkey")
	delete(where, "channel")
	where["state"] = constant.ItemPhotosStateNormal
	item.Photos, err = serv.dao.GetPhotos(where)
	where["state"] = constant.ItemParametersStateNormal
	item.Parameters, err = serv.dao.GetParameters(where)
	where["state"] = constant.ItemPropsStateNormal
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
			State:    constant.ItemSkuStateNormal,
		})
	}

	baseItems := &model.Items{
		ItemId:  0,
		Appkey:  item.Appkey,
		Channel: item.Channel,
		Name:    item.Name,
		Photo:   item.Photo,
		Detail:  item.Detail,
		State:   constant.ItemStateNormal,
	}
	itemId, err = serv.dao.InsertItem(baseItems)
	if err != nil {
		log.ErrorLogger.Error("add item error", zap.String("err", err.Error()))
		err = ecode.InsertItemErr
		return
	}
	serv.addSkus(itemId, item.Skus)

	go func() {
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemInsertData{
			ItemId: itemId,
		})
		pub.Send(rabbitmq.ItemInsert, pubData)
	}()
	go func() {
		serv.addProps(itemId, item.Props)
	}()
	go func() {
		serv.addPhotos(itemId, item.Photos)
	}()
	go func() {
		serv.addParameters(itemId, item.Parameters)
	}()

	return
}

// 添加sku
func (serv *Service) addSkus(itemId int, skus []*model.ItemSkus) {
	for _, sku := range skus {
		sku.ItemId = itemId
	}
	if err := serv.dao.InsertSkus(skus); err != nil {
		log.ErrorLogger.Error("add sku error", zap.Int("item_id", itemId), zap.String("error", err.Error()))
	}
}

// 添加规格属性
func (serv *Service) addProps(itemId int, props []*model.ItemProps) {
	var propValues []*model.ItemPropValues
	for _, prop := range props {
		prop.ItemId = itemId
		prop.State = constant.ItemPropsStateNormal

		for _, propValue := range prop.Values {
			propValue.ItemId = itemId
			propValue.PropName = prop.PropName
			propValue.State = constant.ItemPropsValuesStateNormal
			propValues = append(propValues, propValue)
		}
	}
	if err := serv.dao.InsertProps(props); err != nil {
		log.ErrorLogger.Error("insert props error", zap.Int("item_id", itemId))
	}
	if err := serv.dao.InsertPropValues(propValues); err != nil {
		log.ErrorLogger.Error("insert prop_values error", zap.Int("item_id", itemId))
	}
}

func (serv *Service) addPhotos(itemId int, photos []*model.ItemPhotos) {
	for _, photo := range photos {
		photo.ItemId = itemId
		photo.State = constant.ItemPhotosStateNormal
	}
	if err := serv.dao.InsertPhotos(photos); err != nil {
		log.ErrorLogger.Error("add photos error", zap.Int("item_id", itemId))
	}
}

func (serv *Service) addParameters(itemId int, parameters []*model.ItemParameters) {
	for _, parameter := range parameters {
		parameter.ItemId = itemId
		parameter.State = constant.ItemParametersStateNormal
	}
	if  err := serv.dao.InsertParameters(parameters); err != nil{
		log.ErrorLogger.Error("add parameters error", zap.Int("item_id", itemId))
	}
}

func (serv *Service) SyncSkuInsert(recvData *rabbitmq.SyncSkuInsertData) error {
	itemId := recvData.ItemId
	skuId := recvData.SkuId
	itemSearch := &model.ItemSearches{}
	itemBase, err := serv.dao.GetItem(map[string]interface{}{"item_id": itemId})
	if err != nil {
		return err
	}
	skuData, err := serv.dao.GetSku(map[string]interface{}{"sku_id": skuId})
	if err != nil {
		return err
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

	err = serv.dao.InsertSearch(itemSearch)
	if err != nil {
		log.ErrorLogger.Error("insert search error", zap.Int("sku_id", skuId))
	}
	return err
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

	where := map[string]interface{}{
		"appkey":  item.Appkey,
		"channel": item.Channel,
	}
	err := serv.dao.PutUpdateItem(item.Items, where)
	if err != nil {
		err = ecode.UpdateItemErr
		return err
	}

	where["item_id"] = item.ItemId
	_ = serv.dao.UpdateSkuState(where, constant.ItemSkuStateDeletedSelf)
	var newSku []*model.ItemSkus
	for _, sku := range item.Skus {
		if sku.SkuId > 0 {
			_ = serv.dao.PutUpdateSku(sku, where)
		} else {
			newSku = append(newSku, sku)
		}
	}
	if len(newSku) > 0 {
		serv.addSkus(item.ItemId, newSku)
	}

	go func() {
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemUpdateData{
			ItemId: item.ItemId,
		})
		pub.Send(rabbitmq.ItemUpdate, pubData)
	}()
	go func() {
		_ = serv.dao.DeleteProps(item.ItemId)
		_ = serv.dao.DeletePropValues(item.ItemId)
		serv.addProps(item.ItemId, item.Props)
	}()
	go func() {
		_ = serv.dao.DeletePhotos(item.ItemId)
		serv.addPhotos(item.ItemId, item.Photos)
	}()
	go func() {
		_ = serv.dao.DeleteParameters(item.ItemId)
		serv.addParameters(item.ItemId, item.Parameters)
	}()

	return nil
}

func (serv *Service) SyncSkuUpdate(recvData *rabbitmq.SyncSkuUpdateData) error {
	itemId := recvData.ItemId
	skuId := recvData.SkuId
	itemBase, err := serv.dao.GetItem(map[string]interface{}{"item_id": itemId})
	if err != nil {
		return err
	}
	skuData, err := serv.dao.GetSku(map[string]interface{}{"sku_id": skuId})
	if err != nil {
		return err
	}

	where := map[string]interface{}{
		"item_id": itemId,
		"sku_id":  skuId,
		"appkey":  itemBase.Appkey,
		"channel": itemBase.Channel,
	}

	itemSearch := &model.ItemSearches{
		SkuName:   skuData.SkuName,
		BarCode:   skuData.BarCode,
		SkuCode:   skuData.SkuCode,
		ItemState: itemBase.State,
		SkuState:  skuData.State,
	}
	err = serv.dao.PutUpdateSearch(itemSearch, where)
	return err
}

func (serv *Service) DeleteItem(itemId int, isFinalDelete bool, tokenData *token.MyCustomClaims) error {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	if valid.HasErrors() {
		err := helper.GetEcodeValidParam(valid.Errors)
		return err
	}

	where := map[string]interface{}{
		"item_id": itemId,
		"appkey":  tokenData.AppKey,
		"channel": tokenData.Channel,
	}
	var state int
	if isFinalDelete {
		state = constant.ItemStateDeletedReal
	} else {
		state = constant.ItemStateDeleted
	}
	err := serv.dao.UpdateItemState(where, state)
	if err != nil {
		return err
	}
	err = serv.dao.UpdateSkuState(where, state)
	if err != nil {
		return err
	}
	// 发布商品更新消息
	go func() {
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemUpdateData{
			ItemId: itemId,
		})
		pub.Send(rabbitmq.ItemUpdate, pubData)
	}()

	return nil
}

func (serv *Service) SyncItemUpdate(recvData *rabbitmq.SyncItemUpdateData) error {
	itemId := recvData.ItemId
	where := map[string]interface{}{
		"item_id": itemId,
	}
	itemBase, err := serv.dao.GetItem(where)
	return errors.New("s")
	if err != nil {
		return err
	}
	where["state"] = itemBase.State
	skuList, err := serv.dao.GetSkuList(where)
	if err != nil {
		return err
	}

	for _, sku := range skuList {
		where := map[string]interface{}{
			"item_id": itemId,
			"sku_id":  sku.SkuId,
		}
		itemSearch := &model.ItemSearches{
			SkuName:   sku.SkuName,
			BarCode:   sku.BarCode,
			SkuCode:   sku.SkuCode,
			ItemState: itemBase.State,
			SkuState:  sku.State,
		}
		err = serv.dao.PutUpdateSearch(itemSearch, where)
		if err != nil {
			log.ErrorLogger.Error("put update search error", zap.Int("sku_id", sku.SkuId), zap.String("error", err.Error()))
		}
	}
	return err
}

func (serv *Service) SyncItemInsert(recvData *rabbitmq.SyncItemInsertData) error {
	itemId := recvData.ItemId

	where := map[string]interface{}{
		"item_id": itemId,
	}
	itemBase, err := serv.dao.GetItem(where)
	if err != nil {
		return err
	}
	where["state"] = itemBase.State
	skuList, err := serv.dao.GetSkuList(where)
	if err != nil {
		return err
	}

	var searchList []*model.ItemSearches
	for _, sku := range skuList {
		itemSearch := &model.ItemSearches{}
		itemSearch.ItemId = itemBase.ItemId
		itemSearch.Channel = itemBase.Channel
		itemSearch.Appkey = itemBase.Appkey
		itemSearch.ItemState = itemBase.State
		itemSearch.SkuId = sku.SkuId
		itemSearch.SkuName = sku.SkuName
		itemSearch.SkuCode = sku.SkuCode
		itemSearch.BarCode = sku.BarCode
		itemSearch.SkuState = sku.State

		searchList = append(searchList, itemSearch)
	}

	err = serv.dao.InsertSearches(searchList)
	return err
}

func (serv *Service) RecoverItem(itemId int, tokenData *token.MyCustomClaims) error {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	if valid.HasErrors() {
		err := helper.GetEcodeValidParam(valid.Errors)
		return err
	}

	where := map[string]interface{}{
		"item_id": itemId,
		"appkey":  tokenData.AppKey,
		"channel": tokenData.Channel,
	}
	err := serv.dao.UpdateItemState(where, constant.ItemStateNormal)
	if err != nil {
		log.ErrorLogger.Error("update item state error", zap.Int("item_id", itemId), zap.String("error", err.Error()))
		err = ecode.UpdateItemErr
		return err
	}
	err = serv.dao.RecoverSku(where)
	if err != nil {
		log.ErrorLogger.Error("recover item error", zap.Int("item_id", itemId), zap.String("error", err.Error()))
		err = ecode.UpdateItemErr
		return err
	}

	// 发布商品更新消息
	go func() {
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemUpdateData{
			ItemId: itemId,
		})
		pub.Send(rabbitmq.ItemUpdate, pubData)
	}()

	return nil
}
