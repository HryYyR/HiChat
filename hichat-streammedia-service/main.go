package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"hichat-streammedia-service/config"
	"hichat-streammedia-service/models"
	"hichat-streammedia-service/rpcserver"
	"hichat-streammedia-service/service"
	"hichat-streammedia-service/service_registry"
	"hichat-streammedia-service/util"
	"log"
	"time"
)

func main() {
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	engine.GET("/ws", service.Connectws) //用户连接
	models.ServiceCenter = models.NewHub("1")
	go models.ServiceCenter.Run()

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fmt.Println(time.Now().Format(time.RFC3339))
			fmt.Print("在线房间: ")
			for roomid, _ := range models.ServiceCenter.Room {
				fmt.Print(roomid + " ")
			}
			fmt.Println("")
			fmt.Print("在线用户: ")
			for uid, c := range models.ServiceCenter.Clients {
				fmt.Printf("%v(%v) ", uid, c.BelongRoomUUID)
			}
			fmt.Println("")
		}
	}()

	//服务注册
	dis := service_registry.DiscoveryConfig{
		ID:      util.GenerateUUID(),
		Name:    "hichat-streammedia-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: util.GetIP(),
	}
	err := service_registry.RegisterService(dis)
	if err != nil {
		log.Fatalln(err)
	}

	go rpcserver.ListenNoticeVideoStreamRpcServer()

	fmt.Println("service run in ", 3009)
	err = engine.Run(":3009")
	if err != nil {
		fmt.Println(err)
		return
	}
}
