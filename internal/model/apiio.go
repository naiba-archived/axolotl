package model

type Response struct {
	Code uint        `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
type RunCodeRequest struct {
	Container string `json:"container,omitempty"`
	Code      string `json:"code,omitempty"`
}
