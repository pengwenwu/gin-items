package rabbitmq

type SyncSkuInsertData struct {
	ItemId int
	SkuId  int
}

type SyncSkuUpdateData struct {
	ItemId int
	SkuId  int
}

type SyncItemInsertData struct {
	ItemId int
}

type SyncItemUpdateData struct {
	ItemId int
}
