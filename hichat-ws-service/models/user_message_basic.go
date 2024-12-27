package models

import (
	"encoding/json"
	"log"
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
	Msg               string `xorm:"varchar(2048) notnull"`
	MsgType           int
	IsReply           bool //是否是回复消息
	ReplyMsgID        int  //如果是,被回复的消息的id
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

	sendClientList, ok1 := ServiceCenter.Clients[u.UserID]
	ReceiveClientList, ok2 := ServiceCenter.Clients[u.ReceiveUserID]
	//1501为视频通话,只需传给接收方
	if ok1 && u.MsgType < 1500 {
		for i, client := range sendClientList {
			if client.Status {
				log.Println("发送给send用户", u.UserID)
				ServiceCenter.Clients[u.UserID][i].Send <- bytes
			} else {
				log.Println("send用户", u.UserID, "不在线")

			}
		}

	}
	if ok2 {
		for i, client := range ReceiveClientList {
			if client.Status {
				log.Println("发送给Receive用户", u.ReceiveUserID)
				ServiceCenter.Clients[u.ReceiveUserID][i].Send <- bytes
			} else {
				log.Println("Receive用户", u.ReceiveUserID, "不在线")

			}
		}

	}

	return nil
}
