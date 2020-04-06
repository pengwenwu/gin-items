package ecode

// 为防止定义重复，以及方便查找，各个服务错误码可直接选用一个段号
// item ecode interval is [1000, 1999]

var (
	ItemIllegalItemId = New(1000)
)
