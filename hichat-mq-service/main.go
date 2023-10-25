package main

import (
	adb "HiChat/hichat-mq-service/ADB"
	"HiChat/hichat-mq-service/config"
	"HiChat/hichat-mq-service/models"
	"encoding/json"
	"fmt"
)

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
			var msgstruct *models.Message
			if err := json.Unmarshal(d.Body, &msgstruct); err != nil {
				fmt.Println(err.Error())
			}
			fmt.Printf("message type:%v\n", msgstruct.MsgType)

			if msgstruct.MsgType == config.MsgTypeDefault {
				err = msgstruct.SaveToDb()
				if err != nil {
					fmt.Printf("保存进数据库失败%s\n", err)
				}
			}
			if msgstruct.MsgType == config.MsgTypeImage {
				err = msgstruct.SaveToDb()
				if err != nil {
					fmt.Printf("保存进数据库失败%s\n", err)
				}
			}
			if msgstruct.MsgType == config.MsgTypeAudio {
				err = msgstruct.SaveToDb()
				if err != nil {
					fmt.Printf("保存进数据库失败%s\n", err)
				}
			}

			if msgstruct.MsgType == config.MsgTypeSyncMsg {
				err = msgstruct.SyncToDb()
				if err != nil {
					fmt.Printf("同步进数据库失败%s\n", err)
				}
			}
			if msgstruct.MsgType == config.MsgTypeClearSyncMsg {
				err = msgstruct.ClearMsgNum()
				if err != nil {
					fmt.Printf("清除未读数据库失败%s\n", err)
				}
			}

		}
	}()

	<-forever

}
