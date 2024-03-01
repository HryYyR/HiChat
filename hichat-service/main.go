// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	systeminit "go-websocket-server/SystemInit"
	"go-websocket-server/config"
	_ "go-websocket-server/log"
	"go-websocket-server/models"
	"go-websocket-server/service"
	"go-websocket-server/service_registry"
	"go-websocket-server/util"
	util2 "go-websocket-server/util"
	"log"
	"strconv"
)

func main() {

	flag.Parse()

	adb.InitMySQL()
	adb.InitRedis()
	adb.InitMQ()
	if err := systeminit.InitClientsToGrouplist(); err != nil {
		fmt.Println(err)
		log.Println(err)
		panic(err.Error())
	}

	models.ServiceCenter = models.NewHub(util.GenerateUUID())
	go models.ServiceCenter.Run()

	if err := systeminit.InitUserToClient(); err != nil {
		fmt.Println(err)
		log.Println(err)
		panic(err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	engine.GET("/ws", service.Connectws) //用户连接

	usergroup := engine.Group("user", service.IdentityCheck)
	usergroup.POST("/creategroup", service.CreateGroup)         //创建群聊
	usergroup.POST("/handlejoingroup", service.HandleJoinGroup) //处理加入群聊
	usergroup.POST("/applyjoingroup", service.ApplyJoinGroup)   //申请加入群聊
	usergroup.POST("/exitgroup", service.ExitGroup)             //退出群聊
	usergroup.POST("/searchGroup", service.SearchGroup)         //搜索群聊

	// usergroup.POST("/RefreshGroupList", service.RefreshGroupList) //获取用户信息

	usergroup.POST("/applyadduser", service.ApplyAddUser)                         //申请添加好友
	usergroup.POST("/handleadduser", service.HandleAddUser)                       //处理添加好友
	usergroup.POST("/startusertouservideocall", service.StartUserToUserVideoCall) //检查指定用户登录状态

	//go rpcserver.ListenGetUserGroupListRpcServer()

	//服务注册
	addressIP := util2.GetIP()
	dis := service_registry.DiscoveryConfig{
		ID:      util2.GenerateUUID(),
		Name:    "hichat-ws-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: addressIP,
	}
	err := service_registry.RegisterService(dis)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("service run in ", config.ServerPort)
	serverpost := fmt.Sprintf(":%s", strconv.Itoa(config.ServerPort))
	err = engine.Run(serverpost)
	if err != nil {
		log.Fatalln(err)
	}
}
