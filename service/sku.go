package service

import (
	"gin-items/library/token"
	"gin-items/model"
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
	itemSearchList, total, err := serv.dao.GetItemSearches(whereMap, param.WhereIn, like, param.Order, param.GroupBy, param.Page, param.PageSize)
	if err != nil {
		return
	}
	return
}
