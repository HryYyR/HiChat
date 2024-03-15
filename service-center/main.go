package main

import (
	"HiChat/service-center/config"
	"HiChat/service-center/models"
	"HiChat/service-center/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	dis := models.DiscoveryConfig{
		ID:      Args[0],
		Name:    "server-center",
		Tags:    []string{},
		Port:    config.ConsulPort,
		Address: config.ConsulAddress, //通过ifconfig查看本机的eth0的ipv4地址
	}
	err = service.RegisterService(dis)
	if err != nil {
		panic(err)
	}

	fmt.Println("Service started to:", config.ServicePort)

	engine := gin.New()
	engine.Use(service.Cors())
}
