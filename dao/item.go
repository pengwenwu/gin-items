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

func (dao *Dao) GetItem(where map[string]interface{}) (item *model.Items, err error) {
	item = &model.Items{}
	err = dao.MasterServiceItems.
		Table(item.TableName()).
		Where(where).
		Limit(1).
		Find(&item).Error
	return
}

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

func (dao *Dao) GetPhotos(where map[string]interface{}) (photos []*model.ItemPhotos, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemPhotos{}.TableName()).
		Where(where).
		Find(&photos).Error
	return
}

func (dao *Dao) GetParameters(where map[string]interface{}) (parameters []*model.ItemParameters, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemParameters{}.TableName()).
		Where(where).
		Find(&parameters).Error
	return
}

func (dao *Dao) GetProps(where map[string]interface{}) (props []*model.ItemProps, err error) {
	err = dao.MasterServiceItems.
		Table(model.ItemProps{}.TableName()).
		Where(where).
		Find(&props).Error
	return
}

func (dao *Dao) GetPropValues(where map[string]interface{}) (propValues []*model.ItemPropValues, err error) {
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

func (dao *Dao) InsertSku(sku *model.ItemSkus) error {
	dao.MasterServiceItems.Create(&sku)
	if sku.SkuId == 0 {
		err := ecode.InsertItemErr
		return err
	}
	return nil
}

func (dao *Dao) InsertProp(prop *model.ItemProps) error {
	dao.MasterServiceItems.Create(&prop)
	if prop.Id == 0 {
		err := ecode.InsertPropErr
		return err
	}
	return nil
}

func (dao *Dao) InsertPropValue(propValue *model.ItemPropValues) error {
	dao.MasterServiceItems.Create(&propValue)
	if propValue.Id == 0 {
		err := ecode.InsertPropValueErr
		return err
	}
	return nil
}

func (dao *Dao) InsertPhoto(photo *model.ItemPhotos) error {
	dao.MasterServiceItems.Create(&photo)
	if photo.Id == 0 {
		err := ecode.InsertPhotoErr
		return err
	}
	return nil
}

func (dao *Dao) InsertParameter(parameter *model.ItemParameters) error {
	dao.MasterServiceItems.Create(&parameter)
	if parameter.Id == 0 {
		err := ecode.InsertParameterErr
		return err
	}
	return nil
}

func (dao *Dao) InsertSearches(search *model.ItemSearches) error {
	dao.MasterServiceItems.Create(&search)
	if search.Id == 0 {
		err := ecode.InsertSearchErr
		return err
	}
	return nil
}

func (dao *Dao) UpdateItem(item *model.Items, where map[string]interface{}) error {
	return dao.MasterServiceItems.
		Model(&item).
		Debug().
		Select("name").
		Where(where).
		Limit(1).
		Updates(&item).
		Error
}
