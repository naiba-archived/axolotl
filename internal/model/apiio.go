package model

const (
	MsgTypePeer                      = iota // 建立连接
	MsgTypeChooseProgrammingLanguage        // 选择编程语言
	MsgTypeExecResult                       // 执行结果
	MsgTypeCreatePeer                       // 新成员加入，创建 Peer
	MsgTypeInitRoom                         // 初始化会议室、创建 Peer 与 选择语言
)

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

type WsMsg struct {
	Type uint        `json:"type"`
	From string      `json:"from"`
	Data interface{} `json:"data"`
}
