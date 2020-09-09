package rabbitmq

type RouteKey string

var (
	ItemInsert = RouteKey("service_item.syncItemInsert")
	ItemUpdate = RouteKey("service_item.syncItemUpdate")
	SkuInsert  = RouteKey("service_item.syncSkuInsert")
	SkuUpdate  = RouteKey("service_item.syncSkuUpdate")

	ItemSearches = RouteKey("service_item.sync")
)
