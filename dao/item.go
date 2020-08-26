package dao

import (
	"github.com/pkg/errors"

	"gin-items/library/ecode"
	"gin-items/model"
)

func (dao *Dao) GetSearchItemIds(fields string, where map[string]interface{}, whereIn model.WhereIn, like map[string]string, order, groupBy string, page, pageSize int) (itemIds []int, err error) {
	offset := (page - 1) * pageSize
	query := dao.MasterServiceItems.
		Table(model.ItemSearches{}.TableName()).
		Select(fields).
		Where(where)
	if len(whereIn.ItemId) > 0 {
		query = query.Where("item_id in (?)", whereIn.ItemId)
	}
	if len(like) > 0 {
		for k, v := range like {
			query = query.Where(k+" like ?", "%"+v+"%")
		}
	}
	rows, err := query.Offset(offset).Limit(pageSize).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var itemId int
		if err = rows.Scan(&itemId); err != nil {
			err = errors.WithStack(err)
			return
		}
		itemIds = append(itemIds, itemId)
	}
	return
}

func (dao *Dao) GetSearchItemTotal(fields string, where map[string]interface{}, whereIn model.WhereIn, like map[string]string, order string) (total int, err error) {
	query := dao.MasterServiceItems
	if len(fields) > 0 {
		query = query.Select(fields)
	}
	query = query.Where(where)
	if len(whereIn.ItemId) > 0 {
		query = query.Where("item_id in (?)", whereIn.ItemId)
	}
	if len(like) > 0 {
		for k, v := range like {
			query = query.Where(k+" like ?", "%"+v+"%")
		}
	}
	err = query.Model(&model.ItemSearches{}).Count(&total).Error
	return
}

func (dao *Dao) GetItem(itemId int, where map[string]interface{}) (item *model.Items, err error) {
	item = &model.Items{}
	err = dao.MasterServiceItems.
		Table(item.TableName()).
		Where(where).
		Limit(1).
		Find(&item).Error
	return
}

func (dao *Dao) GetSku(skuId int, where map[string]interface{}) (sku *model.ItemSkus, err error) {
	sku = &model.ItemSkus{}
	err = dao.MasterServiceItems.
		Table(sku.TableName()).
		Where(where).
		Limit(1).
		Find(&sku).Error
	return
}

func (dao *Dao) GetSkus(itemId int, where map[string]interface{}) (skus []*model.ItemSkus, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemSkus{}.TableName()).
		Where(where).
		Find(&skus).Error
	return
}

func (dao *Dao) GetPhotos(itemId int, where map[string]interface{}) (photos []*model.ItemPhotos, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemPhotos{}.TableName()).
		Where(where).
		Find(&photos).Error
	return
}

func (dao *Dao) GetParameters(itemId int, where map[string]interface{}) (parameters []*model.ItemParameters, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemParameters{}.TableName()).
		Where(where).
		Find(&parameters).Error
	return
}

func (dao *Dao) GetProps(itemId int, where map[string]interface{}) (props []*model.ItemProps, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemProps{}.TableName()).
		Where(where).
		Find(&props).Error
	return
}

func (dao *Dao) GetPropValues(itemId int, where map[string]interface{}) (propValues []*model.ItemPropValues, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemPropValues{}.TableName()).
		Where(where).
		Find(&propValues).Error
	return
}

func (dao *Dao) InsertItem(item *model.Items) (itemId int, err error) {
	dao.MasterServiceItems.Create(&item)
	if item.ItemId == 0 {
		err = ecode.InsertItemErr
		return
	}
	itemId = item.ItemId
	return
}

func (dao *Dao) InsertSku(sku *model.ItemSkus) (skuId int, err error) {
	dao.MasterServiceItems.Create(&sku)
	if sku.SkuId == 0 {
		err = ecode.InsertItemErr
		return
	}
	skuId = sku.SkuId
	return
}

func (dao *Dao) InsertProp(prop *model.ItemProps) (id int, err error) {
	dao.MasterServiceItems.Create(&prop)
	if prop.Id == 0 {
		err = ecode.InsertPropErr
		return
	}
	id = prop.Id
	return
}

func (dao *Dao) InsertPropValue(propValue *model.ItemPropValues) (id int, err error) {
	dao.MasterServiceItems.Create(&propValue)
	if propValue.Id == 0 {
		err = ecode.InsertPropValueErr
		return
	}
	id = propValue.Id
	return
}

func (dao *Dao) InsertPhoto(photo *model.ItemPhotos) (id int, err error) {
	dao.MasterServiceItems.Create(&photo)
	if photo.Id == 0 {
		err = ecode.InsertPhotoErr
		return
	}
	id = photo.Id
	return
}

func (dao *Dao) InsertParameter(parameter *model.ItemParameters) (id int, err error) {
	dao.MasterServiceItems.Create(&parameter)
	if parameter.Id == 0 {
		err = ecode.InsertParameterErr
		return
	}
	id = parameter.Id
	return
}
