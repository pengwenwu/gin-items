package item

import (
	"gin-items/model"
)

type Photos struct {
	model.Model

	ItemId int `gorm:"index"`
	Photo string `gorm:"type:varchar(255)"`
	Sort int
	State int
}