package model

import (
	"gin-items/library/constant"
	"gin-items/library/token"
)

// ParamValidator .
type ParamValidator interface {
	Validate() bool
}

// ArgItemSearch param struct
type ArgItemSearch struct {
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

func NewArgItemSearch() *ArgItemSearch {
	return &ArgItemSearch{
		ItemState: constant.ItemStateNormal,
		SkuState:  constant.ItemSkuStateNormal,
		Page:      constant.Page,
		PageSize:  constant.PageSize,
		Order:     "item_id desc",
		GroupBy:   "item_id",
	}
}

// Validate .
func (a *ArgItemSearch) Validate() bool {
	return true
}

func (a *ArgItemSearch) GetWhereMap(tokenData *token.MyCustomClaims) (whereMap map[string]interface{}) {
	whereMap = make(map[string]interface{})
	whereMap["appkey"] = tokenData.AppKey
	whereMap["channel"] = tokenData.Channel
	whereMap["item_state"] = a.ItemState
	whereMap["sku_state"] = a.SkuState
	if a.ItemId > 0 {
		whereMap["item_id"] = a.ItemId
	}
	if a.SkuId > 0 {
		whereMap["sku_id"] = a.SkuId
	}
	if a.BarCode != "" {
		whereMap["bar_code"] = a.BarCode
	}
	if a.SkuCode != "" {
		whereMap["sku_code"] = a.SkuCode
	}
	return
}

type ArgSkuList struct {
	ItemId    int               `json:"item_id"`
	SkuId     int               `json:"sku_id"`
	BarCode   string            `json:"bar_code"`
	SkuCode   string            `json:"sku_code"`
	ItemState int               `json:"item_state"`
	SkuState  int               `json:"sku_state"`
	WhereIn   WhereIn           `json:"where_in"`
	Like      map[string]string `json:"like"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
}

func NewArgSkuList() *ArgSkuList {
	return &ArgSkuList{
		ItemState: constant.ItemStateNormal,
		SkuState:  constant.ItemSkuStateNormal,
		Page:      constant.Page,
		PageSize:  constant.PageSize,
	}
}
