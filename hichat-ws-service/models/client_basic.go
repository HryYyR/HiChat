package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-websocket-server/config"
	"go-websocket-server/util"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// UserClient 用户客户端
type UserClient struct {
	ClientID string
	UserID   int
	UserUUID string
	UserName string
	Conn     *websocket.Conn
	Status   bool
	Send     chan []byte
	Groups   map[int]Group //群聊列表  key:group_id  value:group
	// CachingMessages map[int]int   // key:group_id  value:未读数量
	Mutex            *sync.RWMutex // 互斥锁     多个结构体实例可以共享同一个锁时用指针,此处只会创建一个,所以不用指针
	HoldEncryptedKey bool          //是否持有key,没key不接收消息
	EncryptedKey     []byte
	Device           config.Device
	UserAgent        string
}

// ReadPump 读取用户发送的信息
func (c *UserClient) ReadPump() {
	defer func() {
		//fmt.Println("close reader")
		ServiceCenter.Loginout <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			//fmt.Println(err)
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ServiceCenter.Loginout <- c
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("IsUnexpectedCloseError: %v\n", err)
			}

			break
		}
		//如果没有key  阻止消息,,直到获取到key为止
		if !c.HoldEncryptedKey {
			//fmt.Println("正在验证key")
			encryptedKey, err := base64.StdEncoding.DecodeString(string(message))
			decryptRSABase64, err := util.DecryptRSA(encryptedKey, ServiceCenter.GetPrivateKey())
			//fmt.Println("验证成功 用户: ", c.UserName, "的aeskey: ", string(decryptRSABase64))
			AesKey, err := base64.StdEncoding.DecodeString(string(decryptRSABase64))
			if err != nil {
				log.Println(err)
				continue
			}
			c.EncryptedKey = AesKey
			c.HoldEncryptedKey = true
			continue
		}

		//fmt.Println(message)

		var data EncryptedData
		if err = json.Unmarshal(message, &data); err != nil {
			fmt.Println("Unmarshal rawData error ", err)
			continue
		}
		//fmt.Println("用户", c.UserName, " 的aeskey: ", base64.StdEncoding.EncodeToString(c.EncryptedKey))
		//fmt.Println("用户", c.UserName, " 发送的 Message: ", data.Message)
		//fmt.Println("用户", c.UserName, " 发送的 Iv: ", data.Iv)
		msg, err := base64.StdEncoding.DecodeString(data.Message)
		iv, err := base64.StdEncoding.DecodeString(data.Iv)
		decryptedData, err := util.DecryptAES(msg, iv, c.EncryptedKey)
		if err != nil {
			fmt.Println("DecryptAES error", err)
			continue
		}
		//fmt.Println("iv : ", data.Iv)
		//fmt.Println("Message: ", string(decryptedData))
		ServiceCenter.Broadcast <- decryptedData
	}
}

type EncryptedData struct {
	Iv      string
	Message string
}

// WritePump 给用户发送信息
func (c *UserClient) WritePump() {
	ticker := time.NewTicker(PingPeriod)

	defer func() {
		//fmt.Println("close writer")
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

			var AfterHandleMessage []byte
			var EncryptMessageJson util.EncryptedData
			if c.HoldEncryptedKey {
				//fmt.Printf("用户 %v 的aeskey是: %v\n", c.UserName, base64.StdEncoding.EncodeToString(c.EncryptedKey))
				//加密数据
				EncryptMessageJson, err = util.EncryptAESCBC(message, c.EncryptedKey)
				if err != nil {
					fmt.Println("encryptMessage Error", err)
					return
				}
				//fmt.Println("用户", c.UserName, " 的aeskey: ", base64.StdEncoding.EncodeToString(c.EncryptedKey))
				//fmt.Println("用户 ", c.UserName, "收到的base64 message", base64.StdEncoding.EncodeToString(EncryptMessageJson.Message))
				//fmt.Println("用户 ", c.UserName, "收到的base64 iv", base64.StdEncoding.EncodeToString(EncryptMessageJson.Iv))
				AfterHandleMessage, err = json.Marshal(EncryptMessageJson)
				if err != nil {
					fmt.Println("Marshal encryptMessage Error", err)
					return
				}
				//fmt.Println("数据已加密")
			} else {
				//fmt.Println("数据未加密")
				AfterHandleMessage = message
			}

			_, err = w.Write(AfterHandleMessage)
			if err != nil {
				fmt.Println("WriterError", err)
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
