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

func GetItemList(offset, limit int, maps interface{}) ([]map[string]interface{}, error)  {
	//var items []*Items
	//err := db.Table("items").Select("name").Offset(offset).Limit(limit).Find(&items).Error
	//if err != nil {
	//	return nil, err
	//}
	//return items, nil
	rows, err := db.Table("items").Select("name,photo,detail").Offset(offset).Limit(limit).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	length := len(columns)
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		current := makeResultReceiver(length)
		if err := rows.Scan(current...); err != nil {
			return nil, err
		}
		value := make(map[string]interface{})
		for i := 0; i < length; i++ {
			k := columns[i]
			v := current[i]
			value[k] = v
		}
		result = append(result, value)
	}
	return result, nil
}

func makeResultReceiver(length int) []interface{}  {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}

func GetItemTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Items{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
