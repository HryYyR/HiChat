package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hichat-streammedia-service/models"
	"hichat-streammedia-service/util"
	"log"
	"net/http"
	"strconv"
)

// 初始化将http升级为ws协议的配置信息
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, //5M
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 连接ws
func Connectws(c *gin.Context) {
	ustr := c.Query("id")
	uint, _ := strconv.Atoi(ustr)

	Conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //升级协议
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	fmt.Println("用户", ustr, "进入了")

	var RoomUUID string
	for k, v := range models.ServiceCenter.Room {
		if v.StartUserID == uint || v.ReceivedUserID == uint {
			RoomUUID = k
			break
		}
	}
	if len(RoomUUID) == 0 {
		util.H(c, http.StatusBadRequest, "非法访问", nil)
		return
	}

	client := models.UserClient{
		ClientID:       util.GenerateUUID(),
		UserID:         uint,
		UserName:       ustr,
		Conn:           Conn,
		Send:           make(chan []byte, 256),
		BelongRoomUUID: RoomUUID,
	}

	models.ServiceCenter.Clients[uint] = client

	go client.WritePump()
	go client.ReadPump()
}
