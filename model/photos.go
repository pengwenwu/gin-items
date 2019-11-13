package model

type Photos struct {
	Model

	ItemId int `gorm:"index"`
	Photo string `gorm:"type:varchar(255)"`
	Sort int
	State int
}