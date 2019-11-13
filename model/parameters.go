package model

type Parameters struct {
	Model

	ItemId int `gorm:"index"`
	Parameters string `gorm:"type:varchar(30)"`
	Value string `gorm:"type:varchar(150)"`
	State int
	Sort int
}
