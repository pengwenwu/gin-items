package ecode

var MsgMaps = map[int]string{
	Success:              "ok",
	Error:                "fail",
	InvalidParams:        "请求参数错误",
	UnsupportedMediaType: "错误的请求格式",
	ErrorGetItemListFail: "获取商品item列表失败",
	ErrorGetItemFail:     "获取商品item失败",
}

func GetMsg(code int) string {
	msg, ok := MsgMaps[code]
	if ok {
		return msg
	}
	return MsgMaps[Error]
}
