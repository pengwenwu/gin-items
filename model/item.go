package model

type Items struct {
	Model

	ID int `gorm:"-"` // 忽略id字段
	ItemId int `gorm:"primary_key;AUTO_INCREMENT"`
	AppKey string `gorm:"type:char(32)"`
	Channel int
	Name string `gorm:"type:varchar(255)"`
	Photo string `gorm:"type:varchar(512)"`
	Detail string `gorm:"type:text"`
	State int

	Photos Photos
	Parameters Parameters
	Skus Skus
	Props Props
}


