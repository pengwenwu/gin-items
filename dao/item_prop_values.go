package dao

import (
	"gin-items/library/constant"
	"gin-items/model"
)

func (dao *Dao) GetPropValues(where map[string]interface{}) (propValues []*model.ItemPropValues, err error) {
	err = dao.MasterServiceItems.
		Where(where).
		Find(&propValues).Error
	return
}

func (dao *Dao) InsertPropValue(propValue *model.ItemPropValues) error {
	return dao.MasterServiceItems.Create(&propValue).Error
}

func (dao *Dao) InsertPropValues(propValues []*model.ItemPropValues) error {
	return dao.MasterServiceItems.Create(propValues).Error
}

func (dao *Dao) DeletePropValues(itemId int) error {
	propValue := &model.ItemPropValues{State: constant.ItemPropsValuesStateDeleted}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, constant.ItemPropsValuesStateNormal).
		Limit(constant.CommonLimit).
		Updates(&propValue).
		Error
}
