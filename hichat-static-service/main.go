package main

import (
	"flag"
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/config"
	_ "hichat_static_server/log"
	"hichat_static_server/service"
	"hichat_static_server/service_registry"
	"hichat_static_server/tool"
	"hichat_static_server/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	var port int
	var env string
	flag.IntVar(&port, "p", config.ServerPort, "")
	flag.StringVar(&env, "d", config.ENV, "运行环境")
	flag.Parse()

	config.SetEnvironment(env)

	fmt.Println(env)

	fmt.Println(config.ConsulAddress)

	adb.InitMySQL()
	adb.InitRedis()
	//adb.InitMQ()
	defer adb.NebulaInstance.CloseNebula()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())

	engine.POST("/test", service.Test) //test

	usergroup := engine.Group("user", service.IdentityCheck, service.FlowControl)
	usergroup.POST("/login", service.Login)               //登录
	usergroup.POST("/register", service.Register)         //注册
	usergroup.POST("/emailcode", service.EmailCode)       //邮箱验证码
	usergroup.POST("/edituserdata", service.EditUserData) //修改用户信息
	usergroup.POST("/getuserdata", service.GetUserData)   //获取用户基本信息
	usergroup.POST("/searchuser", service.SearchUser)     //搜索用户

	usergroup.POST("/getusergrouplist", service.GetUserGroupList)   //获取用户的群聊列表
	usergroup.POST("/getuserfriendlist", service.GetUserFriendList) //获取用户的好友列表

	usergroup.POST("/getuserapplyaddfriendlist", service.GetUserApplyAddFriendList) //获取用户的好友申请列表
	usergroup.POST("/getuserapplyjoingrouplist", service.GetUserApplyJoinGroupList) //获取用户的群聊通知列表

	usergroup.POST("/getusermessagelist", service.GetUserMessageList) //获取用户之间的消息(限定条数)

	usergroup.POST("/aimessage", service.AiMessage) //Ai问答

	groupgroup := engine.Group("group", service.IdentityCheck, service.FlowControl)
	groupgroup.POST("/searchgroup", service.SearchGroup)                 //搜索群聊
	groupgroup.POST("/getgroupmessagelist", service.GetGroupMessageList) //获取指定群聊的消息(限定条数)
	groupgroup.POST("/getgroupmemberdata", service.GetGroupMemberList)   //获取指定群聊的成员数据

	//服务注册
	addressIP := tool.GetIP()
	dis := service_registry.DiscoveryConfig{
		ID:   util.GenerateUUID(),
		Name: config.ServerName,
		Tags: []string{
			"traefik.enable=true",
			"traefik.http.routers.group-router.rule=PathPrefix(`/group`)",
			"traefik.http.routers.user-router.rule=PathPrefix(`/user`)",
			fmt.Sprintf("traefik.http.services.%s.loadBalancer.server.port=%d", config.ServerName, config.ServerPort),
		}, // 标签开启服务暴露
		Port:    config.ServerPort,
		Address: addressIP,
	}
	err := service_registry.ConsulRegisterService(dis)
	if err != nil {
		panic(err)
	}

	fmt.Println("service run in ", config.ServerPort)
	serverpost := fmt.Sprintf(":%s", strconv.Itoa(config.ServerPort))
	err = engine.Run(serverpost)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
