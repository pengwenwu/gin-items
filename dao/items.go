package dao

import (
	"gin-items/library/constant"
	"gin-items/model"
)

var updateItemFields = []string{"name", "photo", "detail", "last_dated"}

func (dao *Dao) GetItem(where map[string]interface{}) (item *model.Items, err error) {
	item = &model.Items{}
	err = dao.MasterServiceItems.
		Where(where).
		Take(&item).Error
	return
}

func (dao *Dao) InsertItem(item *model.Items) (itemId int, err error) {
	err = dao.MasterServiceItems.Create(&item).Error
	itemId = item.ItemId
	return
}

func (dao *Dao) PutUpdateItem(item *model.Items, where map[string]interface{}) error {
	return dao.MasterServiceItems.
		Select(updateItemFields).
		Where(where).
		Limit(1).
		Updates(&item).
		Error
}

func (dao *Dao) UpdateItemState(where map[string]interface{}, state int) error {
	item := &model.Items{State:state}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where(where).
		Limit(constant.CommonLimit).
		Updates(&item).
		Error
}
