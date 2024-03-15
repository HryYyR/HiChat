package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 用户客户端
type UserClient struct {
	ClientID string
	UserID   int
	UserUUID string
	UserName string
	Conn     *websocket.Conn
	Status   bool
	Send     chan []byte
	Groups   map[int]Group //群聊列表  key:group_id  value:group
	// CachingMessages map[int]int   // key:groupid  value:未读数量
	Mutex *sync.RWMutex // 互斥锁     多个结构体实例可以共享同一个锁时用指针,此处只会创建一个,所以不用指针
}

// 读取用户发送的信息
func (c *UserClient) ReadPump() {
	defer func() {
		fmt.Println("close reader")
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
				fmt.Println("用户退出了")
				ServiceCenter.Loginout <- c
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("IsUnexpectedCloseError: %v\n", err)
			}
			break
		}

		// message = bytes.TrimSpace(bytes.Replace(message, Newline, Space, -1))
		ServiceCenter.Broadcast <- message
	}
}

// 给用户发送信息
func (c *UserClient) WritePump() {
	ticker := time.NewTicker(PingPeriod)

	defer func() {
		fmt.Println("close writer")
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		if c.Status == false {
			break
		}
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
			// fmt.Println("ping")
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
