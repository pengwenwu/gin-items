package model

// ParamValidator .
type ParamValidator interface {
	Validate() bool
}

// ArgItemSearch param struct
type ArgItemSearch struct {
	ItemId    int                      `form:"item_id"`
	SkuId     int                      `form:"sku_id"`
	BarCode   string                   `form:"bar_code"`
	SkuCode   string                   `form:"sku_code"`
	ItemState int                      `form:"item_state"`
	SkuState  int                      `form:"sku_state"`
	Name      string                   `form:"name"`
	WhereIn   map[string][]interface{} `form:"where_in"`
	Like      map[string][]interface{} `form:"like"`
	Fields    string                   `form:"fields"`
	Page      int                      `form:"page" default:"1"`
	PageSize  int                      `form:"page_size" default:"20"`
	Order     string                   `form:"order" default:"item_id"`
	Desc      string                   `form:"desc" default:"desc"`
}

// Validate .
func (p *ArgItemSearch) Validate() bool {
	return true
}
