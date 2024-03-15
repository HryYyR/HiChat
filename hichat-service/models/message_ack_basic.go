package models

type AckMessage struct {
	MsgSort   int `json:"MsgSort"`   //1好友 2群聊
	MsgType   int `json:"MsgType"`   //888
	AckStatus int `json:"AckStatus"` //0失败 1成功
	UserId    int `json:"UserId"`    //发送人的id
}
