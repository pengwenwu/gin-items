package item

import (
	"gin-items/model"
)

type ItemSearches struct {
	model.Model

	Id int `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Appkey string `gorm:"type:char(32)" json:"appkey"`
	Channel int `json:"channel"`
	ItemId int `json:"item_id"`
	SkuId int `json:"sku_id"`
	SkuName string `gorm:"type:varchar(255)" json:"sku_name"`
	BarCode string `gorm:"type:varchar(20)" json:"bar_code"`
	SkuCode string
}