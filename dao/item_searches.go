package dao

import (
	"fmt"
	"gin-items/library/ecode"
	"gin-items/model"
	"github.com/pkg/errors"
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

func (dao *Dao) InsertSearch(search *model.ItemSearches) error {
	dao.MasterServiceItems.Create(&search)
	if search.Id == 0 {
		err := ecode.InsertSearchErr
		return err
	}
	return nil
}

func (dao *Dao) InsertSearches(searchList []*model.ItemSearches) error {
	for _, v := range searchList {
		fmt.Printf("%+v\n", v)
	}
	return dao.MasterServiceItems.Debug().Create(&searchList).Error
}

func (dao *Dao) UpdateSearch(where, update map[string]interface{}) error {
	return dao.MasterServiceItems.
		Model(&model.ItemSearches{}).
		Where(where).
		Limit(updateCommonLimit).
		Updates(update).
		Error
}