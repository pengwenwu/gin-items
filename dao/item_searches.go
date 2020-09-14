package dao

import (
	"gin-items/model"
)

var (
	updateItemSearchesFields = []string{"sku_name", "bar_code", "sku_code", "item_state", "sku_state", "last_dated"}
)

func (dao *Dao) GetItemSearchList(where map[string]interface{}, whereIn model.WhereIn, like map[string]string, order, groupBy string, page, pageSize int) (itemSearchList []*model.ItemSearches, total int64, err error) {
	offset := (page - 1) * pageSize
	query := dao.MasterServiceItems.
		Model(&model.ItemSearches{}).
		Where(where)
	if len(whereIn.ItemId) > 0 {
		query = query.Where("item_id in (?)", whereIn.ItemId)
	}
	if len(whereIn.SkuId) > 0 {
		query = query.Where("sku_id in (?)", whereIn.SkuId)
	}
	if len(like) > 0 {
		for k, v := range like {
			query = query.Where(k+" like ?", "%"+v+"%")
		}
	}
	if order != "" {
		query = query.Order(order)
	}
	if groupBy != "" {
		query = query.Group(groupBy)
	}
	err = query.Count(&total).Error
	err = query.Offset(offset).Limit(pageSize).Find(&itemSearchList).Error
	return
}

func (dao *Dao) InsertSearch(search *model.ItemSearches) error {
	return dao.MasterServiceItems.Create(&search).Error
}

func (dao *Dao) InsertSearches(searchList []*model.ItemSearches) error {
	return dao.MasterServiceItems.Create(&searchList).Error
}

func (dao *Dao) PutUpdateSearch(itemSearch *model.ItemSearches, where map[string]interface{}) error {
	return dao.MasterServiceItems.
		Select(updateItemSearchesFields).
		Where(where).
		Limit(1).
		Updates(&itemSearch).
		Error
}