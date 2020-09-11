package dao

import (
	"gin-items/library/define"
	"gin-items/library/ecode"
	"gin-items/model"
)

func (dao *Dao) GetPropValues(where map[string]interface{}) (propValues []*model.ItemPropValues, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemPropValues{}.TableName()).
		Where(where).
		Find(&propValues).Error
	return
}

func (dao *Dao) InsertPropValue(propValue *model.ItemPropValues) error {
	dao.MasterServiceItems.Create(&propValue)
	if propValue.Id == 0 {
		err := ecode.InsertPropValueErr
		return err
	}
	return nil
}

func (dao *Dao) InsertPropValues(propValues []*model.ItemPropValues) error {
	return dao.MasterServiceItems.Model(&model.ItemPropValues{}).Create(propValues).Error
}

func (dao *Dao) DeletePropValues(itemId int) error {
	propValue := &model.ItemPropValues{State: define.ItemPropsValuesStateDeleted}
	return dao.MasterServiceItems.
		Model(&propValue).
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, define.ItemPropsValuesStateNormal).
		Limit(commonLimit).
		Updates(&propValue).
		Error
}
