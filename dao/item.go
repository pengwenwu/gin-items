package dao

import (
	"gin-items/model"
)

func (d *Dao) SearchItem(arg *model.ArgItemSearch) (itemIds []int, total int, err error)  {
	d.DB.Table(model.Item.TableName()).Rows()
}

func GetItemList(fields string, offset, limit int, where, like, order interface{}) ([]map[string]string, error)  {
	rows, err := db.Table("item_searches").Where(where).Offset(offset).Limit(limit).Rows()
	if err != nil {
		return nil, err
	}
	result := Rows2SliceMap(rows)
	return result, nil
}

func GetItemTotal(where, like interface{}) (int, error) {
	var count int
	err := db.Model(&item.ItemSearches{}).Where(where).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
