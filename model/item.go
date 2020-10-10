package model

import (
    "fmt"
    "reflect"
    "strings"

    "github.com/astaxie/beego/validation"
    "gorm.io/gorm"

    "gin-items/helper"
    "gin-items/library/constant"
)

type Items struct {
	ItemId  int    `gorm:"column:item_id;primary_key" json:"item_id"`
	Appkey  string `gorm:"column:appkey" json:"-"`
	Channel int    `gorm:"column:channel" json:"-"`
	Name    string `gorm:"column:name" json:"name"`
	Photo   string `gorm:"column:photo" json:"photo"`
	Detail  string `gorm:"column:detail" json:"detail"`
	State   int    `gorm:"column:state" json:"-"`
	Model
}

type ItemSkus struct {
	SkuId      int    `gorm:"column:sku_id;primary_key" json:"sku_id"`
	ItemId     int    `gorm:"column:item_id" json:"item_id"`
	Appkey     string `gorm:"column:appkey" json:"-"`
	Channel    int    `gorm:"column:channel" json:"-"`
	ItemName   string `gorm:"column:item_name" json:"item_name"`
	SkuName    string `gorm:"column:sku_name" json:"sku_name"`
	SkuPhoto   string `gorm:"column:sku_photo" json:"sku_photo"`
	SkuCode    string `gorm:"column:sku_code" json:"sku_code"`
	BarCode    string `gorm:"column:bar_code" json:"bar_code"`
	Properties string `gorm:"column:properties" json:"properties"`
	State      int    `gorm:"column:state" json:"-"`
	Model
}

type ItemProps struct {
	Id        int    `gorm:"column:id;primary_key" json:"id"`
	ItemId    int    `gorm:"column:item_id" json:"item_id"`
	PropName  string `gorm:"column:prop_name" json:"prop_name"`
	Sort      int    `gorm:"column:sort" json:"sort"`
	HavePhoto int    `gorm:"column:have_photo" json:"have_photo"`
	PropDesc  string `gorm:"column:prop_desc" json:"prop_desc"`
	State     int    `gorm:"column:state" json:"-"`
	Values    []*ItemPropValues `gorm:"-" json:"values"`
	Model
}

type ItemPropValues struct {
	Id            int    `gorm:"column:id;primary_key" json:"id"`
	ItemId        int    `gorm:"column:item_id" json:"item_id"`
	PropName      string `gorm:"column:prop_name" json:"prop_name"`
	PropValueName string `gorm:"column:prop_value_name" json:"prop_value_name"`
	Sort          int    `gorm:"column:sort" json:"sort"`
	PropPhoto     string `gorm:"column:prop_photo" json:"prop_photo"`
	PropDesc      string `gorm:"column:prop_desc" json:"prop_desc"`
	AssistedNum   int    `gorm:"column:assisted_num" json:"assisted_num"`
	State         int    `gorm:"column:state" json:"-"`
	Model
}

type ItemPhotos struct {
	Id     int    `gorm:"column:id;primary_key" json:"id"`
	ItemId int    `gorm:"column:item_id" json:"item_id"`
	Photo  string `gorm:"column:photo" json:"photo"`
	Sort   int    `gorm:"column:sort" json:"sort"`
	State  int    `gorm:"column:state" json:"-"`
	Model
}

type ItemParameters struct {
	Id         int    `gorm:"column:id;primary_key" json:"id"`
	ItemId     int    `gorm:"column:item_id" json:"item_id"`
	Parameters string `gorm:"column:parameters" json:"parameters"`
	Value      string `gorm:"column:value" json:"value"`
	State      int    `gorm:"column:state" json:"-"`
	Sort       int    `gorm:"column:sort" json:"sort"`
	Model
}

type ItemSearches struct {
	Id        int    `gorm:"column:id;primary_key" json:"id"`
	Appkey    string `gorm:"column:appkey" json:"-"`
	Channel   int    `gorm:"column:channel" json:"-"`
	ItemId    int    `gorm:"column:item_id" json:"item_id"`
	SkuId     int    `gorm:"column:sku_id" json:"sku_id"`
	SkuName   string `gorm:"column:sku_name" json:"sku_name"`
	BarCode   string `gorm:"column:bar_code" json:"bar_code"`
	SkuCode   string `gorm:"column:sku_code" json:"sku_code"`
	ItemState int    `gorm:"column:item_state" json:"-"`
	SkuState  int    `gorm:"column:sku_state" json:"-"`
	Model
}

type Item struct {
	*Items
	Photos     []*ItemPhotos `json:"photos"`
	Parameters []*ItemParameters `json:"parameters"`
	Skus       []*ItemSkus `json:"skus"`
	Props      []*ItemProps `json:"props"`
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

func BeforeCreate(db *gorm.DB) {
    fmt.Println(db.Statement.ReflectValue.Kind())
    if db.Statement.Schema != nil {
        switch db.Statement.ReflectValue.Kind() {
        case reflect.Slice, reflect.Array:
            for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
                fmt.Println(111, db.Statement.ReflectValue.Index(i), db.Statement.ReflectValue.Index(i).Len())
                t := db.Statement.ReflectValue.Index(i).FieldByName("dated")
                fmt.Println(222, t.IsNil())
            }
        case reflect.Struct:
            if field := db.Statement.Schema.LookUpField("dated"); field != nil {
                _ = field.Set(db.Statement.ReflectValue, helper.FormatDateTimeNow())
            }
            if field := db.Statement.Schema.LookUpField("last_dated"); field != nil {
                _ = field.Set(db.Statement.ReflectValue, helper.FormatDateTimeZero())
            }
        }
    }
}

func (items *Items) BeforeCreate(tx *gorm.DB) error {
	items.Dated = helper.FormatDateTimeNow()
	items.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (items *Items) BeforeUpdate(tx *gorm.DB) error {
	items.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (sku *ItemSkus) BeforeCreate(tx *gorm.DB) error {
	sku.Dated = helper.FormatDateTimeNow()
	sku.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (sku *ItemSkus) BeforeUpdate(tx *gorm.DB) error {
	sku.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (prop *ItemProps) BeforeCreate(tx *gorm.DB) error {
	prop.Dated = helper.FormatDateTimeNow()
	prop.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (prop *ItemProps) BeforeUpdate(tx *gorm.DB) error {
	prop.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (propValue *ItemPropValues) BeforeCreate(tx *gorm.DB) error {
	propValue.Dated = helper.FormatDateTimeNow()
	propValue.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (propValue *ItemPropValues) BeforeUpdate(tx *gorm.DB) error {
	propValue.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (photo *ItemPhotos) BeforeCreate(tx *gorm.DB) error {
	photo.Dated = helper.FormatDateTimeNow()
	photo.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (photo *ItemPhotos) BeforeUpdate(tx *gorm.DB) error {
	photo.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (parameter *ItemParameters) BeforeCreate(tx *gorm.DB) error {
	parameter.Dated = helper.FormatDateTimeNow()
	parameter.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (parameter *ItemParameters) BeforeUpdate(tx *gorm.DB) error {
	parameter.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (search *ItemSearches) BeforeCreate(tx *gorm.DB) error {
	search.Dated = helper.FormatDateTimeNow()
	search.LastDated = helper.FormatDateTimeZero()
	return nil
}
func (search *ItemSearches) BeforeUpdate(tx *gorm.DB) error {
	search.LastDated = helper.FormatDateTimeNow()
	return nil
}

func (item *Item) Valid(v *validation.Validation) {
	v.Required(item.Appkey, "appkey")
	v.Required(item.Channel, "channel")
	v.Required(item.Name, "name")

	var propValues []*ItemPropValues
	if len(item.Props) > 0 {
		for _, prop := range item.Props {
			v.Required(prop.PropName, "props.prop_name")
			v.Required(prop.Values, "props.values")
			for _, propValue := range prop.Values {
				v.Required(propValue.PropValueName, "props.values.prop_value_name")
				propValue.PropName = prop.PropName
				propValues = append(propValues, propValue)
			}
		}
	}


	if len(item.Skus) > 0 {
		for _, sku := range item.Skus {
			if sku.SkuName == "" && sku.Properties == "" {
				_ = v.SetError("skus.properties", "缺少sku_name")
			}
			skuName := item.Name
			if sku.Properties != "" {
				skuProps := strings.Split(sku.Properties, ";")
				for _, v := range skuProps {
					skuProp := strings.Split(v, ":")
					skuName += " " + skuProp[1]

					for _, propValue := range propValues {
						if propValue.PropName == skuProp[0] && propValue.PropValueName == skuProp[1] && propValue.PropPhoto != "" {
							sku.SkuPhoto = propValue.PropPhoto
						}
					}
				}
			}

			if sku.SkuName == "" {
				sku.SkuName = skuName
			}
			sku.Appkey = item.Appkey
			sku.Channel = item.Channel
			sku.ItemName = item.Name
			sku.State = constant.ItemSkuStateNormal
		}
	}

	if len(item.Photos) > 0 {
		for _, photo := range item.Photos {
			v.Required(photo.Photo, "photos.photo")
		}
	}

	if len(item.Parameters) > 0 {
		for _, parameter := range item.Parameters {
			v.Required(parameter.Value, "parameters.parameters")
			v.Required(parameter.Value, "parameters.value")
		}
	}

	// 当主图没有时，轮播图第一张图设置为主图
	if len(item.Photos) > 0 && item.Photo == "" {
		item.Photo = item.Photos[0].Photo
	}
	// 当没有轮播图、主图的时候，设置默认图
	if len(item.Photos) == 0 && item.Photo == "" {
		item.Photo = constant.ItemDefaultPhoto
	}
	// 当没有轮播图的时候，选设置的第一张默认图
	if len(item.Photos) == 0 {
		item.Photos = append(item.Photos, &ItemPhotos{
			Photo: item.Photo,
		})
	}
}
