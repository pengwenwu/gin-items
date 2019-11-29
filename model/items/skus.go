package items

import "gin-items/model"

type Skus struct {
	model.Model

	Id int `gorm:"-"`

	SkuId int `gorm:"primary_key;AUTO_INCREMENT"`
	ItemId int `gorm:"index"`
	Appkey string `gorm:type:char(32)`
	Channel int
	ItemName string `gorm:"type:varchar(255)"`
	SkuName string `gorm:"type:varchar(255)"`
	SkuPhoto string `gorm:"type:varchar(512)"`
	SkuCode string `gorm:"type:varchar(50)"`
	BarCode string `gorm:"type:varchar(50)"`
	Properties string `gorm:"type:varchar(512)"`
	State int
}
