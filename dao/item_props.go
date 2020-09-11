package dao

import (
	"gin-items/library/define"
	"gin-items/library/ecode"
	"gin-items/model"
)

func (dao *Dao) GetProps(where map[string]interface{}) (props []*model.ItemProps, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemProps{}.TableName()).
		Where(where).
		Find(&props).Error
	return
}

func (dao *Dao) InsertProp(prop *model.ItemProps) error {
	dao.MasterServiceItems.Create(&prop)
	if prop.Id == 0 {
		err := ecode.InsertPropErr
		return err
	}
	return nil
}

func (dao *Dao) InsertProps(props []*model.ItemProps) error {
	return dao.MasterServiceItems.Model(&model.ItemProps{}).Create(props).Error
}

func (dao *Dao) DeleteProps(itemId int) error {
	props := &model.ItemProps{State: define.ItemPropsStateDeleted}
	return dao.MasterServiceItems.
		Model(&props).
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, define.ItemPropsStateNormal).
		Limit(commonLimit).
		Updates(props).
		Error
}
