package item

import (
	"gin-items/model"
)

type Props struct {
	model.Model

	ItemId int `gorm:"index"`
	PropName string `gorm:"type:varchar(255)"`
	Sort int
	HavePhoto int
	PropDesc string `gorm:"type:varchar(255)"`
	State int
}