package models

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/util"
	"log"
	"strconv"
	"sync"
)

var ServiceCenter *Hub

type MessageTransmitter interface {
	Transmit() error
}

type Hub struct {
	HubID      string                  //HUb的id
	Clients    map[int][]*UserClient   //用户列表  key:userid value:userclient
	Broadcast  chan []byte             //广播列表
	Loginout   chan *UserClient        //退出登录的列表
	Transmit   chan MessageTransmitter //转发列表
	Mutex      *sync.RWMutex           // 互斥锁     用指针时多个结构体实例共享同一个锁,否则每个实例有属于自己的锁
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
		Clients:    make(map[int][]*UserClient),
		Broadcast:  make(chan []byte),
		Transmit:   make(chan MessageTransmitter),
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
			if UC.Status == false {
				continue
			}
			for _, client := range ServiceCenter.Clients[UC.UserID] {
				fmt.Println(client.Status, client.UserID, client.Device)
				if client.Device == UC.Device && client.Status {
					client.Mutex.Lock()
					// 获取当前设备标志，处理可能的错误
					deviceSign, _ := adb.Rediss.HGet("UserClient", strconv.Itoa(UC.UserID)).Int()

					// 更新设备标志
					newDeviceSign := deviceSign &^ int(UC.Device)
					adb.Rediss.HSet("UserClient", strconv.Itoa(UC.UserID), newDeviceSign)
					client.Mutex.Unlock()

				}
			}

		// 消息广播给指定用户
		case message := <-h.Broadcast:
			var basicMsg *BasicMessage
			if err := json.Unmarshal(message, &basicMsg); err != nil {
				fmt.Printf("msgtype:%d ,非法msg %s\n", basicMsg.Type, err)
				continue
			}
			switch basicMsg.Type {
			case 1:
				// 群聊消息
				var msg Message
				if err := json.Unmarshal(basicMsg.Data, &msg); err != nil {
					log.Printf("解析 GroupMessage 失败: %s\n", err)
					continue
				}
				fmt.Println(msg)
				if len(strconv.Itoa(msg.MsgType)) < 4 {
					if msgfun, ok := HandleGroupMsgMap[msg.MsgType]; ok {
						go func(msgfunc GroupMsgfun, types int) {
							err := msgfunc(&msg, basicMsg.Data)
							if err != nil {
								log.Println("HandleGroupMsgMap error: ", err)
							} else {
								if types < 399 {
									TransmitMsg(message, types) //群聊消息保存成功后,转发消息
								}
							}

						}(msgfun, msg.MsgType)
					}
					continue
				}
			case 2:
				var msg UserMessage
				if err := json.Unmarshal(basicMsg.Data, &msg); err != nil {
					log.Printf("解析 UserMessage 失败: %s\n", err)
					continue
				}
				fmt.Println(msg)
				if msgfun, ok := HandleFriendMsgMap[msg.MsgType]; ok {
					go func(msgfunc FriendMsgfun, types int) {
						err := msgfunc(&msg, basicMsg.Data)
						if err != nil {
							log.Println("HandleFriendMsgMap", err)
						} else {
							if types < 1399 {
								TransmitMsg(message, types) //用户保存成功后,转发消息
							}

						}
					}(msgfun, msg.MsgType)
				} else {
					log.Printf("处理方法不存在 type:%d\n", msg.MsgType)
				}
			default:
				log.Printf("非法的消息类型 type:%d\n", basicMsg.Type)
			}
			//
			//var groupmsgstruct *Message
			//err := json.Unmarshal(message, &groupmsgstruct)
			//if err == nil && len(strconv.Itoa(groupmsgstruct.MsgType)) < 4 {
			//	if msgfun, ok := HandleGroupMsgMap[groupmsgstruct.MsgType]; ok {
			//		go func(msgfunc GroupMsgfun, types int) {
			//			err := msgfunc(groupmsgstruct, message)
			//			if err != nil {
			//				log.Println("HandleGroupMsgMap error: ", err)
			//			} else {
			//				if types < 399 {
			//					TransmitMsg(message, types) //群聊消息保存成功后,转发消息
			//				}
			//			}
			//
			//		}(msgfun, groupmsgstruct.MsgType)
			//	}
			//	continue
			//}

			// 好友消息
			//var usermsgstruct *UserMessage
			//err = json.Unmarshal(message, &usermsgstruct)
			////log.Printf("%+v\n", usermsgstruct)
			//if err == nil {
			//	//log.Println("friendmsg:", usermsgstruct.MsgType)
			//	if msgfun, ok := HandleFriendMsgMap[usermsgstruct.MsgType]; ok {
			//		go func(msgfunc FriendMsgfun, types int) {
			//			err := msgfunc(usermsgstruct, message)
			//			if err != nil {
			//				log.Println("HandleFriendMsgMap", err)
			//			} else {
			//				if types < 1399 {
			//					TransmitMsg(message, types) //用户保存成功后,转发消息
			//				}
			//
			//			}
			//		}(msgfun, usermsgstruct.MsgType)
			//	} else {
			//		log.Println("处理方法不存在")
			//	}
			//}
		case msg := <-h.Transmit:
			go func(m MessageTransmitter) {
				err := m.Transmit()
				if err != nil {
					log.Println("HandleTransmit error:", err)
				}
			}(msg)

		}
	}
}
