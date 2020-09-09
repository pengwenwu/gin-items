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

func (dao *Dao) DeleteParameters(itemId int) error {
	return dao.MasterServiceItems.
		Model(&model.ItemParameters{}).
		Where("item_id = ?", itemId).
		Limit(updateCommonLimit).
		Update("state", define.ItemPhotosStateDeleted).
		Error
}