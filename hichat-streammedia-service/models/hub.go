package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

var ServiceCenter *Hub

type Hub struct {
	HubID     string                    //HUb的id
	Clients   map[int]UserClient        //用户列表  key:userid value:userclient
	Room      map[string]UserToUserRoom //房间列表 key:rromuuid value:usertouserroom
	Broadcast chan []byte               //广播列表
	Loginout  chan *UserClient          //退出登录的列表
	Mutex     sync.RWMutex              // 互斥锁     多个结构体实例可以共享同一个锁时用指针,此处只会创建一个,所以不用指针
}

func NewHub(HubID string) *Hub {
	return &Hub{
		HubID:     HubID,
		Clients:   make(map[int]UserClient),
		Room:      make(map[string]UserToUserRoom),
		Broadcast: make(chan []byte),
		Loginout:  make(chan *UserClient),
		Mutex:     sync.RWMutex{},
	}
}

type sdp struct {
	Type   string `json:"type"`
	Userid string `json:"userid"`
	Sdp    string `json:"sdp"`
}

func (h *Hub) Run() {
	defer func() {
		close(h.Broadcast)
		close(h.Loginout)
	}()
	for {
		select {
		// 消息广播到指定group
		case message := <-h.Broadcast:
			var data sdp
			err := json.Unmarshal(message, &data)
			if err != nil {
				fmt.Println("非法json", err)
			}

			uid, err := strconv.Atoi(data.Userid)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Printf("user: %v send: %v \n", data.Userid, data.Type)
			//fmt.Printf("%+v\n", data.Type)
			for _, room := range ServiceCenter.Room {
				if uid == room.StartUserID {
					ServiceCenter.Clients[room.ReceivedUserID].Send <- message
					break
				}
				if uid == room.ReceivedUserID {
					ServiceCenter.Clients[room.StartUserID].Send <- message
					break
				}
			}
		case c := <-h.Loginout:
			room, has := ServiceCenter.Room[c.BelongRoomUUID]

			fmt.Println("用户", c.UserID, "退出了")
			if has {
				msg := &sdp{
					Type:   "LoginOut",
					Userid: strconv.Itoa(c.UserID),
					Sdp:    "{}",
				}
				bytes, err := json.Marshal(msg)
				if err != nil {
					fmt.Println(err)
				}
				ServiceCenter.Clients[room.StartUserID].Send <- bytes
				ServiceCenter.Clients[room.ReceivedUserID].Send <- bytes
			}
			h.Mutex.Lock()
			delete(ServiceCenter.Room, c.BelongRoomUUID)
			delete(ServiceCenter.Clients, room.StartUserID)
			delete(ServiceCenter.Clients, room.ReceivedUserID)
			h.Mutex.Unlock()
		}
	}
}
