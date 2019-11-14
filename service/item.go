package service

import (
	"gin-items/model"
)

type Items struct {
	ItemId int
	OffSet int
	PageSize int
}

func (i *Items) GetItemList() ([]*model.Items, error) {
	var items []*model.Items

	items, err := model.GetItemList(i.OffSet, i.PageSize, i.getMaps())
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *Items) Count() (int, error)  {
	return model.GetItemTotal(i.getMaps())
}

func (i *Items) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	return maps
}