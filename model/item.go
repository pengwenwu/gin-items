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

func GetItemList(offset, limit int, maps interface{}) ([]*Items, error)  {
	var items []*Items
	err := db.Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func GetItemTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Items{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
