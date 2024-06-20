package main

import (
	adb "HiChat/hichat-mq-service/ADB"
	"HiChat/hichat-mq-service/config"
	"HiChat/hichat-mq-service/models"
	"HiChat/hichat-mq-service/service_registry"
	"HiChat/hichat-mq-service/util"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log"
	"strconv"
)

type GroupMsgfun func() error
type FriendMsgfun func() error

// type RedisMsgfun func() error

func main() {
	adb.InitMQ()
	adb.InitMySQL()
	adb.InitRedis()

	go func() {
		models.BI.InsertMessages()
	}()

	//初始化雪花算法
	node, err := snowflake.NewNode(int64(config.ServerPort % 1023))
	if err != nil {
		log.Fatal("init snowflake error")
	}
	models.SnowFlakeNode = node

	// 获取接收消息的Delivery通道
	msgConsume, err := adb.MQc.Consume(
		adb.MQq.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		fmt.Println("创建消费者失败!")
		log.Fatal(err)
	}

	//服务注册
	addressIP := util.GetIP()
	dis := service_registry.DiscoveryConfig{
		ID:      util.GenerateUUID(),
		Name:    "hichat-mq-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: addressIP,
	}
	err = service_registry.RegisterService(dis)
	if err != nil {
		panic(err)
	}
	serverpost := fmt.Sprintf(":%s", strconv.Itoa(config.ServerPort))
	go service_registry.StartTcp(serverpost) //启动端口监听
	fmt.Println("服务注册在", serverpost)

	//go func() {
	fmt.Println("开始消费!")
	for d := range msgConsume {
		// 群聊消息
		var msgstruct *models.Message
		err := json.Unmarshal(d.Body, &msgstruct)
		if err == nil && msgstruct.MsgType < 500 {
			//fmt.Println(msgstruct.MsgType)
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
		}
		// 私聊消息
		var usermsgstruct *models.UserMessage
		err = json.Unmarshal(d.Body, &usermsgstruct)
		if err == nil && msgstruct.MsgType > 1000 && msgstruct.MsgType < 1500 {
			//fmt.Println(msgstruct.MsgType)
			HandleMap := map[int]FriendMsgfun{
				config.MsgTypeFriendDefault:      usermsgstruct.SaveFriendMsgToDb,
				config.MsgTypeFriendImage:        usermsgstruct.SaveFriendMsgToDb,
				config.MsgTypeFriendAudio:        usermsgstruct.SaveFriendMsgToDb,
				config.MsgTypeSyncFriendMsg:      usermsgstruct.SyncFriendMsgToDb,
				config.MsgTypeClearSyncFriendMsg: usermsgstruct.ClearFriendMsgNum,
			}
			if err := HandleMap[msgstruct.MsgType](); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("json Unmarshal error: ", err)
		}
	}
	//}()

}
