package dao

import (
	"gin-items/library/constant"
	"gin-items/model"
)

func (dao *Dao) GetProps(where map[string]interface{}) (props []*model.ItemProps, err error) {
	err = dao.MasterServiceItems.
		Where(where).
		Find(&props).Error
	return
}

func (dao *Dao) InsertProp(prop *model.ItemProps) error {
	return dao.MasterServiceItems.Create(&prop).Error
}

func (dao *Dao) InsertProps(props []*model.ItemProps) error {
	return dao.MasterServiceItems.Create(props).Error
}

func (dao *Dao) DeleteProps(itemId int) error {
	props := &model.ItemProps{State: constant.ItemPropsStateDeleted}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, constant.ItemPropsStateNormal).
		Limit(constant.CommonLimit).
		Updates(props).
		Error
}
