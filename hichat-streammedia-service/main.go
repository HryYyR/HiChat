package main

import (
	"flag"
	"fmt"
	"hichat-streammedia-serivce/models"
	"hichat-streammedia-serivce/service"

	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	engine.GET("/ws", service.Connectws) //用户连接
	models.ServiceCenter = models.NewHub("1")
	go models.ServiceCenter.Run()
	fmt.Println("service run in ", 3009)
	err := engine.Run(":3009")
	if err != nil {
		fmt.Println(err)
		return
	}
}
