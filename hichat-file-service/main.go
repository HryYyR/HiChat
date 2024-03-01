package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	adb "hichat-file-service/ADB"
	"hichat-file-service/config"
	"hichat-file-service/middleware"
	"hichat-file-service/service"
	"hichat-file-service/service_registry"
	"hichat-file-service/util"
)

func main() {
	adb.InitMySQL()

	gin.SetMode("release")
	engine := gin.New()
	engine.Use(middleware.Cors())
	g := engine.Group("file")
	g.POST("/uploadfile", service.UploadFile)

	engine.Static("/static", "./file")

	//服务注册
	addressIP := util.GetIP()
	dis := service_registry.DiscoveryConfig{
		ID:      util.GenerateUUID(),
		Name:    "hichat-file-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: addressIP,
	}
	err := service_registry.RegisterService(dis)
	if err != nil {
		panic(err)
	}

	fmt.Println("service run in 3006")
	engine.Run(":3006")
}
