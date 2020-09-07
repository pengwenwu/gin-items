package service

import (
	"fmt"
	"gin-items/library/rabbitmq"
	"strings"

	"github.com/astaxie/beego/validation"

	"gin-items/helper"
	"gin-items/library/define"
	"gin-items/model"
)

func (serv *Service) GetItemList(params *model.ArgItemSearch) (itemList []*model.Item, total int, err error) {
	//valid := validation.Validation{}
	//valid.Required(params.Fields, "fields")
	//if valid.HasErrors() {
	//	err = helper.GetEcodeValidParam(valid.Errors)
	//	return
	//}
	//fields := params.Fields
	whereMap := params.GetWhereMap()
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
		item, err := serv.GetItemByItemId(itemId)
		if err != nil {
			continue
		}
		itemList = append(itemList, item)
	}
	return
}

func (serv *Service) GetItemBaseByItemId(itemId int) (item *model.Items, err error) {
	valid := validation.Validation{}
	valid.Min(itemId, 1, "item_id")
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}
	item, err = serv.dao.GetItem(itemId, map[string]interface{}{"item_id": itemId})
	return
}

func (serv *Service) GetItemByItemId(itemId int) (item *model.Item, err error) {
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
	}
	item.Items, err = serv.dao.GetItem(itemId, where)

	if err != nil {
		return
	}
	where["state"] = item.Items.State
	item.Skus, err = serv.dao.GetSkus(itemId, where)
	where["state"] = define.ItemPhotosStateNormal
	item.Photos, err = serv.dao.GetPhotos(itemId, where)
	where["state"] = define.ItemParametersStateNormal
	item.Parameters, err = serv.dao.GetParameters(itemId, where)
	where["state"] = define.ItemPropsStateNormal
	item.Props, err = serv.dao.GetProps(itemId, where)
	if len(item.Props) > 0 {
		for k, prop := range item.Props {
			prop.Values, err = serv.dao.GetPropValues(itemId, where)
			item.Props[k] = prop
		}
	}
	return
}

func (serv *Service) Add(item *model.Item) (itemId int, err error) {
	valid := validation.Validation{}
	valid.Valid(&item)
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}

	var propValues []model.ItemPropValues
	if len(item.Props) > 0 {
		for _, prop := range item.Props {
			for _, propValue := range prop.Values {
				propValue.PropName = prop.PropName
				propValues = append(propValues, *propValue)
			}
		}
	}

	if len(item.Skus) > 0 {
		for k, sku := range item.Skus {
			skuName := item.Name
			if sku.Properties != "" {
				skuProps := strings.Split(sku.Properties, ";")
				for _, v := range skuProps {
					skuProp := strings.Split(v, ":")
					skuName += " " + skuProp[1]

					for _, propValue := range propValues {
						if propValue.PropName == skuProp[0] && propValue.PropValueName == skuProp[1] && propValue.PropPhoto != "" {
							sku.SkuPhoto = propValue.PropPhoto
						}
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
		item.Photos = append(item.Photos, &model.ItemPhotos{
			Photo: item.Photo,
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
		Model:   model.Model{},
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
		serv.dao.InsertSku(sku)
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

func (serv *Service) addPhotos(itemId int, photos []*model.ItemPhotos) {
	for _, photo := range photos {
		photo.ItemId = itemId
		photo.State = define.ItemPhotosStateNormal
		serv.dao.InsertPhoto(photo)
		// todo: 报警校验
	}
}

func (serv *Service) addParameters(itemId int, parameters []*model.ItemParameters) {
	for _, parameter := range parameters {
		parameter.ItemId = itemId
		parameter.State = define.ItemParametersStateNormal
		serv.dao.InsertParameter(parameter)
		// todo: 报警校验
	}
}

func (serv *Service) SyncSkuInsert(recvData *rabbitmq.SyncSkuInsertData) {
	fmt.Printf("%+v", recvData)
}
