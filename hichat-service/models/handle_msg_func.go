package models

import (
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"

	"github.com/streadway/amqp"
)

type GroupMsgfun func(msgstruct *Message, msg []byte)

var HandleGroupMsgMap = map[int]GroupMsgfun{
	1:   HandleDefaultGroupMsg,   //群聊文字
	2:   HandleDefaultGroupMsg,   //群聊图片
	3:   HandleDefaultGroupMsg,   //群聊音频
	201: HandleDefaultGroupMsg,   //群聊退出
	202: HandleDefaultGroupMsg,   //群聊加入
	401: HandleGroupClearSyncMsg, //群聊清除同步库
}

// 1 默认消息
func HandleDefaultGroupMsg(msgstruct *Message, msg []byte) {
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
	GroupWriteSyncMsg(msgstruct)

	// 给这个列表里的用户发送消息
	for _, userid := range useridlist {
		if ServiceCenter.Clients[userid].Status {
			ServiceCenter.Clients[userid].Send <- msg
		}
	}
}

// 401 清除同步消息
func HandleGroupClearSyncMsg(msgstruct *Message, msg []byte) {
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

// 群聊写入同步消息
func GroupWriteSyncMsg(msgstruct *Message) {
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

type FriendMsgfun func(msgstruct *UserMessage, msg []byte)

var HandleFriendMsgMap = map[int]FriendMsgfun{
	1001: HandleDefaultFriendMsg,   //好友文字
	1002: HandleDefaultFriendMsg,   //好友图片
	1003: HandleDefaultFriendMsg,   //好友音频
	1401: HandleFriendClearSyncMsg, //好友清除同步库
}

// 1001  默认消息
func HandleDefaultFriendMsg(msgstruct *UserMessage, msg []byte) {
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

	FriendWriteSyncMsg(msgstruct)

	ServiceCenter.Clients[msgstruct.UserID].Send <- msg
	ServiceCenter.Clients[msgstruct.ReceiveUserID].Send <- msg
}

// 1401 清除同步消息
func HandleFriendClearSyncMsg(msgstruct *UserMessage, msg []byte) {
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

// 好友写入同步消息
func FriendWriteSyncMsg(usermsgstruct *UserMessage) {
	usermsgstruct.MsgType = config.MsgTypeSyncFriendMsg
	// 写入同步库
	// 4.将消息发布到声明的队列
	go func(syncmsg UserMessage) {
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
	}(*usermsgstruct)
}
