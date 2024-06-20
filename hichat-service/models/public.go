package models

import (
	"encoding/json"
	"time"
)

//var House sync.Map
//var RoomMutexes = make(map[string]*sync.Mutex) //房间锁
//var MutexForRoomMutexes = new(sync.Mutex)      //全局锁

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	PingPeriod = (PongWait * 9) / 10

	MaxMessageSize = int64(1024 * 1024) //消息最大容量
)

var (
	Newline = []byte{'\n'}
	Space   = []byte{' '}
)

// ConversionMsgBytesToMsgStruct 将消息字节转换为结构体
func ConversionMsgBytesToMsgStruct(byte []byte) (MessageTransmitter, error) {
	var groupMsgstruct *Message
	err := json.Unmarshal(byte, &groupMsgstruct)
	if err == nil {
		return groupMsgstruct, nil
	}

	var userMsgStruct *UserMessage
	err = json.Unmarshal(byte, &userMsgStruct)
	if err == nil {
		return userMsgStruct, nil
	}

	return nil, err
}
