package models

import (
	adb "HiChat/hichat-mq-service/ADB"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log"
	"strconv"
	"sync"
	"time"
)

var M sync.Mutex
var SnowFlakeNode *snowflake.Node

type Message struct {
	ID         int    `xorm:"pk autoincr"`
	UserID     int    `xorm:"notnull"`
	UserName   string `xorm:"notnull"`
	UserAvatar string
	UserCity   string
	UserAge    int
	GroupID    int `xorm:"notnull"`
	Msg        string
	MsgType    int  `xorm:"notnull default(1)"` //1 文字 2 音频 3 视频 4 文件
	IsReply    bool //是否是回复消息
	ReplyMsgID int  //如果是,被回复的用户id
	Context    []byte
	UpdatedAt  time.Time `xorm:"updated"`
	CreatedAt  time.Time `xorm:"created"`
	DeletedAt  time.Time `xorm:"deleted"`
}

func (*Message) TableName() string {
	return "group_message"
}

// AccordingToGroupidGetUserlist 根据 groupid 获取用户列表
func (m *Message) AccordingToGroupidGetUserlist() ([]int, error) {
	var useridlist []int
	if err := adb.Ssql.Cols("user_id").Table("group_user_relative").Where("group_id=?", m.GroupID).Find(&useridlist); err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		return nil, err
	}
	return useridlist, nil
}

func (m *Message) SaveGroupMsgToDb() error {
	// fmt.Printf("%+v\n", m)
	//if _, err := adb.Ssql.Table("group_message").Insert(&m); err != nil {
	//	fmt.Println(err.Error())
	//	return err
	//}
	BI.AddGroupMsg(m)

	//m.ID = int(SnowFlakeNode.Generate())
	jsondata, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err.Error())
		//todo
	}
	err = adb.Rediss.RPush(fmt.Sprintf("gm%s", strconv.Itoa(m.GroupID)), string(jsondata)).Err()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (m *Message) SyncGroupMsgToDb() error {
	// fmt.Printf("%+v\n", m)
	useridlist, err := m.AccordingToGroupidGetUserlist() //获取群里所有用户id
	if err != nil {
		fmt.Printf("查询用户id列表失败%s\n", err)
		return err
	}
	for _, userid := range useridlist {
		// 除了发送者本身,其余都同步
		if userid == m.UserID {
			continue
		}
		// 判断 未读表里是否存在 该用户和群关系
		var groupunreadmessage GroupUnreadMessage
		isexit, err := adb.Ssql.Table("group_unread_message").Where("group_id=? and user_id=?", m.GroupID, userid).Get(&groupunreadmessage)
		if err != nil {
			fmt.Printf("查询失败%s\n", err)
			return err
		}
		// 存在,就+1
		if isexit {
			M.Lock()
			_, err := adb.Ssql.Table("group_unread_message").ID(groupunreadmessage.ID).Cols("unread_number").Update(GroupUnreadMessage{
				UnreadNumber: groupunreadmessage.UnreadNumber + 1,
			})
			M.Unlock()
			if err != nil {
				fmt.Printf("更新失败%s\n", err)
				return err
			}
		} else {
			// 不存在,先查他的用户名称再插入
			var userdata Users
			exit, err := adb.Ssql.Cols("id,user_name").Table("users").Where("id=?", userid).Get(&userdata)
			if err != nil {
				fmt.Printf("查询用户数据失败%s\n", err)
				return err
			}
			if exit {
				newmsg := GroupUnreadMessage{
					UserName:     userdata.UserName,
					UserID:       userdata.ID,
					GroupID:      m.GroupID,
					UnreadNumber: 1,
				}
				_, err = adb.Ssql.Table("group_unread_message").Insert(&newmsg)
				if err != nil {
					fmt.Printf("插入失败%s\n", err)
					return err
				}
			} else {
				fmt.Printf("查询用户数据不存在%s\n", err)
				return err
			}
		}
	}
	return nil
}

func (m *Message) ClearGroupMsgNum() error {
	var willcleardata GroupUnreadMessage
	is, err := adb.Ssql.Table("group_unread_message").Where("user_id = ? and group_id=?", m.UserID, m.GroupID).Get(&willcleardata)
	if err != nil {
		return err
	}
	if is {
		_, err := adb.Ssql.Table("group_unread_message").ID(willcleardata.ID).Cols("unread_number").Update(GroupUnreadMessage{
			UnreadNumber: 0,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
