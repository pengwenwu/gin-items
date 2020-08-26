package ecode

var CommonMsgMaps = map[Code]string{
	OK:               "success",
	IllegalParams:    "参数非法",
	NoAuthorization:  "未获取到token",
	AuthorizationErr: "token非法",
	NoAppKey:         "缺少appKey",
	RecordNotFound:	  "未查询到记录",
	RequestErr:       "请求错误",
	ServerErr:        "服务器错误",
}
