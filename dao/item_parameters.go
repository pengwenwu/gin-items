package dao

import (
	"gin-items/library/constant"
	"gin-items/model"
)

func (dao *Dao) GetParameters(where map[string]interface{}) (parameters []*model.ItemParameters, err error) {
	err = dao.MasterServiceItems.
		Where(where).
		Find(&parameters).Error
	return
}

func (dao *Dao) InsertParameter(parameter *model.ItemParameters) error {
	return dao.MasterServiceItems.Create(&parameter).Error
}

func (dao *Dao) InsertParameters(parameters []*model.ItemParameters) error {
	return dao.MasterServiceItems.Create(&parameters).Error
}

func (dao *Dao) DeleteParameters(itemId int) error {
	parameter := &model.ItemParameters{State: constant.ItemParametersStateDeleted}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, constant.ItemParametersStateNormal).
		Limit(constant.CommonLimit).
		Updates(&parameter).
		Error
}