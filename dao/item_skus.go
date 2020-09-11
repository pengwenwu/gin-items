package dao

import (
	"gin-items/library/ecode"
	"gin-items/model"
)

var (
	updateSkuFields = []string{"item_name", "sku_name", "sku_photo", "sku_code", "bar_code", "properties", "state", "last_dated"}
)

func (dao *Dao) GetSku(where map[string]interface{}) (sku *model.ItemSkus, err error) {
	sku = &model.ItemSkus{}
	err = dao.MasterServiceItems.
		Table(sku.TableName()).
		Where(where).
		Limit(1).
		Find(&sku).Error
	return
}

func (dao *Dao) GetSkus(where map[string]interface{}) (skus []*model.ItemSkus, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemSkus{}.TableName()).
		Where(where).
		Find(&skus).Error
	return
}

func (dao *Dao) InsertSku(sku *model.ItemSkus) error {
	dao.MasterServiceItems.Create(&sku)
	if sku.SkuId == 0 {
		err := ecode.InsertItemErr
		return err
	}
	return nil
}

func (dao *Dao) InsertSkus(skus []*model.ItemSkus) error {
	return dao.MasterServiceItems.Model(&model.ItemSkus{}).Create(skus).Error
}

func (dao *Dao) UpdateSkus(where, update map[string]interface{}) error {
	return dao.MasterServiceItems.
		Model(&model.ItemSkus{}).
		Where(where).
		Limit(commonLimit).
		Updates(update).
		Error
}

func (dao *Dao) UpdateSku(sku *model.ItemSkus, where map[string]interface{}) error {
	return dao.MasterServiceItems.
		Model(&sku).
		Select(updateSkuFields).
		Where(where).
		Limit(1).
		Updates(&sku).
		Error
}

func (dao *Dao) DeleteSkus(where map[string]interface{}, state int) error {
	sku := &model.ItemSkus{State: state}
	return dao.MasterServiceItems.
		Model(&sku).
		Select("state", "last_dated").
		Where(where).
		Limit(commonLimit).
		Updates(&sku).
		Error
}
