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
	var env string
	flag.IntVar(&port, "p", config.ServerPort, "端口号")
	flag.StringVar(&env, "d", config.ENV, "运行环境")
	flag.Parse()

	config.SetEnvironment(env)
	serveraddress := util2.GetIP() //生产环境使用
	//serveraddress := "192.168.137.1"

	//服务注册
	regsvconf := service_registry.DiscoveryConfig{
		ID:   util2.GenerateUUID(),
		Name: config.ServerName,
		Tags: []string{
			"traefik.enable=true",
			"traefik.http.routers.ws-router.rule=PathPrefix(`/ws`)",
			fmt.Sprintf("traefik.http.services.%s.loadBalancer.server.port=%d", config.ServerName, config.ServerPort),
		}, // 标签开启服务暴露
		Port:    port,
		Address: serveraddress,
	}
	err := service_registry.RegisterService(regsvconf)
	if err != nil {
		log.Fatalln(err)
	}
	err = adb.GetMySQLConn().Sync(new(models.UsersFile), new(models.Users), new(models.UserUserRelative), new(models.UserUnreadMessage), new(models.UserMessage), new(models.GroupUserRelative), new(models.GroupUnreadMessage), new(models.GroupMessage), new(models.Group), new(models.ApplyJoinGroup), new(models.ApplyAddUser))
	if err != nil {
		log.Println("同步数据表失败: ", err)
		panic(err)
	}

	//初始化基础资源
	adb.InitRedis()
	adb.MqHub.InitMQ()
	go models.RunReceiveMQMsg() //启动消费消息列表

	//消息处理中心
	models.ServiceCenter = models.NewHub(util.GenerateUUID())
	go models.ServiceCenter.Run()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	//初始化依赖
	sqlconn := adb.GetMySQLConn()
	userRepository := UsersScripts.NewUserRepository(sqlconn)
	groupRepository := GroupScripts.NewGroupRepository(sqlconn)

	//路由注册
	engine.Use(service.DependencyInjection(userRepository, groupRepository)) //依赖注入
	engine.GET("/ws", service.Connectws)                                     //用户连接
	usergrouprouter := engine.Group("ws/user", service.IdentityCheck, service.FlowControl)
	Route.InItUserGroupRouter(usergrouprouter)

	//启动服务
	serverPort := fmt.Sprintf(":%v", port)
	fmt.Println("服务运行在:  ", serveraddress, serverPort)
	err = engine.Run(serverPort)
	if err != nil {
		log.Fatalln(err)
	}
}
