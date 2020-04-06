package model

// ParamValidator .
type ParamValidator interface {
	Validate() bool
}

// ArgItemSearch param struct
type ArgItemSearch struct {
	ItemId    int                      `json:"item_id"`
	SkuId     int                      `json:"sku_id"`
	BarCode   string                   `json:"bar_code"`
	SkuCode   string                   `json:"sku_code"`
	ItemState int                      `json:"item_state"`
	SkuState  int                      `json:"sku_state"`
	Name      string                   `json:"name"`
	WhereIn   WhereIn                  `json:"where_in"`
	Like      map[string][]interface{} `json:"like"`
	Fields    string                   `json:"fields"`
	Page      int                      `json:"page"`
	PageSize  int                      `json:"page_size"`
	Order     string                   `json:"order"`
	Desc      string                   `json:"desc"`
}

type WhereIn struct {
	ItemId []int `json:"item_id"`
}

type ArgGetItemById struct {
	Fields string `json:"fields"`
}

// Validate .
func (a ArgItemSearch) Validate() bool {
	return true
}

func (a ArgItemSearch) GetWhereMap() (whereMap map[string]interface{}) {
	whereMap = make(map[string]interface{})
	if a.BarCode != "" {
		whereMap["bar_code"] = a.BarCode
	}
	if a.Name != "" {
		whereMap["sku_name"] = a.Name
	}
	return
}

func (a ArgItemSearch) TableName() string {
	return "item_searches"
}
