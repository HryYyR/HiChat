package main

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/config"
	_ "hichat_static_server/log"
	"hichat_static_server/service"
	"hichat_static_server/service_registry"
	"hichat_static_server/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	adb.InitMySQL()
	adb.InitRedis()
	adb.InitMQ()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())
	engine.POST("/login", service.Login)         //登录
	engine.POST("/register", service.Register)   //注册
	engine.POST("/emailcode", service.EmailCode) //邮箱验证码

	engine.POST("/test", service.Test) //test

	usergroup := engine.Group("user", service.IdentityCheck)
	usergroup.POST("/edituserdata", service.EditUserData) //修改用户信息
	usergroup.POST("/getuserdata", service.GetUserData)   //获取用户基本信息
	usergroup.POST("/searchuser", service.SearchUser)     //搜索用户

	usergroup.POST("/getusergrouplist", service.GetUserGroupList)   //获取用户的群聊列表
	usergroup.POST("/getuserfriendlist", service.GetUserFriendList) //获取用户的好友列表

	usergroup.POST("/getuserapplyaddfriendlist", service.GetUserApplyAddFriendList) //获取用户的好友申请列表
	usergroup.POST("/getuserapplyjoingrouplist", service.GetUserApplyJoinGroupList) //获取用户的群聊通知列表

	groupgroup := engine.Group("group", service.IdentityCheck)
	groupgroup.POST("/searchgroup", service.SearchGroup)                 //搜索群聊
	groupgroup.POST("/getgroupmessagelist", service.GetGroupMessageList) //获取指定群聊的消息(限定条数)

	//服务注册
	addressIP := util.GetIP()
	dis := service_registry.DiscoveryConfig{
		ID:      util.GenerateUUID(),
		Name:    "hichat-static-server",
		Tags:    nil,
		Port:    config.ServerPort,
		Address: addressIP,
	}
	err := service_registry.RegisterService(dis)
	if err != nil {
		panic(err)
	}

	fmt.Println("service run in ", config.ServerPort)
	serverpost := fmt.Sprintf(":%s", strconv.Itoa(config.ServerPort))
	err = engine.Run(serverpost)
	if err != nil {
		panic(err)
	}
}
