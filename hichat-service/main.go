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
	"strconv"
)

func main() {

	flag.Parse()

	adb.InitMySQL()
	adb.InitRedis()
	adb.InitMQ()
	defer func(SqlStruct *adb.Sql) {
		err := SqlStruct.CloseConn()
		if err != nil {
			log.Println("mysql close error: ", err)
		}
	}(adb.SqlStruct)

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
	dis := service_registry.DiscoveryConfig{
		ID:      util2.GenerateUUID(),
		Name:    "hichat-ws-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: util2.GetIP(),
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
