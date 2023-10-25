package models

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"

	"github.com/streadway/amqp"
)

type Msgfun func(msgstruct *Message, msg []byte)

var HandleMsgMap = map[int]Msgfun{
	1:   HandleDefaultMsg,   //群聊文字
	2:   HandleDefaultMsg,   //群聊图片
	3:   HandleDefaultMsg,   //群聊音频
	201: HandleDefaultMsg,   //群聊退出
	202: HandleDefaultMsg,   //群聊加入
	401: HandleClearSyncMsg, //群聊清除同步库
}

// 1
func HandleDefaultMsg(msgstruct *Message, msg []byte) {
	// 保存消息进数据库
	go func(msg []byte) {
		err := adb.MQc.Publish(
			"",           // exchange
			adb.MQq.Name, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msg,
			})
		if err != nil {
			fmt.Printf("上传消息到队列失败!%s\n", err)
		}
	}(msg)

	//根据 groupid 获取用户id列表
	useridlist, err := msgstruct.AccordingToGroupidGetUserlist()
	if err != nil {
		fmt.Printf("获取用户id列表失败!%s\n", err)
	}

	// 写入同步消息
	WriteSyncMsg(msgstruct)

	// 给这个列表里的用户发送消息
	for _, userid := range useridlist {
		if ServiceCenter.Clients[userid].Status {
			ServiceCenter.Clients[userid].Send <- msg
		}
	}
	// for clientid, UserClient := range ServiceCenter.Clients {
	// 	for _, userid := range useridlist {
	// 		if clientid == userid {
	// 			if UserClient.Status {
	// 				ServiceCenter.Clients[clientid].Send <- msg
	// 			} else {
	// 				v, ok := ServiceCenter.Clients[clientid].CachingMessages[msgstruct.GroupID]
	// 				if !ok {
	// 					v = 0
	// 				}
	// 				v++
	// 				// ServiceCenter.Clients[clientid].Mutex.Lock()
	// 				ServiceCenter.Clients[clientid].CachingMessages[msgstruct.GroupID] = v
	// 				// ServiceCenter.Clients[clientid].Mutex.Unlock()
	// 			}
	// 			// 除了发布者以外,写入同步消息
	// 			if userid != msgstruct.UserID {
	// 				WriteSyncMsg(msgstruct, &UserClient)
	// 			}
	// 		}
	// 	}
	// }

}

// 401 清除同步消息
func HandleClearSyncMsg(msgstruct *Message, msg []byte) {
	go func(msg []byte) {
		err := adb.MQc.Publish(
			"",           // exchange
			adb.MQq.Name, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msg,
			})
		if err != nil {
			fmt.Printf("上传消息到队列失败!%s\n", err)
		}
	}(msg)
}

// 写入同步消息
func WriteSyncMsg(msgstruct *Message) {
	msgstruct.MsgType = config.MsgTypeSyncMsg
	// 写入同步库
	// 4.将消息发布到声明的队列
	go func(syncmsg Message) {
		msgbyte, err := json.Marshal(syncmsg)
		if err != nil {
			fmt.Printf("同步消息转换byte失败!%s\n", err.Error())
		}
		err = adb.MQc.Publish(
			"",           // exchange
			adb.MQq.Name, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msgbyte,
			})
		if err != nil {
			fmt.Printf("上传消息到队列失败!%s\n", err)
		}
	}(*msgstruct)
}
