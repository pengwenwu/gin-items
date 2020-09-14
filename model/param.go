package model

import (
	"gin-items/library/constant"
	"gin-items/library/token"
)

// ParamItemSearch param struct
type ParamItemSearch struct {
	ItemId    int               `json:"item_id"`
	SkuId     int               `json:"sku_id"`
	BarCode   string            `json:"bar_code"`
	SkuCode   string            `json:"sku_code"`
	ItemState int               `json:"item_state"`
	SkuState  int               `json:"sku_state"`
	Name      string            `json:"name"`
	WhereIn   WhereIn           `json:"where_in"`
	Like      map[string]string `json:"like"`
	Fields    string            `json:"fields"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	Order     string            `json:"order"`
	GroupBy   string            `json:"group_by"`
}

type WhereIn struct {
	ItemId []int `json:"item_id"`
	SkuId  []int `json:"sku_id"`
}

func NewParamItemSearch() *ParamItemSearch {
	return &ParamItemSearch{
		ItemState: constant.ItemStateNormal,
		SkuState:  constant.ItemSkuStateNormal,
		Page:      constant.Page,
		PageSize:  constant.PageSize,
	}
}

func (param *ParamItemSearch) GetWhereMap(tokenData *token.MyCustomClaims) (whereMap map[string]interface{}) {
	whereMap = make(map[string]interface{})
	whereMap["appkey"] = tokenData.AppKey
	whereMap["channel"] = tokenData.Channel
	whereMap["item_state"] = param.ItemState
	whereMap["sku_state"] = param.SkuState
	if param.ItemId > 0 {
		whereMap["item_id"] = param.ItemId
	}
	if param.SkuId > 0 {
		whereMap["sku_id"] = param.SkuId
	}
	if param.BarCode != "" {
		whereMap["bar_code"] = param.BarCode
	}
	if param.SkuCode != "" {
		whereMap["sku_code"] = param.SkuCode
	}
	return
}
