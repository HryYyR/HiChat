// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	adb "go-websocket-server/ADB"
	systeminit "go-websocket-server/SystemInit"
	_ "go-websocket-server/log"
	"go-websocket-server/models"
	"go-websocket-server/rpcserver"
	"go-websocket-server/service"
	"go-websocket-server/util"

	"github.com/gin-gonic/gin"
)

func main() {

	flag.Parse()

	go systeminit.PrintRoomInfo()

	adb.InitMySQL()
	adb.InitRedis()
	if err := systeminit.InitClientsToGrouplist(); err != nil {
		panic(err.Error())
	}

	models.ServiceCenter = models.NewHub(util.GenerateUUID())
	go models.ServiceCenter.Run()

	if err := systeminit.InitUserToClient(); err != nil {
		panic(err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())
	// engine.POST("/login", service.Login)         //登录
	engine.POST("/register", service.Register)   //注册
	engine.POST("/emailcode", service.EmailCode) //邮箱验证码

	engine.GET("/ws", service.Connectws) //用户连接

	// wsgroup := engine.Group("ws")
	// wsgroup.GET("/:roomid", service.ServeWs) //用户进入房间

	usergroup := engine.Group("user", service.IdentityCheck)
	usergroup.POST("/creategroup", service.CreateGroup)           //创建群聊
	usergroup.POST("/joingroup", service.JoinGroup)               //加入群聊
	usergroup.POST("/exitgroup", service.ExitGroup)               //退出群聊
	usergroup.POST("/RefreshGroupList", service.RefreshGroupList) //获取用户信息
	usergroup.POST("/searchGroup", service.SearchGroup)           //获取用户信息

	go rpcserver.ListenGetUserGroupListRpcServer()

	fmt.Println("service run in 3004")
	engine.Run(":3004")

}
