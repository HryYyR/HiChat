package models

import "time"

// 群聊消息
type GroupMessage struct {
	ID          int `xorm:"pk autoincr"`
	UserID      int `xorm:"notnull index"`
	UserUUID    string
	UserName    string
	UserAvatar  string
	UserCity    string
	UserAge     int
	GroupID     int    `xorm:"notnull index"`
	Msg         string `xorm:"notnull"`
	MsgType     int
	IsReply     bool //是否是回复消息
	ReplyUserID int  //如果是,被回复的用户id
	Context     []byte
	CreatedAt   time.Time `xorm:"created"`
	DeletedAt   time.Time `xorm:"deleted"`
	UpdatedAt   time.Time `xorm:"updated"`
}

func (GroupMessage) TableName() string {
	return "group_message"
}
