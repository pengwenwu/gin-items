package define

// items 商品表
const (
	ITEM_STATE_DELETED      = iota // 删除（放入回收站） 0
	ITEM_STATE_NORMAL              // 正常
	ITEM_STATE_DELETED_REAL        // 彻底删除
)

// item_skus 商品sku表
const (
	ITEM_SKU_STATE_DELETED      = iota // 删除（放入回收站）
	ITEM_SKU_STATE_NORMAL              // 正常
	ITEM_SKU_STATE_DELETED_REAL        // 彻底删除
	ITEM_SKU_STATE_DELETED_SELF        // 业务上本身删除
)

// item_props 商品规格表
const (
	ITEM_PROPS_STATE_DELETED = iota // 已删除
	ITEM_PROPS_STATE_NORMAL         // 正常
)

// item_prop_values 商品规格对应的属性表
const (
	ITEM_PROPS_VALUES_STATE_DELETED = iota // 已删除
	ITEM_PROPS_VALUES_STATE_NORMAL         // 正常
)

// item_parameters 商品产品属性表
const (
	ITEM_PARAMETERS_STATE_DELETED = iota // 已删除
	ITEM_PARAMETERS_STATE_NORMAL         // 正常
)

// item_photos 商品轮播图表
const (
	ITEM_PHOTOS_STATE_DELETED = iota // 已删除
	ITEM_PHOTOS_STATE_NORMAL         // 正常
)
