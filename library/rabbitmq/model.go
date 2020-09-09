package rabbitmq

type SyncSkuInsertData struct {
	ItemId int
	SkuId  int
}

type SyncSkuUpdateData struct {
	ItemId int
	SkuId  int
}

type SyncItemSearchesData struct {
	ItemId int
	SkuId int
}
