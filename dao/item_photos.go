package dao

import (
	"gin-items/library/define"
	"gin-items/library/ecode"
	"gin-items/model"
)

func (dao *Dao) GetPhotos(where map[string]interface{}) (photos []*model.ItemPhotos, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemPhotos{}.TableName()).
		Where(where).
		Find(&photos).Error
	return
}

func (dao *Dao) InsertPhoto(photo *model.ItemPhotos) error {
	dao.MasterServiceItems.Create(&photo)
	if photo.Id == 0 {
		err := ecode.InsertPhotoErr
		return err
	}
	return nil
}

func (dao *Dao) InsertPhotos(photos []*model.ItemPhotos) error {
	return dao.MasterServiceItems.Model(&model.ItemPhotos{}).Create(photos).Error
}

func (dao *Dao) DeletePhotos(itemId int) error {
	return dao.MasterServiceItems.
		Model(&model.ItemPhotos{}).
		Where("item_id = ?", itemId).
		Limit(updateCommonLimit).
		Update("state", define.ItemPhotosStateDeleted).
		Error
}