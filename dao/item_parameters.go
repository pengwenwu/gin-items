package dao

import (
	"gin-items/library/define"
	"gin-items/library/ecode"
	"gin-items/model"
)

func (dao *Dao) GetParameters(where map[string]interface{}) (parameters []*model.ItemParameters, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemParameters{}.TableName()).
		Where(where).
		Find(&parameters).Error
	return
}

func (dao *Dao) InsertParameter(parameter *model.ItemParameters) error {
	dao.MasterServiceItems.Create(&parameter)
	if parameter.Id == 0 {
		err := ecode.InsertParameterErr
		return err
	}
	return nil
}

func (dao *Dao) InsertParameters(parameters []*model.ItemParameters) error {
	return dao.MasterServiceItems.Model(&model.ItemParameters{}).Create(parameters).Error
}

func (dao *Dao) DeleteParameters(itemId int) error {
	parameter := &model.ItemParameters{State: define.ItemParametersStateDeleted}
	return dao.MasterServiceItems.
		Model(&parameter).
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, define.ItemParametersStateNormal).
		Limit(commonLimit).
		Updates(&parameter).
		Error
}