package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/util"
	"log"
	"strconv"
	"strings"
	"time"
)

// Message 用户传输的消息结构体
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

func (m *Message) TableName() string {
	return "group_message"
}

// AccordingToGroupidGetUserlist 根据groupid获取用户列表
func (m *Message) AccordingToGroupidGetUserlist() ([]int, error) {
	var useridlist []int
	strgroupid := strconv.Itoa(m.GroupID)
	result := adb.Rediss.HGet("GroupToUserMap", strgroupid).Val()
	if len(result) == 0 {
		if err := adb.SqlStruct.Conn.Cols("user_id").Table("group_user_relative").Where("group_id=?", m.GroupID).Find(&useridlist); err != nil {
			fmt.Println(err.Error())
			log.Println(err.Error())
			return nil, err
		}
		//存 redis
		stringuseridlist := util.IntArrToStrArr(useridlist)
		joinstr := strings.Join(stringuseridlist, ",")
		adb.Rediss.HSet("GroupToUserMap", strgroupid, joinstr)

	} else {
		strarr := strings.Split(result, ",")
		useridlist = util.StrArrToIntArr(strarr)
	}
	return useridlist, nil

}

//func (m *Message) SaveToDb() error {
//	fmt.Printf("%+v\n", m)
//	if _, err := adb.SqlStruct.Conn.Table("group_message").Insert(&m); err != nil {
//		fmt.Println(err.Error())
//		log.Println(err.Error())
//		return err
//	}
//	return nil
//}
