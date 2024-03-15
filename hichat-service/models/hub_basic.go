package models

import (
	"encoding/json"
	"fmt"
	"go-websocket-server/config"
	"log"
	"strconv"
	"sync"
)

var ServiceCenter *Hub

type Hub struct {
	HubID     string             //HUb的id
	Clients   map[int]UserClient //用户列表  key:userid value:userclient
	Broadcast chan []byte        //广播列表
	Loginout  chan *UserClient   //退出登录的列表
	Mutex     sync.RWMutex       // 互斥锁     多个结构体实例可以共享同一个锁时用指针,此处只会创建一个,所以不用指针
}

func NewHub(HubID string) *Hub {
	return &Hub{
		HubID:     HubID,
		Clients:   make(map[int]UserClient),
		Broadcast: make(chan []byte),
		Loginout:  make(chan *UserClient),
		Mutex:     sync.RWMutex{},
	}
}

func (h *Hub) Run() {
	defer func() {
		close(h.Broadcast)
		close(h.Loginout)
	}()
	for {
		select {
		// 退出登录
		case UC := <-h.Loginout:
			client := ServiceCenter.Clients[UC.UserID]
			client.Status = false
			client.Conn = nil
			ServiceCenter.Clients[UC.UserID].Mutex.Lock()
			ServiceCenter.Clients[UC.UserID] = client
			ServiceCenter.Clients[UC.UserID].Mutex.Unlock()

		// 消息广播给指定用户
		case message := <-h.Broadcast:

			// 群聊消息
			var msgstruct *Message
			err := json.Unmarshal(message, &msgstruct)
			if err == nil && len(strconv.Itoa(msgstruct.MsgType)) < 4 {
				fmt.Println("groupmsg:", msgstruct.MsgType)
				err := HandleGroupMsgMap[msgstruct.MsgType](msgstruct, message)
				if err != nil {
					log.Println(err)
				}
				//todo
				//if msgstruct.MsgType < 100 {
				//	if err != nil {
				//		sendAckMsg(2, msgstruct.UserID, 0)
				//	} else {
				//		sendAckMsg(2, msgstruct.UserID, 1)
				//	}
				//}
				continue
			}

			// 好友消息
			var usermsgstruct *UserMessage
			err = json.Unmarshal(message, &usermsgstruct)
			if err == nil {
				fmt.Println("friendmsg:", msgstruct.MsgType)
				err := HandleFriendMsgMap[msgstruct.MsgType](usermsgstruct, message)
				if err != nil {
					log.Println(err)
				}
				//todo
				//if usermsgstruct.MsgType < 1100 {
				//	if err != nil {
				//		sendAckMsg(1, usermsgstruct.UserID, 0)
				//	} else {
				//		sendAckMsg(1, usermsgstruct.UserID, 1)
				//	}
				//
				//}

			} else {
				log.Println(err)
				fmt.Println("解析消息体失败:error", err)
			}
		}
	}
}

func sendAckMsg(msgsort, uid, status int) {
	ackmsg := &AckMessage{
		MsgType:   config.MsgTypeAckMsg,
		AckStatus: status,
		UserId:    uid,
		MsgSort:   msgsort,
	}
	ackbytes, _ := json.Marshal(ackmsg)
	ServiceCenter.Clients[uid].Send <- ackbytes
}
