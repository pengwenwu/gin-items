package dao

import (
	"gin-items/library/constant"
	"gin-items/model"
)

func (dao *Dao) GetPhotos(where map[string]interface{}) (photos []*model.ItemPhotos, err error) {
	err = dao.MasterServiceItems.
		Where(where).
		Find(&photos).Error
	return
}

func (dao *Dao) InsertPhoto(photo *model.ItemPhotos) error {
	return dao.MasterServiceItems.Create(&photo).Error
}

func (dao *Dao) InsertPhotos(photos []*model.ItemPhotos) error {
	return dao.MasterServiceItems.Create(photos).Error
}

func (dao *Dao) DeletePhotos(itemId int) error {
	photo := &model.ItemPhotos{State: constant.ItemPhotosStateDeleted}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where("item_id = ? and state = ?", itemId, constant.ItemPhotosStateNormal).
		Limit(constant.CommonLimit).
		Updates(&photo).
		Error
}
