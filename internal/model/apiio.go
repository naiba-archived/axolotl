package model

type Response struct {
	Code uint        `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
type RunCodeRequest struct {
	Room      string `json:"room,omitempty"`
	Container string `json:"container,omitempty"`
	Code      string `json:"code,omitempty"`
}

const (
	MsgTypePeer                      = iota // 建立连接
	MsgTypeChooseProgrammingLanguage        // 选择编程语言
	MsgTypeExecResult                       // 执行结果
)

type WsMsg struct {
	Type uint        `json:"type"`
	Data interface{} `json:"data"`
}
