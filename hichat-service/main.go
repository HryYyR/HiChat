// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	"go-websocket-server/Route"
	systeminit "go-websocket-server/SystemInit"
	"go-websocket-server/config"
	_ "go-websocket-server/log"
	"go-websocket-server/models"
	"go-websocket-server/service"
	"go-websocket-server/service_registry"
	"go-websocket-server/util"
	util2 "go-websocket-server/util"
	"log"
)

func main() {

	flag.Parse()
	adb.InitMySQL()
	adb.InitRedis()
	adb.MqHub.InitMQ()
	defer func(SqlStruct *adb.Sql) {
		err := SqlStruct.CloseConn()
		if err != nil {
			log.Println("mysql close error: ", err)
		}
	}(adb.SqlStruct)
	go models.RunReceiveMQMsg() //启动消费消息列表

	if err := systeminit.InitClientsToGrouplist(); err != nil {
		log.Println(err)
		panic(err.Error())
	}

	models.ServiceCenter = models.NewHub(util.GenerateUUID())
	go models.ServiceCenter.Run()

	//todo delete
	if err := systeminit.InitUserToClient(); err != nil {
		log.Println(err)
		panic(err.Error())
	}

	//go systeminit.PrintRoomInfo()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())
	engine.GET("/ws", service.Connectws) //用户连接
	usergroup := engine.Group("user", service.IdentityCheck)
	Route.InItUserGroupRouter(usergroup)

	//go rpcserver.ListenGetUserGroupListRpcServer()

	//go func() {
	//	ticker := time.NewTicker(10 * time.Second)
	//	defer ticker.Stop()
	//	for range ticker.C {
	//		for gid, useridgoup := range models.GroupUserList {
	//			fmt.Printf("%v:%v \n", gid, useridgoup)
	//		}
	//	}
	//}()

	//服务注册
	regsvconf := service_registry.DiscoveryConfig{
		ID:      util2.GenerateUUID(),
		Name:    "hichat-ws-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: util2.GetIP(),
	}
	err := service_registry.RegisterService(regsvconf)
	if err != nil {
		log.Fatalln(err)
	}

	serverpost := fmt.Sprintf(":%d", config.ServerPort)
	log.Println("server run in ", util2.GetIP(), serverpost)
	err = engine.Run(serverpost)
	if err != nil {
		log.Fatalln(err)
	}
}
