package define

// items 商品表
const (
	ItemStateDeleted     = iota // 删除（放入回收站） 0
	ItemStateNormal             // 正常
	ItemStateDeletedReal        // 彻底删除
)

// item_skus 商品sku表
const (
	ItemSkuStateDeleted     = iota // 删除（放入回收站）
	ItemSkuStateNormal             // 正常
	ItemSkuStateDeletedReal        // 彻底删除
	ItemSkuStateDeletedSelf        // 业务上本身删除
)

// item_props 商品规格表
const (
	ItemPropsStateDeleted = iota // 已删除
	ItemPropsStateNormal         // 正常
)

// item_prop_values 商品规格对应的属性表
const (
	ItemPropsValuesStateDeleted = iota // 已删除
	ItemPropsValuesStateNormal         // 正常
)

// item_parameters 商品产品属性表
const (
	ItemParametersStateDeleted = iota // 已删除
	ItemParametersStateNormal         // 正常
)

// item_photos 商品轮播图表
const (
	ItemPhotosStateDeleted = iota // 已删除
	ItemPhotosStateNormal         // 正常
)
