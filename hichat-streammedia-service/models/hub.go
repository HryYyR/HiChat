package models

import (
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
		// 消息广播到指定group
		case message := <-h.Broadcast:
			for _, v := range ServiceCenter.Clients {
				v.Send <- message
			}
		}
	}
}
