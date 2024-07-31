package models

import (
	"context"
	"encoding/json"
	"errors"
	adb "go-websocket-server/ADB"
	"go-websocket-server/config"
	"log"
)

//	type GroupMsgfun interface {
//		func(msgstruct *Message, msg_model []byte) error
//	}
type GroupMsgfun func(msgstruct *Message, msg []byte) error

var HandleGroupMsgMap = map[int]GroupMsgfun{
	config.MsgTypeDefault:      HandleDefaultGroupMsg,   //群聊文字	1
	config.MsgTypeImage:        HandleDefaultGroupMsg,   //群聊图片	2
	config.MsgTypeAudio:        HandleDefaultGroupMsg,   //群聊音频	3
	config.MsgTypeQuitGroup:    HandleDefaultGroupMsg,   //群聊退出	201
	config.MsgTypeJoinGroup:    HandleDefaultGroupMsg,   //群聊加入	202
	config.MsgTypeClearSyncMsg: HandleGroupClearSyncMsg, //群聊清除同步库	401
}

// HandleDefaultGroupMsg 1 默认消息
func HandleDefaultGroupMsg(msgstruct *Message, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 保存消息进数据库
	go func(msg []byte) {
		if err := adb.MqHub.PublishToNonImmediateTasksQueue(msg); err != nil {
			log.Printf("上传消息到队列失败!%s\n", err)
			cancelFunc()
		}
	}(msg)

	// 写入同步消息
	err := GroupWriteSyncMsg(msgstruct)
	if err != nil {
		return err
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

// GroupWriteSyncMsg 群聊写入同步消息
func GroupWriteSyncMsg(msgstruct *Message) error {
	msgstruct.MsgType = config.MsgTypeSyncMsg

	ctx, cancelFunc := context.WithCancel(context.Background())

	// 写入同步库
	// 4.将消息发布到声明的队列
	go func(syncmsg Message) {
		msgbyte, err := json.Marshal(syncmsg)
		if err != nil {
			log.Printf("同步消息转换byte失败!%s\n", err.Error())
			cancelFunc()
		}
		// 保存消息进数据库
		go func(msg []byte) {
			if err := adb.MqHub.PublishToNonImmediateTasksQueue(msg); err != nil {
				log.Printf("上传消息到队列失败!%s\n", err)
				cancelFunc()
			}
		}(msgbyte)
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

// HandleGroupClearSyncMsg 401 清除同步消息
func HandleGroupClearSyncMsg(msgstruct *Message, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 保存消息进数据库
	go func(msg []byte) {
		if err := adb.MqHub.PublishToNonImmediateTasksQueue(msg); err != nil {
			log.Printf("上传消息到队列失败!%s\n", err)
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

type FriendMsgfun func(msgstruct *UserMessage, msg []byte) error

var HandleFriendMsgMap = map[int]FriendMsgfun{
	config.MsgTypeFriendDefault:      HandleDefaultFriendMsg,   //好友文字	1001
	config.MsgTypeFriendImage:        HandleDefaultFriendMsg,   //好友图片	1002
	config.MsgTypeFriendAudio:        HandleDefaultFriendMsg,   //好友音频	1003
	config.MsgTypeClearSyncFriendMsg: HandleFriendClearSyncMsg, //好友清除同步库 1400
}

// HandleDefaultFriendMsg 1001  默认消息
func HandleDefaultFriendMsg(msgstruct *UserMessage, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 保存消息进数据库
	go func(msg []byte) {
		if err := adb.MqHub.PublishToNonImmediateTasksQueue(msg); err != nil {
			log.Printf("上传消息到队列失败!%s\n", err)
			cancelFunc()
		}
	}(msg)

	if err := FriendWriteSyncMsg(msgstruct); err != nil {
		log.Println("上传好友同步消息到队列失败,err")
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

// HandleFriendClearSyncMsg 1401 清除同步消息
func HandleFriendClearSyncMsg(msgstruct *UserMessage, msg []byte) error {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 保存消息进数据库
	go func(msg []byte) {
		if err := adb.MqHub.PublishToNonImmediateTasksQueue(msg); err != nil {
			log.Printf("上传消息到队列失败!%s\n", err)
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
	go func(usermsgstruct UserMessage) {
		msgbyte, err := json.Marshal(usermsgstruct)
		if err != nil {
			log.Printf("同步消息转换byte失败!%s\n", err.Error())
		}
		if err := adb.MqHub.PublishToNonImmediateTasksQueue(msgbyte); err != nil {
			log.Printf("上传消息到队列失败!%s\n", err)
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
