package dao

import (
	"gin-items/library/constant"
	"gin-items/library/ecode"
	"gin-items/model"
)

var (
	updateSkuFields = []string{"item_name", "sku_name", "sku_photo", "sku_code", "bar_code", "properties", "state", "last_dated"}
)

func (dao *Dao) GetSku(where map[string]interface{}) (sku *model.ItemSkus, err error) {
	sku = &model.ItemSkus{}
	err = dao.MasterServiceItems.
		Where(where).
		Take(&sku).Error
	return
}

func (dao *Dao) GetSkuList(where map[string]interface{}) (skuList []*model.ItemSkus, err error) {
	err = dao.MasterServiceItems.
		Where(where).
		Limit(constant.CommonLimit).
		Find(&skuList).Error
	return
}

func (dao *Dao) GetSkuListBySkuIds(skuIds []int) (skuList []*model.ItemSkus, err error) {
	err = dao.MasterServiceItems.
		Where("sku_id in (?)", skuIds).
		Limit(constant.CommonLimit).
		Find(&skuList).
		Error
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
	return dao.MasterServiceItems.Create(skus).Error
}

func (dao *Dao) UpdateSkus(where, update map[string]interface{}) error {
	return dao.MasterServiceItems.
		Model(&model.ItemSkus{}).
		Where(where).
		Limit(constant.CommonLimit).
		Updates(update).
		Error
}

func (dao *Dao) PutUpdateSku(sku *model.ItemSkus, where map[string]interface{}) error {
	return dao.MasterServiceItems.
		Select(updateSkuFields).
		Where(where).
		Limit(1).
		Updates(&sku).
		Error
}

func (dao *Dao) UpdateSkuState(where map[string]interface{}, state int) error {
	sku := &model.ItemSkus{State: state}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where(where).
		Limit(constant.CommonLimit).
		Updates(&sku).
		Error
}

func (dao *Dao) RecoverSku(where map[string]interface{}) error {
	sku := &model.ItemSkus{State: constant.ItemSkuStateNormal}
	return dao.MasterServiceItems.
		Select("state", "last_dated").
		Where(where).
		Where("state in ?", []int{constant.ItemSkuStateDeleted, constant.ItemSkuStateDeletedReal}).
		Limit(constant.CommonLimit).
		Updates(&sku).
		Error
}
