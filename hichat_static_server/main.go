package main

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/service"

	"github.com/gin-gonic/gin"
)

func main() {
	adb.InitMySQL()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())
	engine.POST("/login", service.Login) //登录

	fmt.Println("service run in 3005")
	engine.Run(":3005")
}
