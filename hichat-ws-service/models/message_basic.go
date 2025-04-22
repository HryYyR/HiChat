package models

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"go-websocket-server/util"
	"log"
	"strconv"
	"strings"
	"time"
)

type BasicMessage struct {
	Type uint32          `json:"type"`
	Data json.RawMessage `json:"data"`
}

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
	MsgType    int  `xorm:"notnull default(1)"` //1文字 2 图片 3 音频 4 文件
	IsReply    bool //是否是回复消息
	ReplyMsgID int  //如果是,被回复的用户id
	Context    []byte
	UpdatedAt  time.Time `xorm:"updated"`
	CreatedAt  time.Time `xorm:"created"`
	DeletedAt  time.Time `xorm:"deleted"`
}

type GroupMsgTransmitFun func(groupmsg Message, msgbytes []byte) error

var TransmitGroupMsgMap = map[int]GroupMsgTransmitFun{
	config.MsgTypeDefault:               TransmitToAllFunc,         //1  群聊文字消息
	config.MsgTypeImage:                 TransmitToAllFunc,         //2  群聊图片消息
	config.MsgTypeAudio:                 TransmitToAllFunc,         //3  群聊音频消息
	config.MsgTypeQuitGroup:             TransmitToAllFunc,         //201  退出群聊
	config.MsgTypeJoinGroup:             TransmitToAllFunc,         //202  加入群聊
	config.MsgTypeApplyJoinGroup:        TransmitToUserIDFunc,      //203  申请加入群聊
	config.MsgTypeDissolveGroup:         TransmitDissolveGroupFunc, //204  解散群聊
	config.MsgTypeRefreshGroupAndNotice: TransmitToUserIDFunc,      //500  刷新群聊列表
	config.MsgTypeRefreshGroupNotice:    TransmitToUserIDFunc,      //502  刷新群聊通知列表
}

func TransmitDissolveGroupFunc(g Message, gbytes []byte) error {
	useridlist, err := util.BytesToInts(g.Context)
	if err != nil {
		log.Println(err)
		return err
	}
	//fmt.Printf("群id列表:%+v\n", useridlist)
	ServiceCenter.Mutex.Lock()
	// 给这个列表里的用户发送消息
	for _, userid := range useridlist {
		if clientlist, ok := ServiceCenter.Clients[userid]; ok {
			for i, client := range clientlist {
				if client.Status {
					//log.Println("给用户", userid, "群发信息", g.MsgType)
					ServiceCenter.Clients[userid][i].Send <- gbytes
				}
			}

		}
	}
	ServiceCenter.Mutex.Unlock()
	return nil
}

// TransmitToUserIDFunc 将此群聊消息转发给指定ID的用户
func TransmitToUserIDFunc(g Message, gbytes []byte) error {
	if clientlist, ok := ServiceCenter.Clients[g.UserID]; ok {
		for i, client := range clientlist {
			if client.Status {
				fmt.Println("成功转发给", g.UserID)
				ServiceCenter.Clients[g.UserID][i].Send <- gbytes
			}
		}

	}
	return nil
}

// TransmitToAllFunc 将此群聊消息转发给该群所有人
func TransmitToAllFunc(g Message, gbytes []byte) error {
	useridlist, err := g.AccordingToGroupidGetUserlist()
	//fmt.Printf("群id列表:%+v\n", useridlist)
	if err != nil {
		return err
	}
	// 给这个列表里的用户发送消息
	for _, userid := range useridlist {
		if clientlist, ok := ServiceCenter.Clients[userid]; ok {
			for i, client := range clientlist {
				if client.Status {
					//log.Println("给用户", userid, "群发信息", g.MsgType)
					ServiceCenter.Clients[userid][i].Send <- gbytes
				}
			}

		}
	}
	return nil
}

// Transmit 处理转发过来的消息
func (m Message) Transmit() error {
	msgbytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Println("处理群聊消息", m.MsgType)
	if fun, ok := TransmitGroupMsgMap[m.MsgType]; ok {
		err := fun(m, msgbytes)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// AccordingToGroupidGetUserlist 根据groupid获取用户列表
func (m *Message) AccordingToGroupidGetUserlist() ([]int, error) {
	var useridlist []int
	strgroupid := strconv.Itoa(m.GroupID)
	result := adb.Rediss.HGet("GroupToUserMap", strgroupid).Val()
	if len(result) == 0 {
		if err := adb.SqlStruct.Conn.Cols("user_id").Table("group_user_relative").Where("group_id=?", m.GroupID).Find(&useridlist); err != nil {
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
