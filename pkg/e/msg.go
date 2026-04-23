package e

var MsgFlagMap = map[int]string{
	SUCCESS:                 "ok",
	ERROR:                   "fail",
	InvalidParams:           "请求参数错误",
	ErrorAuthCheckTokenFail: "Token鉴权失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlagMap[code]
	if ok {
		return msg
	}
	return MsgFlagMap[ERROR]
}
