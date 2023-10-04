package models

import (
	"encoding/json"
	"fmt"
)

var ServiceCenter *Hub

type Hub struct {
	HubID     string             //HUb的id
	Clients   map[int]UserClient //用户列表  key:userid value:userclient
	Broadcast chan []byte        //广播列表
	Loginout  chan *UserClient   //退出登录的列表
}

func NewHub(HubID string) *Hub {
	return &Hub{
		HubID:     HubID,
		Clients:   make(map[int]UserClient),
		Broadcast: make(chan []byte),
		Loginout:  make(chan *UserClient),
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
		case UserClient := <-h.Loginout:
			client := ServiceCenter.Clients[UserClient.UserID]
			client.Status = false
			ServiceCenter.Clients[UserClient.UserID] = client

		// 消息广播到指定group
		case message := <-h.Broadcast:
			var msgstruct *Message
			if err := json.Unmarshal(message, &msgstruct); err != nil {
				fmt.Println(err)
			}
			useridlist, err := msgstruct.AccordingToGroupidGetUserlist() //根据 groupid 获取用户id列表
			if err != nil {
				fmt.Println(err.Error())
			}

			// 给这个列表里的用户发送消息
			for clientid, UserClient := range ServiceCenter.Clients {
				for _, userid := range useridlist {
					if clientid == userid {
						if UserClient.Status {
							ServiceCenter.Clients[clientid].Send <- message
						} else if msgstruct.MsgType == 1 { //除了默认消息,其他消息不缓存
							v, ok := ServiceCenter.Clients[clientid].CachingMessages[msgstruct.GroupID]
							if !ok {
								v = 0
							}
							v++
							ServiceCenter.Clients[clientid].CachingMessages[msgstruct.GroupID] = v
						}
					}
				}

			}

		}
	}
}
