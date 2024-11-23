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
	var port int
	var env string
	flag.IntVar(&port, "p", config.ServerPort, "端口号")
	flag.StringVar(&env, "d", config.ENV, "运行环境")

	flag.Parse()

	config.SetEnvironment(env)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	engine.GET("/streammedia", service.Connectws) //用户连接
	models.ServiceCenter = models.NewHub("1")
	go models.ServiceCenter.Run()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
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
		ID:   util.GenerateUUID(),
		Name: config.ServerName,
		Tags: []string{
			"traefik.enable=true",
			"traefik.http.routers.streammedia-router.rule=PathPrefix(`/streammedia`)",
			fmt.Sprintf("traefik.http.services.%s.loadBalancer.server.port=%d", config.ServerName, config.ServerPort),
		}, // 标签开启服务暴露
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
