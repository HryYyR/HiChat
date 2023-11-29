package service

import (
	"fmt"
	"go-websocket-server/models"
	"go-websocket-server/util"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 初始化将http升级为ws协议的配置信息
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 连接ws
func Connectws(c *gin.Context) {
	token := c.Query("token")

	userdata, err := util.DecryptToken(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("用户%v加入了房间\n", userdata.ID)

	grouplist, err := models.GetUserGroupList(userdata.ID)
	uuid := util.GenerateUUID()
	if err != nil {
		fmt.Println(err)
		return
	}

	Conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) //升级协议
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}

	client := models.UserClient{
		ClientID: uuid,
		UserID:   userdata.ID,
		UserUUID: userdata.UUID,
		UserName: userdata.UserName,
		Conn:     Conn,
		Send:     make(chan []byte, 256),
		Status:   true,
		Groups:   grouplist,
		// CachingMessages: make(map[int]int, 0),
		Mutex: &sync.RWMutex{},
	}

	models.ServiceCenter.Mutex.Lock()
	models.ServiceCenter.Clients[userdata.ID] = client
	models.ServiceCenter.Mutex.Unlock()

	go client.WritePump()
	go client.ReadPump()
}
