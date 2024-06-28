package models

import (
	"encoding/json"
	"time"
)

// UserMessage 数据库的用户消息
type UserMessage struct {
	ID                int    `xorm:"pk autoincr index"`
	UserID            int    `xorm:"notnull"`
	UserName          string `xorm:"notnull"`
	UserAvatar        string
	ReceiveUserID     int    `xorm:"notnull"`
	ReceiveUserName   string `xorm:"notnull"`
	ReceiveUserAvatar string
	Msg               string `xorm:"notnull"`
	MsgType           int
	IsReply           bool //是否是回复消息
	ReplyUserID       int  //如果是,被回复的用户id
	Context           []byte
	CreatedAt         time.Time `xorm:"created"`
	DeletedAt         time.Time `xorm:"deleted"`
	UpdatedAt         time.Time `xorm:"updated"`
}

func (u UserMessage) Transmit() error {
	bytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	client1, ok1 := ServiceCenter.Clients[u.UserID]
	client2, ok2 := ServiceCenter.Clients[u.ReceiveUserID]
	//1501为视频通话,只需传给接收方
	if ok1 && u.MsgType < 1500 {
		if client1.Status {
			//log.Println("发送给用户", u.UserID)
			ServiceCenter.Clients[u.UserID].Send <- bytes
		}
	}
	if ok2 {
		if client2.Status {
			//log.Println("发送给用户", u.ReceiveUserID)
			ServiceCenter.Clients[u.ReceiveUserID].Send <- bytes
		}
	}

	return nil
}
