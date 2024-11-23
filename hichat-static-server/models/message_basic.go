package models

import "time"

// Message 用户传输的群聊消息结构体
type Message struct {
	ID         int    `xorm:"pk autoincr"`
	UserID     int    `xorm:"notnull"`
	UserName   string `xorm:"notnull"`
	UserAvatar string
	UserCity   string
	UserAge    int
	GroupID    int `xorm:"notnull"`
	Msg        string
	MsgType    int  `xorm:"notnull default(1)"` //1文字 2图片 3音频 4文件
	IsReply    bool //是否是回复消息
	ReplyMsgID int  //如果是,被回复的用户id
	Context    []byte
	UpdatedAt  time.Time `xorm:"updated"`
	CreatedAt  time.Time `xorm:"created"`
	DeletedAt  time.Time `xorm:"deleted"`
}
