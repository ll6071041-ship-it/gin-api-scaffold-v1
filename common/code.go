package common

// ResCode 自定义业务状态码
type ResCode int64

const (
	CodeSuccess         ResCode = 1000 + iota
	CodeInvalidParam            // 1001
	CodeUserExist               // 1002
	CodeUserNotExist            // 1003
	CodeInvalidPassword         // 1004
	CodeServerBusy              // 1005

	CodeNeedLogin
	CodeInvalidToken
)

// codeMsgMap 状态码映射
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的Token",
}

// Msg 方法：获取状态码对应的提示信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
