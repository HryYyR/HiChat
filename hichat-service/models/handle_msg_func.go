package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"

	"github.com/streadway/amqp"
)

//	type GroupMsgfun interface {
//		func(msgstruct *Message, msg []byte) error
//	}
type GroupMsgfun func(msgstruct *Message, msg []byte) error

var HandleGroupMsgMap = map[int]GroupMsgfun{
	1:   HandleDefaultGroupMsg,   //群聊文字
	2:   HandleDefaultGroupMsg,   //群聊图片
	3:   HandleDefaultGroupMsg,   //群聊音频
	201: HandleDefaultGroupMsg,   //群聊退出
	202: HandleDefaultGroupMsg,   //群聊加入
	401: HandleGroupClearSyncMsg, //群聊清除同步库
}

// HandleDefaultGroupMsg 1 默认消息
func HandleDefaultGroupMsg(msgstruct *Message, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

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
			cancelFunc()
		}
	}(msg)

	//根据 groupid 获取用户id列表
	useridlist, err := msgstruct.AccordingToGroupidGetUserlist()
	if err != nil {
		fmt.Printf("获取用户id列表失败!%s\n", err)
		return err
	}

	// 写入同步消息
	err = GroupWriteSyncMsg(msgstruct)
	if err != nil {
		return err
	}

	// 给这个列表里的用户发送消息
	for _, userid := range useridlist {
		fmt.Println("给用户发信息", userid)
		if ServiceCenter.Clients[userid].Status {
			ServiceCenter.Clients[userid].Send <- msg
		}
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return err
		}
	default:
		return nil
	}
	return nil
}

// HandleGroupClearSyncMsg 401 清除同步消息
func HandleGroupClearSyncMsg(msgstruct *Message, msg []byte) error {

	ctx, cancelFunc := context.WithCancel(context.Background())
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
			cancelFunc()
		}
	}(msg)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return err
		}
	default:
		return nil
	}
	return nil
}

// GroupWriteSyncMsg 群聊写入同步消息
func GroupWriteSyncMsg(msgstruct *Message) error {
	msgstruct.MsgType = config.MsgTypeSyncMsg

	ctx, cancelFunc := context.WithCancel(context.Background())

	// 写入同步库
	// 4.将消息发布到声明的队列
	go func(syncmsg Message) {
		msgbyte, err := json.Marshal(syncmsg)
		if err != nil {
			fmt.Printf("同步消息转换byte失败!%s\n", err.Error())
			cancelFunc()
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
			cancelFunc()
		}
	}(*msgstruct)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return err
		}
	default:
		return nil
	}
	return nil
}

type FriendMsgfun func(msgstruct *UserMessage, msg []byte) error

var HandleFriendMsgMap = map[int]FriendMsgfun{
	1001: HandleDefaultFriendMsg,   //好友文字
	1002: HandleDefaultFriendMsg,   //好友图片
	1003: HandleDefaultFriendMsg,   //好友音频
	1401: HandleFriendClearSyncMsg, //好友清除同步库
}

// HandleDefaultFriendMsg 1001  默认消息
func HandleDefaultFriendMsg(msgstruct *UserMessage, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	//bytes, _ := json.Marshal(msgstruct)
	// 保存消息进数据库
	go func() {
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
			cancelFunc()
		}

	}()

	FriendWriteSyncMsg(msgstruct)

	//fmt.Println("未加密的消息", string(bytes))

	ServiceCenter.Clients[msgstruct.UserID].Send <- msg
	ServiceCenter.Clients[msgstruct.ReceiveUserID].Send <- msg

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return err
		}
	default:
		return nil
	}
	return nil
}

// HandleFriendClearSyncMsg 1401 清除同步消息
func HandleFriendClearSyncMsg(msgstruct *UserMessage, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

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
			cancelFunc()
		}
	}(msg)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return err
		}
	default:
		return nil
	}

	return nil
}

// FriendWriteSyncMsg 好友写入同步消息
func FriendWriteSyncMsg(usermsgstruct *UserMessage) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

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
			cancelFunc()
		}
	}(*usermsgstruct)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.Canceled) {
			return err
		}
	default:
		return nil
	}

	return nil
}
