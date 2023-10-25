package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"log"
	"time"
)

type Message struct {
	ID          int    `xorm:"pk autoincr"`
	UserID      int    `xorm:"notnull"`
	UserName    string `xorm:"notnull"`
	UserAvatar  string
	UserCity    string
	UserAge     string
	GroupID     int `xorm:"notnull"`
	Msg         string
	MsgType     int  `xorm:"notnull default(1)"` //1文字 2图片 3音频 4文件
	IsReply     bool //是否是回复消息
	ReplyUserID int  //如果是,被回复的用户id
	Context     []byte
	UpdatedAt   time.Time `xorm:"updated"`
	CreatedAt   time.Time `xorm:"created"`
	DeletedAt   time.Time `xorm:"deleted"`
}

func (Message) TableName() string {
	return "group_message"
}

// 根据 groupid 获取用户列表
func (m *Message) AccordingToGroupidGetUserlist() ([]int, error) {
	var useridlist []int
	if err := adb.Ssql.Cols("user_id").Table("group_user_relative").Where("group_id=?", m.GroupID).Find(&useridlist); err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		return nil, err
	}
	return useridlist, nil
}

func (m *Message) SaveToDb() error {
	fmt.Printf("%+v\n", m)
	if _, err := adb.Ssql.Table("group_message").Insert(&m); err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		return err
	}
	return nil
}
