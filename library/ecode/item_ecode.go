package ecode

// 为防止定义重复，以及方便查找，各个服务错误码可直接选用一个段号
// item ecode interval is [1000, 1999]

var (
	InsertItemErr = New(1000) // 新增item失败
	InsertSkuErr = New(1001) // 新增sku失败
	InsertPropErr = New(1002) // 新增prop失败
	InsertPropValueErr = New(1003) // 新增propValue失败
	InsertPhotoErr = New(1004) // 新增photo失败
	InsertParameterErr = New(1005) // 新增parameter失败
)
