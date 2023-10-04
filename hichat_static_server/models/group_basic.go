package models

import "time"

type GroupDetail struct {
	GroupInfo   Group
	MessageList []GroupMessage
}

type Group struct {
	ID          int    `xorm:"pk autoincr notnull index"`
	UUID        string `xorm:"notnull unique"`
	CreaterID   int    `xorm:"notnull"`
	CreaterName string `xorm:"notnull"`
	GroupName   string `xorm:"notnull"`
	Avatar      string
	Grade       int `xorm:"default(1)"`
	MemberCount int
	CreatedAt   time.Time `xorm:"created"`
	DeletedAt   time.Time `xorm:"deleted"`
	UpdatedAt   time.Time `xorm:"updated"`
}

// 群聊消息
type GroupMessage struct {
	ID          int `xorm:"pk autoincr"`
	UserID      int `xorm:"notnull"`
	UserUUID    string
	UserName    string
	GroupID     int    `xorm:"notnull"`
	Msg         string `xorm:"notnull"`
	MsgType     int
	IsReply     bool //是否是回复消息
	ReplyUserID int  //如果是,被回复的用户id
	Context     []byte
	CreatedAt   time.Time `xorm:"created"`
	DeletedAt   time.Time `xorm:"deleted"`
	UpdatedAt   time.Time `xorm:"updated"`
}

type GroupUserRelative struct {
	ID        int `xorm:"pk autoincr notnull index"`
	UserID    int
	GroupID   int
	GroupUUID string
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
	UpdatedAt time.Time `xorm:"updated"`
}
