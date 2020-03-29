package item

import (
	"gin-items/model"
)

type Items struct {
	model.Model

	ItemId int `gorm:"primary_key;AUTO_INCREMENT" json:"item_id"`
	Appkey string `gorm:"type:char(32)" json:"appkey"`
	Channel int `json:"channel"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Photo string `gorm:"type:varchar(512)" json:"photo"`
	Detail string `gorm:"type:text" json:"detail"`
	State int `json:"state"`
}


