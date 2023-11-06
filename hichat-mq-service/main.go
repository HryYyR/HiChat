package main

import (
	adb "HiChat/hichat-mq-service/ADB"
	"HiChat/hichat-mq-service/config"
	"HiChat/hichat-mq-service/models"
	"encoding/json"
	"fmt"
	"strconv"
)

type GroupMsgfun func() error
type FriendMsgfun func() error

func main() {
	adb.InitMQ()
	adb.InitMySQL()

	// 获取接收消息的Delivery通道
	msgs, err := adb.MQc.Consume(
		adb.MQq.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		fmt.Println("接收失败!")
	}
	forever := make(chan bool)

	go func() {
		fmt.Println("开始消费!")
		for d := range msgs {

			// 群聊消息
			var msgstruct *models.Message
			err := json.Unmarshal(d.Body, &msgstruct)
			if err == nil && len(strconv.Itoa(msgstruct.MsgType)) < 4 {

				fmt.Println(msgstruct.MsgType)
				HandleMap := map[int]GroupMsgfun{
					config.MsgTypeDefault:      msgstruct.SaveGroupMsgToDb,
					config.MsgTypeImage:        msgstruct.SaveGroupMsgToDb,
					config.MsgTypeAudio:        msgstruct.SaveGroupMsgToDb,
					config.MsgTypeSyncMsg:      msgstruct.SyncGroupMsgToDb,
					config.MsgTypeClearSyncMsg: msgstruct.ClearGroupMsgNum,
				}
				err := HandleMap[msgstruct.MsgType]()
				if err != nil {
					fmt.Println(err.Error())
				}
				continue
				// if msgstruct.MsgType == config.MsgTypeDefault {
				// 	err = msgstruct.SaveGroupMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("保存群聊消息失败%s\n", err)
				// 	}
				// }
				// if msgstruct.MsgType == config.MsgTypeImage {
				// 	err = msgstruct.SaveGroupMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("保存群聊消息失败%s\n", err)
				// 	}
				// }
				// if msgstruct.MsgType == config.MsgTypeAudio {
				// 	err = msgstruct.SaveGroupMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("保存群聊消息失败%s\n", err)
				// 	}
				// }

				// if msgstruct.MsgType == config.MsgTypeSyncMsg {
				// 	err = msgstruct.SyncGroupMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("同步进数据库失败%s\n", err)
				// 	}
				// }
				// if msgstruct.MsgType == config.MsgTypeClearSyncMsg {
				// 	err = msgstruct.ClearGroupMsgNum()
				// 	if err != nil {
				// 		fmt.Printf("清除未读数据库失败%s\n", err)
				// 	}
				// }
				// continue
			}

			// 私聊消息
			var usermsgstruct *models.UserMessage
			err = json.Unmarshal(d.Body, &usermsgstruct)
			if err == nil {

				fmt.Println(msgstruct.MsgType)

				HandleMap := map[int]FriendMsgfun{
					config.MsgTypeFriendDefault:      usermsgstruct.SaveFriendMsgToDb,
					config.MsgTypeFriendImage:        usermsgstruct.SaveFriendMsgToDb,
					config.MsgTypeFriendAudio:        usermsgstruct.SaveFriendMsgToDb,
					config.MsgTypeSyncFriendMsg:      usermsgstruct.SyncFriendMsgToDb,
					config.MsgTypeClearSyncFriendMsg: usermsgstruct.ClearFriendMsgNum,
				}
				err := HandleMap[msgstruct.MsgType]()
				if err != nil {
					fmt.Println(err.Error())
				}
				continue

				// if msgstruct.MsgType == config.MsgTypeFriendDefault {
				// 	err = usermsgstruct.SaveFriendMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("保存好友消息失败%s\n", err)
				// 	}
				// }

				// if msgstruct.MsgType == config.MsgTypeFriendImage {
				// 	err = usermsgstruct.SaveFriendMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("保存好友消息失败%s\n", err)
				// 	}
				// }

				// if msgstruct.MsgType == config.MsgTypeFriendAudio {
				// 	err = usermsgstruct.SaveFriendMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("保存好友消息失败%s\n", err)
				// 	}
				// }

				// if msgstruct.MsgType == config.MsgTypeSyncFriendMsg {
				// 	err = usermsgstruct.SyncFriendMsgToDb()
				// 	if err != nil {
				// 		fmt.Printf("同步好友信息失败%s\n", err)
				// 	}
				// }
				// if msgstruct.MsgType == config.MsgTypeClearSyncFriendMsg {
				// 	err = usermsgstruct.ClearFriendMsgNum()
				// 	if err != nil {
				// 		fmt.Printf("清除好友同步信息失败%s\n", err)
				// 	}
				// }
				// continue
			}

		}
	}()

	<-forever

}
