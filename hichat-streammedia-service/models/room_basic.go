package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserToUserRoom struct {
	RoomUUID string `json:"room_id"`
	RoomName string `json:"room_name"`
	RoomType int    `json:"room_type"`

	StartUserID   int    `json:"start_user_id"`
	StartUserName string `json:"start_user_name"`

	ReceivedUserID   int    `json:"received_user_id"`
	ReceivedUserName string `json:"received_user_name"`
	CreateTime       time.Time
}

// CheckUserLive 检查一个房间里的用户有没有连接超时
func (r *UserToUserRoom) CheckUserLive() {
	time.Sleep(time.Second * 5)
	sign := true

	if _, ok := ServiceCenter.Clients[r.StartUserID]; !ok {
		sign = false
	}
	if _, ok := ServiceCenter.Clients[r.ReceivedUserID]; !ok {
		sign = false
	}

	fmt.Println("双方均已连接")

	msg := sdp{
		Type:   "timeout",
		Userid: "",
		Sdp:    "",
	}
	msgbytes, _ := json.Marshal(msg)
	if !sign {
		if client, ok := ServiceCenter.Clients[r.StartUserID]; ok {
			ServiceCenter.Clients[r.StartUserID].Conn.Close()
			client.Send <- msgbytes
		}
		if client, ok := ServiceCenter.Clients[r.ReceivedUserID]; ok {
			ServiceCenter.Clients[r.ReceivedUserID].Conn.Close()
			client.Send <- msgbytes
		}
	}
}
