// Copyright 2023 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	adb "go-websocket-server/ADB"
	GroupScripts "go-websocket-server/ADB/MysqlScripts/GroupsScripts"
	"go-websocket-server/ADB/MysqlScripts/UsersScripts"
	"go-websocket-server/Route"
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
	var port int
	flag.IntVar(&port, "p", config.ServerPort, "端口号")
	flag.Parse()
	adb.InitRedis()
	adb.MqHub.InitMQ()
	go models.RunReceiveMQMsg() //启动消费消息列表
	//adb.InitMySQL()
	//defer func(SqlStruct *adb.Sql) {
	//	SqlStruct.CloseConn()
	//}(adb.SqlStruct)

	models.ServiceCenter = models.NewHub(util.GenerateUUID())
	go models.ServiceCenter.Run()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	sqlconn := adb.GetMySQLConn()
	userRepository := UsersScripts.NewUserRepository(sqlconn)
	groupRepository := GroupScripts.NewGroupRepository(sqlconn)

	engine.Use(service.DependencyInjection(userRepository, groupRepository)) //依赖注入
	engine.GET("/ws", service.Connectws)                                     //用户连接
	usergroup := engine.Group("ws/user", service.IdentityCheck, service.FlowControl)
	Route.InItUserGroupRouter(usergroup)

	serveraddress := util2.GetIP() //生产环境使用
	//serveraddress := "192.168.137.1"

	serverPort := fmt.Sprintf(":%v", port)

	//服务注册
	regsvconf := service_registry.DiscoveryConfig{
		ID:      util2.GenerateUUID(),
		Name:    "hichat-ws-server",
		Tags:    nil,
		Port:    port,
		Address: serveraddress,
	}
	err := service_registry.RegisterService(regsvconf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("服务运行在:  ", serveraddress, serverPort)
	err = engine.Run(serverPort)
	if err != nil {
		log.Fatalln(err)
	}
}
