package models

import "time"

// Friend 不带msglist的friend
type Friend struct {
	Id        int32
	UserName  string
	NikeName  string
	Email     string
	Avatar    string
	City      string
	Age       string
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
}

// FriendResponse 带msglist的friend
type FriendResponse struct {
	Id            int32
	UserName      string
	NikeName      string
	Email         string
	Avatar        string
	City          string
	Age           string
	UnreadMessage int
	MessageList   []UserMessageItem
	CreatedAt     time.Time
	DeletedAt     time.Time
	UpdatedAt     time.Time
}

// UserMessageItem 用于返回的消息
type UserMessageItem struct {
	ID                int
	UserID            int
	UserName          string
	UserAvatar        string
	ReceiveUserID     int
	ReceiveUserName   string
	ReceiveUserAvatar string
	Msg               string
	MsgType           int
	IsReply           bool      //是否是回复消息
	ReplyUserID       int       //如果是,被回复的用户id
	CreatedAt         time.Time `xorm:"created"`
	DeletedAt         time.Time `xorm:"deleted"`
	UpdatedAt         time.Time `xorm:"updated"`
}

// UserMessage user_message models
type UserMessage struct {
	ID                int    `xorm:"pk autoincr index"`
	UUID              string `xorm:"notnull"`
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
