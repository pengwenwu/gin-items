package ecode

var MsgMaps = GetAllMsgMaps(
	CommonMsgMaps,
	ItemMsgMaps,
)

func GetAllMsgMaps(msgMaps ...map[Code]string) (allMsgMaps map[Code]string) {
	allMsgMaps = make(map[Code]string)
	for _, msgMap := range msgMaps {
		for k, v := range msgMap {
			allMsgMaps[k] = v
		}
	}
	return
}

func GetMsg(code int) string {
	var intToCode = Int(code)
	msg, ok := MsgMaps[intToCode]
	if ok {
		return msg
	}
	return ""
}
