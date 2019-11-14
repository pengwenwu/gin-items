package e

var MsgMaps = map[int]string{
	Success:       "ok",
	Error:         " fail",
	InvalidParams: "请求参数错误",

	ErrorGetItemListFail: "获取商品列表失败",
	ErrorGetItemCount: "获取商品列表总数失败",
}

func GetMsg(code int) string {
	msg, ok := MsgMaps[code]
	if ok {
		return msg
	}
	return MsgMaps[Error]
}
