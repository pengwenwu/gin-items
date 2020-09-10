package dao

import (
	"gin-items/library/ecode"
	"gin-items/model"
)

func (dao *Dao) GetItem(where map[string]interface{}) (item *model.Items, err error) {
	item = &model.Items{}
	err = dao.MasterServiceItems.
		Table(item.TableName()).
		Where(where).
		Limit(1).
		Find(&item).Error
	return
}

func (dao *Dao) InsertItem(item *model.Items) (itemId int, err error) {
	dao.MasterServiceItems.Debug().Create(&item)
	if item.ItemId == 0 {
		err = ecode.InsertItemErr
		return
	}
	itemId = item.ItemId
	return
}

func (dao *Dao) UpdateItem(where, update map[string]interface{}) error {
	return dao.MasterServiceItems.
		Model(&model.Items{}).
		Where(where).
		Limit(1).
		Updates(update).
		Error
}
