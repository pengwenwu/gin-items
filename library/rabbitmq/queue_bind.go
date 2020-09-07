package rabbitmq

type queueBind struct {
	name string
	keys []RouteKey
}

var (
	SyncItemInsert = &queueBind{
		name: "service_item.syncItemInsert",
		keys: []RouteKey{ItemInsert},
	}
	SyncItemUpdate = &queueBind{
		name: "service_item.syncItemUpdate",
		keys: []RouteKey{ItemUpdate},
	}
	SyncSkuInsert = &queueBind{
		name: "service_item.syncSkuInsert",
		keys: []RouteKey{SkuInsert},
	}
	SyncSkuUpdate = &queueBind{
		name: "service_item.syncSkuUpdate",
		keys: []RouteKey{SkuUpdate},
	}
)
