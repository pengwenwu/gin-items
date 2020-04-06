package dao

import (
	"gin-items/model"
	"github.com/pkg/errors"
)

//func (d *Dao) SearchItem(arg *model.ArgItemSearch) (itemIds []int, total int, err error)  {
//	d.DB.Table(model.Item{}.TableName()).Rows()
//}

func (dao *Dao) GetSearchItemIds(params model.ArgItemSearch, fields string, offset, limit int, where map[string]interface{}) (itemIds []int, err error)  {
	rows, err := dao.DB.Table(params.TableName()).Where(where).Offset(offset).Limit(limit).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var itemId int
		if err = rows.Scan(&itemId); err != nil {
			err = errors.WithStack(err)
			return
		}
		itemIds = append(itemIds, itemId)
	}
	return
}

func (dao *Dao) GetSearchItemTotal(where interface{}) (int, error) {
	var count int
	err := dao.DB.Model(&model.ItemSearches{}).Where(where).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dao *Dao) GetItemById(itemId int, fields string) (item model.Item, err error) {
	err = dao.DB.Table(item.TableName()).Where("item_id = ?", itemId).Find(&item).Error
	if err != nil {
		return
	}
	return
}
