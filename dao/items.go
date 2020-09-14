package dao

import (
	"gin-items/library/constant"
	"gin-items/library/ecode"
	"gin-items/model"
)

var updateItemFields = []string{"name", "photo", "detail", "last_dated"}

func (dao *Dao) GetItem(where map[string]interface{}) (item *model.Items, err error) {
	item = &model.Items{}
	err = dao.MasterServiceItems.
		Where(where).
		Limit(1).
		Find(&item).Error
	return
}

func (dao *Dao) InsertItem(item *model.Items) (itemId int, err error) {
	dao.MasterServiceItems.Create(&item)
	if item.ItemId == 0 {
		err = ecode.InsertItemErr
		return
	}
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
