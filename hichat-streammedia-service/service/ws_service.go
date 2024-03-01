package service

import (
	"fmt"
	"hichat-streammedia-serivce/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 初始化将http升级为ws协议的配置信息
var upgrader = websocket.Upgrader{
	ReadBufferSize:  5368709120, //5M
	WriteBufferSize: 5368709120,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 连接ws
func Connectws(c *gin.Context) {

	Conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //升级协议
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	fmt.Println("有用户进入了")

	client := models.UserClient{
		ClientID: "123",
		UserID:   123,
		UserName: "niko",
		// ClientID: uuid,
		// UserID:   userdata.ID,
		// UserName: userdata.UserName,
		Conn: Conn,
		Send: make(chan []byte, 256),
	}
	index := 1
	models.ServiceCenter.Clients[index] = client
	index++

	go client.WritePump()
	go client.ReadPump()
}
