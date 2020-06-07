package model

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

type Items struct {
	ItemId  int    `gorm:"column:item_id" json:"item_id"`
	Appkey  string `gorm:"column:appkey" json:"appkey"`
	Channel int    `gorm:"column:channel" json:"channel"`
	Name    string `gorm:"column:name" json:"name"`
	Photo   string `gorm:"column:photo" json:"photo,omitempty"`
	Detail  string `gorm:"column:detail" json:"detail,omitempty"`
	State   int    `gorm:"column:state" json:"state,omitempty"`
	Model
}

type ItemSkus struct {
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
	Model
}

type ItemProps struct {
	Id        int    `gorm:"column:id" json:"id"`
	ItemId    int    `gorm:"column:item_id" json:"item_id"`
	PropName  string `gorm:"column:prop_name" json:"prop_name"`
	Sort      int    `gorm:"column:sort" json:"sort"`
	HavePhoto int    `gorm:"column:have_photo" json:"have_photo"`
	PropDesc  string `gorm:"column:prop_desc" json:"prop_desc"`
	State     int    `gorm:"column:state" json:"state"`
	Values    []ItemPropValues `json:"values"`
	Model
}

type ItemPropValues struct {
	Id            int    `gorm:"column:id" json:"id"`
	ItemId        int    `gorm:"column:item_id" json:"item_id"`
	PropName      string `gorm:"column:prop_name" json:"prop_name"`
	PropValueName string `gorm:"column:prop_value_name" json:"prop_value_name"`
	Sort          int    `gorm:"column:sort" json:"sort"`
	PropPhoto     string `gorm:"column:prop_photo" json:"prop_photo"`
	PropDesc      string `gorm:"column:prop_desc" json:"prop_desc"`
	AssistedNum   int    `gorm:"column:assisted_num" json:"assisted_num"`
	State         int    `gorm:"column:state" json:"state"`
	Model
}

type ItemPhotos struct {
	Id     int    `gorm:"column:id" json:"id"`
	ItemId int    `gorm:"column:item_id" json:"item_id"`
	Photo  string `gorm:"column:photo" json:"photo"`
	Sort   int    `gorm:"column:sort" json:"sort"`
	State  int    `gorm:"column:state" json:"state"`
	Model
}

type ItemParameters struct {
	Id         int    `gorm:"column:id" json:"id"`
	ItemId     int    `gorm:"column:item_id" json:"item_id"`
	Parameters string `gorm:"column:parameters" json:"parameters"`
	Value      string `gorm:"column:value" json:"value"`
	State      int    `gorm:"column:state" json:"state"`
	Sort       int    `gorm:"column:sort" json:"sort"`
	Model
}

type ItemSearches struct {
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
	Model
}

type Item struct {
	Items
	Photos     []*ItemPhotos `json:"photos,omitempty"`
	Parameters []*ItemParameters `json:"parameters,omitempty"`
	Skus       []*ItemSkus `json:"skus,omitempty"`
	Props      []*ItemProps `json:"props,omitempty"`
}

func (Items) TableName() string {
	return "items"
}

func (ItemSkus) TableName() string {
	return "item_skus"
}

func (ItemSearches) TableName() string {
	return "item_searches"
}

func (ItemProps) TableName() string {
	return "item_props"
}

func (ItemPropValues) TableName() string {
	return "item_prop_values"
}

func (ItemPhotos) TableName() string {
	return "item_photos"
}

func (ItemParameters) TableName() string {
	return "item_parameters"
}

func GetFields(i interface{}) (fields []string) {
	t:=reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i:=0;i<t.NumField();i++{
		if v.Field(i).Type().Kind() == reflect.Struct{
			structField := v.Field(i).Type()
			for j :=0 ; j< structField.NumField(); j++ {
				fields = append(fields, structField.Field(j).Tag.Get("json"))
			}
			continue
		}
		sf:=t.Field(i)
		fields = append(fields, sf.Tag.Get("json"))
	}
	return
}

func (items *Items) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("dated", time.Now().Format("2006-01-02 15:04:05"))
	scope.SetColumn("last_dated", "0000-00-00 00:00:00")
	return nil
}
func (items *Items) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("last_dated", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}