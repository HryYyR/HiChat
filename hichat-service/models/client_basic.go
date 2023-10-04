// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// 用户客户端
type UserClient struct {
	ClientID        string
	UserID          int
	UserUUID        string
	UserName        string
	Conn            *websocket.Conn
	Status          bool
	Send            chan []byte
	Groups          map[int]Group //群聊列表  key:group_id  value:group
	CachingMessages map[int]int   // key:groupid  value:未读数量
}

// 读取用户发送的信息
func (c *UserClient) ReadPump() {
	defer func() {
		ServiceCenter.Loginout <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("有用户退出了")
				ServiceCenter.Loginout <- c
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("IsUnexpectedCloseError: %v", err)
			}
			break
		}

		// 保存消息进数据库
		var msgstruct *Message
		if err := json.Unmarshal(message, &msgstruct); err != nil {
			fmt.Println(err)
		}
		err = msgstruct.SaveToDb()
		if err != nil {
			fmt.Println("消息保存失败,取消发送", err)
			return
		}
		// end

		// message = bytes.TrimSpace(bytes.Replace(message, Newline, Space, -1))
		ServiceCenter.Broadcast <- message
	}
}

// 给用户发送信息
func (c *UserClient) WritePump() {
	ticker := time.NewTicker(PingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("NextWriterError", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				fmt.Println("WriterError", err)
			}

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(Newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				fmt.Println("WriterCloseError", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
