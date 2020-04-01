package model

type Items struct {
	Model

	ItemId  int    `gorm:"column:item_id" json:"item_id"`
	Appkey  string `gorm:"column:appkey" json:"appkey"`
	Channel int    `gorm:"column:channel" json:"channel"`
	Name    string `gorm:"column:name" json:"name"`
	Photo   string `gorm:"column:photo" json:"photo"`
	Detail  string `gorm:"column:detail" json:"detail"`
	State   int    `gorm:"column:state" json:"state"`
}

type ItemSkus struct {
	Model

	SkuId      int    `gorm:"column:sku_id" json:"sku_id"`
	ItemId     int    `gorm:"column:item_id" json:"item_id"`
	Appkey     string `gorm:"column:appkey" json:"appkey"`
	Channel    int    `gorm:"column:channel" json:"channel"`
	ItemName   string `gorm:"column:item_name" json:"item_name"`
	SkuName    string `gorm:"column:sku_name" json:"sku_name"`
	SkuPhoto   string `gorm:"column:sku_photo" json:"sku_photo"`
	SkuCode    string `gorm:"column:sku_code" json:"sku_code"`
	BarCode    string `gorm:"column:bar_code" json:"bar_code"`
	Properties string `gorm:"column:properties" json:"properties"`
	State      int    `gorm:"column:state" json:"state"`
}

type ItemProps struct {
	Model

	Id        int    `gorm:"column:id" json:"id"`
	ItemId    int    `gorm:"column:item_id" json:"item_id"`
	PropName  string `gorm:"column:prop_name" json:"prop_name"`
	Sort      int    `gorm:"column:sort" json:"sort"`
	HavePhoto int    `gorm:"column:have_photo" json:"have_photo"`
	PropDesc  string `gorm:"column:prop_desc" json:"prop_desc"`
	State     int    `gorm:"column:state" json:"state"`
}

type ItemPropValue struct {
	Model

	Id            int    `gorm:"column:id" json:"id"`
	ItemId        int    `gorm:"column:item_id" json:"item_id"`
	PropName      string `gorm:"column:prop_name" json:"prop_name"`
	PropValueName string `gorm:"column:prop_value_name" json:"prop_value_name"`
	Sort          int    `gorm:"column:sort" json:"sort"`
	PropPhoto     string `gorm:"column:prop_photo" json:"prop_photo"`
	PropDesc      string `gorm:"column:prop_desc" json:"prop_desc"`
	AssistedNum   int    `gorm:"column:assisted_num" json:"assisted_num"`
	State         int    `gorm:"column:state" json:"state"`
}

type ItemPhotos struct {
	Model

	Id     int    `gorm:"column:id" json:"id"`
	ItemId int    `gorm:"column:item_id" json:"item_id"`
	Photo  string `gorm:"type:varchar(255)"`
	Sort   int    `gorm:"column:sort" json:"sort"`
	State  int    `gorm:"column:state" json:"state"`
}

type ItemParameters struct {
	Model

	Id         int    `gorm:"column:id" json:"id"`
	ItemId     int    `gorm:"column:item_id" json:"item_id"`
	Parameters string `gorm:"column:parameters" json:"parameters"`
	Value      string `gorm:"column:value" json:"value"`
	State      int    `gorm:"column:state" json:"state"`
	Sort       int    `gorm:"column:sort" json:"sort"`
}

type ItemSearches struct {
	Model

	Id        int    `gorm:"column:id" json:"id"`
	Appkey    string `gorm:"column:appkey" json:"appkey"`
	Channel   int    `gorm:"column:channel" json:"channel"`
	ItemId    int    `gorm:"column:item_id" json:"item_id"`
	SkuId     int    `gorm:"column:sku_id" json:"sku_id"`
	SkuName   string `gorm:"column:sku_name" json:"sku_name"`
	BarCode   string `gorm:"column:bar_code" json:"bar_code"`
	SkuCode   string `gorm:"column:sku_code" json:"sku_code"`
	ItemState int    `gorm:"column:item_state" json:"item_state"`
	SkuState  int    `gorm:"column:sku_state" json:"sku_state"`
}
