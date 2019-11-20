package model

type Items struct {
	Model

	ID int `gorm:"-"` // 忽略id字段
	ItemId int `gorm:"primary_key;AUTO_INCREMENT" json:"item_id"`
	Appkey string `gorm:"type:char(32)" json:"appkey"`
	Channel int `json:"channel"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Photo string `gorm:"type:varchar(512)" json:"photo"`
	Detail string `gorm:"type:text" json:"detail"`
	State int `json:"state"`
	//
	//Photos Photos `json:"photos"`
	//Parameters Parameters `json:"parameters"`
	//Skus Skus `json:"skus"`
	//Props Props `json:"props"`
}

func GetItemList(offset, limit int, maps interface{}) ([]map[string]string, error)  {
	rows, err := db.Table("items").Select("name,photo,detail,state").Offset(offset).Limit(limit).Rows()
	if err != nil {
		return nil, err
	}
	result := Rows2SliceMap(rows)
	return result, nil
}

func GetItemTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Items{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
