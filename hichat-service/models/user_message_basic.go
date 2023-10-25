package models

import "time"

type UserMessage struct {
	ID            int    `xorm:"pk autoincr index"`
	UUID          string `xorm:"notnull"`
	UserID        int    `xorm:"notnull"`
	UserName      string `xorm:"notnull"`
	UserAvatar    string
	PreUserID     int    `xorm:"notnull"`
	PreUserName   string `xorm:"notnull"`
	PreUserAvatar string
	Msg           string `xorm:"notnull"`
	MsgType       int
	IsReply       bool //是否是回复消息
	ReplyUserID   int  //如果是,被回复的用户id
	Context       []byte
	CreatedAt     time.Time `xorm:"created"`
	DeletedAt     time.Time `xorm:"deleted"`
	UpdatedAt     time.Time `xorm:"updated"`
}
