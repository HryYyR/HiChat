package models

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// 用户客户端
type UserClient struct {
	ClientID       string
	UserID         int
	BelongRoomUUID string
	UserName       string
	Conn           *websocket.Conn
	Send           chan []byte
}

// 读取用户发送的信息
func (c *UserClient) ReadPump() {

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("用户", c.UserName, "退出了")
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("IsUnexpectedCloseError: %v \n", err)
			}
			ServiceCenter.Loginout <- c
			break
		}
		// fmt.Println(message)
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

			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
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
			// fmt.Println("ping")
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
