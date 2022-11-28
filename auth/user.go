package auth

type UserInfo struct {
	Uid   int64       `json:"uid"`
	Token string      `json:"token"`
	Info  interface{} `json:"info"`
}

type ValidateMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var noLoginMsg = &ValidateMsg{Code: 403, Message: "未登录"}
