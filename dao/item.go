package dao

import (
	"github.com/pkg/errors"

	"gin-items/model"
)

func (dao *Dao) GetSearchItemIds(fields string, where map[string]interface{}, whereIn model.WhereIn, like map[string]string, order, groupBy string, page, pageSize int) (itemIds []int, err error) {
	offset := (page - 1) * pageSize
	query := dao.DB.
		Table(model.ItemSearches{}.TableName()).
		Select(fields).
		Where(where)
	if len(whereIn.ItemId) > 0 {
		query = query.Where("item_id in (?)", whereIn.ItemId)
	}
	if len(like) > 0 {
		for k,v := range like {
			query = query.Where(k + " like ?", "%" + v +"%")
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
	query := dao.DB
	if len(fields) > 0 {
		query =query.Select(fields)
	}
	query = query.Where(where)
	if len(whereIn.ItemId) > 0 {
		query = query.Where("item_id in (?)", whereIn.ItemId)
	}
	if len(like) > 0 {
		for k,v := range like {
			query = query.Where(k + " like ?", "%" + v +"%")
		}
	}
	err = query.Model(&model.ItemSearches{}).Count(&total).Error
	if err != nil {
		return
	}
	return
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
