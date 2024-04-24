package models

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"go-websocket-server/config"
	"go-websocket-server/util"
	"log"
	"strconv"
	"sync"
)

var ServiceCenter *Hub

type Hub struct {
	HubID      string             //HUb的id
	Clients    map[int]UserClient //用户列表  key:userid value:userclient
	Broadcast  chan []byte        //广播列表
	Loginout   chan *UserClient   //退出登录的列表
	Mutex      *sync.RWMutex      // 互斥锁     用指针时多个结构体实例共享同一个锁,否则每个实例有属于自己的锁
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (h *Hub) setPrivateKey(pkey *rsa.PrivateKey) {
	h.privateKey = pkey
}
func (h *Hub) GetPrivateKey() *rsa.PrivateKey {
	return h.privateKey
}
func (h *Hub) setPublicKey(pkey *rsa.PublicKey) {
	h.publicKey = pkey
}
func (h *Hub) GetPublicKey() *rsa.PublicKey {
	return h.publicKey
}

func NewHub(HubID string) *Hub {
	publicKey, privateKey := util.GenerateRsaKey()
	return &Hub{
		HubID:      HubID,
		Clients:    make(map[int]UserClient),
		Broadcast:  make(chan []byte),
		Loginout:   make(chan *UserClient),
		Mutex:      &sync.RWMutex{},
		publicKey:  publicKey,
		privateKey: privateKey,
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
			client.HoldEncryptedKey = false
			client.EncryptedKey = []byte{}
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
				//err := HandleGroupMsgMap[msgstruct.MsgType](msgstruct, message)
				if msgfun, ok := HandleGroupMsgMap[msgstruct.MsgType]; ok {
					err := msgfun(msgstruct, message)
					if err != nil {
						log.Println("HandleGroupMsgMap error: ", err)
						fmt.Println("HandleGroupMsgMap error: ", err)
					}
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
			//fmt.Printf("%+v\n", usermsgstruct)
			if err == nil {
				fmt.Println("friendmsg:", msgstruct.MsgType)
				if msgfun, ok := HandleFriendMsgMap[msgstruct.MsgType]; ok {
					err := msgfun(usermsgstruct, message)
					if err != nil {
						log.Println("HandleFriendMsgMap", err)
						fmt.Println("HandleFriendMsgMap", err)
					}
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
				//log.Println(err)
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
