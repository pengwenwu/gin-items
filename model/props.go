package model

type Props struct {
	Model

	ItemId int `gorm:"index"`
	PropName string `gorm:"type:varchar(255)"`
	Sort int
	HavePhoto int
	PropDesc string `gorm:"type:varchar(255)"`
	State int
}