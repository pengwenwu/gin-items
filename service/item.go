package service

import (
	"gin-items/library/ecode"
	"gin-items/library/token"
	"github.com/astaxie/beego/validation"
	"sync"

	"gin-items/helper"
	"gin-items/library/define"
	"gin-items/library/rabbitmq"
	"gin-items/model"
)

func (serv *Service) GetItemList(params *model.ArgItemSearch, tokenData *token.MyCustomClaims) (itemList []*model.Item, total int64, err error) {
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

	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemInsertData{
			ItemId: itemId,
		})
		pub.Send(rabbitmq.ItemInsert, pubData)
	}()
	go func() {
		defer wg.Done()
		serv.addProps(itemId, item.Props)
	}()
	go func() {
		defer wg.Done()
		serv.addPhotos(itemId, item.Photos)
	}()
	go func() {
		defer wg.Done()
		serv.addParameters(itemId, item.Parameters)
	}()
	wg.Wait()

	return
}

// 添加sku
func (serv *Service) addSkus(itemId int, skus []*model.ItemSkus) {
	for _, sku := range skus {
		sku.ItemId = itemId
	}
	_ = serv.dao.InsertSkus(skus)
	// todo: 报警校验失败
}

// 添加规格属性
func (serv *Service) addProps(itemId int, props []*model.ItemProps) {
	var propValues []*model.ItemPropValues
	for _, prop := range props {
		prop.ItemId = itemId
		prop.State = define.ItemPropsStateNormal

		for _, propValue := range prop.Values {
			propValue.ItemId = itemId
			propValue.PropName = prop.PropName
			propValue.State = define.ItemPropsValuesStateNormal
			propValues = append(propValues, propValue)
		}
	}
	// todo: 报警校验
	_ = serv.dao.InsertProps(props)
	// todo：报警校验
	_ = serv.dao.InsertPropValues(propValues)
}

func (serv *Service) addPhotos(itemId int, photos []*model.ItemPhotos) {
	for _, photo := range photos {
		photo.ItemId = itemId
		photo.State = define.ItemPhotosStateNormal
	}
	// todo: 报警校验
	_ = serv.dao.InsertPhotos(photos)
}

func (serv *Service) addParameters(itemId int, parameters []*model.ItemParameters) {
	for _, parameter := range parameters {
		parameter.ItemId = itemId
		parameter.State = define.ItemParametersStateNormal
	}
	// todo: 报警校验
	_ = serv.dao.InsertParameters(parameters)
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

	err = serv.dao.InsertSearch(itemSearch)
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
	_ = serv.dao.UpdateSkuState(where, define.ItemSkuStateDeletedSelf)
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

	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		pub, _ := rabbitmq.NewProducer()
		pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemUpdateData{
			ItemId: item.ItemId,
		})
		pub.Send(rabbitmq.ItemUpdate, pubData)
	}()
	go func() {
		defer wg.Done()
		_ = serv.dao.DeleteProps(item.ItemId)
		_ = serv.dao.DeletePropValues(item.ItemId)
		serv.addProps(item.ItemId, item.Props)
	}()
	go func() {
		defer wg.Done()
		_ = serv.dao.DeletePhotos(item.ItemId)
		serv.addPhotos(item.ItemId, item.Photos)
	}()
	go func() {
		defer wg.Done()
		_ = serv.dao.DeleteParameters(item.ItemId)
		serv.addParameters(item.ItemId, item.Parameters)
	}()
	wg.Wait()

	return nil
}

func (serv *Service) SyncSkuUpdate(recvData *rabbitmq.SyncSkuUpdateData) {
	itemId := recvData.ItemId
	skuId := recvData.SkuId
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
	_ = serv.dao.PutUpdateSearch(itemSearch, where)
	// todo 错误处理
	return
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
		state = define.ItemStateDeletedReal
	} else {
		state = define.ItemStateDeleted
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
	pub, _ := rabbitmq.NewProducer()
	pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemUpdateData{
		ItemId: itemId,
	})
	pub.Send(rabbitmq.ItemUpdate, pubData)

	return nil
}

func (serv *Service) SyncItemUpdate(recvData *rabbitmq.SyncItemUpdateData) {
	itemId := recvData.ItemId
	where := map[string]interface{}{
		"item_id": itemId,
	}
	itemBase, err := serv.dao.GetItem(where)
	if err != nil {
		// todo: log记录
		return
	}
	where["state"] = itemBase.State
	skuList, err := serv.dao.GetSkus(where)
	if err != nil {
		// todo: log记录
		return
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
		_ = serv.dao.PutUpdateSearch(itemSearch, where)
		// todo 错误处理
	}
	return
}

func (serv *Service) SyncItemInsert(recvData *rabbitmq.SyncItemInsertData) {
	itemId := recvData.ItemId

	where := map[string]interface{}{
		"item_id": itemId,
	}
	itemBase, err := serv.dao.GetItem(where)
	if err != nil {
		// todo: log记录
		return
	}
	where["state"] = itemBase.State
	skuList, err := serv.dao.GetSkus(where)
	if err != nil {
		// todo: log记录
		return
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
	// todo 错误处理
	return
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
		"appkey": tokenData.AppKey,
		"channel": tokenData.Channel,
	}
	_ = serv.dao.UpdateItemState(where, define.ItemStateNormal)
	_ = serv.dao.RecoverSku(where)
	// 发布商品更新消息
	pub, _ := rabbitmq.NewProducer()
	pubData, _ := rabbitmq.MqPack(&rabbitmq.SyncItemUpdateData{
		ItemId: itemId,
	})
	pub.Send(rabbitmq.ItemUpdate, pubData)
	return nil
}
