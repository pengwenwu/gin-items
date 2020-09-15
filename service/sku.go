package service

import (
	"gin-items/helper"
	"gin-items/library/token"
	"gin-items/model"
	"github.com/astaxie/beego/validation"
)

func (serv *Service) GetSkuList(param *model.ParamItemSearch, tokenData *token.MyCustomClaims) (skuList []*model.ItemSkus, total int64, err error) {
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
	var skuIds []int
	for _, itemSearch := range itemSearchList {
		skuIds = append(skuIds, itemSearch.SkuId)
	}
	skuList, err = serv.dao.GetSkuListBySkuIds(skuIds)
	return
}

func (serv *Service) GetSkuBySkuId(skuId int, tokenData *token.MyCustomClaims) (sku *model.ItemSkus, err error) {
	valid := validation.Validation{}
	valid.Min(skuId, 1, "sku_id")
	if valid.HasErrors() {
		err = helper.GetEcodeValidParam(valid.Errors)
		return
	}
	where := map[string]interface{}{
		"appkey": tokenData.AppKey,
		"channel": tokenData.Channel,
		"sku_id": skuId,
	}
	sku, err = serv.dao.GetSku(where)
	return
}