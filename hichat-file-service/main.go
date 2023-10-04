package main

import (
	"fmt"
	adb "hichat-file-service/ADB"
	"hichat-file-service/middleware"
	"hichat-file-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	adb.InitMySQL()

	gin.SetMode("release")
	engine := gin.New()
	engine.Use(middleware.Cors())
	engine.POST("/uploadfile", service.UploadFile)

	engine.Static("/static", "./file")

	fmt.Println("service run in 3006")
	engine.Run(":3006")
}
