package dao

import (
	"github.com/pkg/errors"

	"gin-items/model"
)

func (dao *Dao) GetSearchItemIds(params model.ArgItemSearch, fields string, where map[string]interface{}, offset, limit int) (itemIds []int, err error) {
	rows, err := dao.DB.Table(params.TableName()).Where(where).Offset(offset).Limit(limit).Rows()
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

func (dao *Dao) GetSearchItemTotal(where interface{}) (int, error) {
	var count int
	err := dao.DB.Model(&model.ItemSearches{}).Where(where).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dao *Dao) GetItem(itemId int, fields string) (item map[string]string, err error) {
	rows, err := dao.DB.
		Table(model.Items{}.TableName()).
		Select(fields).
		Where("item_id = ?", itemId).
		Rows()
	if err != nil {
		return
	}
	data := RowsToSliceMap(rows)
	if len(data) > 0 {
		item = data[0]
	}
	return
}

func (dao *Dao) GetItemPhotos(fields string, where map[string]interface{}, order string, page, pageSize int) (data []map[string]string, err error) {
	offSet := (page - 1) * pageSize
	rows, err := dao.DB.
		Table(model.ItemPhotos{}.TableName()).
		Select(fields).
		Where(where).
		Order(order).
		Offset(offSet).
		Limit(pageSize).
		Rows()
	if err != nil {
		return
	}
	data = RowsToSliceMap(rows)
	return
}

func (dao *Dao) GetItemParameters(fields string, where map[string]interface{}, order string, page, pageSize int) (data []map[string]string, err error) {
	offSet := (page - 1) * pageSize
	rows, err := dao.DB.
		Table(model.ItemParameters{}.TableName()).
		Select(fields).
		Where(where).
		Order(order).
		Offset(offSet).
		Limit(pageSize).
		Rows()
	if err != nil {
		return
	}
	data = RowsToSliceMap(rows)
	return
}

func (dao *Dao) GetItemSkus(fields string, where map[string]interface{}, order string, page, pageSize int) (data []map[string]string, err error) {
	offSet := (page - 1) * pageSize
	rows, err := dao.DB.
		Table(model.ItemSkus{}.TableName()).
		Select(fields).
		Where(where).
		Order(order).
		Offset(offSet).
		Limit(pageSize).
		Rows()
	if err != nil {
		return
	}
	data = RowsToSliceMap(rows)
	return
}

func (dao *Dao) GetItemProps(where map[string]interface{}, order string, page, pageSize int) (data []model.ItemProps, err error) {
	offSet := (page - 1) * pageSize
	err = dao.DB.
		Table(model.ItemProps{}.TableName()).
		Where(where).
		Offset(offSet).
		Limit(pageSize).
		Find(&data).
		Error
	if err != nil {
		return
	}
	return
}

func (dao *Dao) GetItemPropValues(where map[string]interface{}, order string, page, pageSize int) (data []model.ItemPropValues, err error) {
	offSet := (page - 1) * pageSize
	err = dao.DB.
		Table(model.ItemPropValues{}.TableName()).
		Where(where).
		Offset(offSet).
		Limit(pageSize).
		Find(&data).
		Error
	if err != nil {
		return
	}
	return
}
