package main

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"hichat_static_server/service"

	"github.com/gin-gonic/gin"
)

func main() {
	adb.InitMySQL()
	adb.InitRedis()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(service.Cors())
	engine.POST("/login", service.Login)         //登录
	engine.POST("/register", service.Register)   //注册
	engine.POST("/emailcode", service.EmailCode) //邮箱验证码

	engine.POST("/test", service.Test) //test

	usergroup := engine.Group("user", service.IdentityCheck)
	usergroup.POST("/edituserdata", service.EditUserData) //邮箱验证码
	usergroup.POST("/getuserdata", service.GetUserData)   //获取用户基本信息
	usergroup.POST("/searchuser", service.SearchUser)     //搜索用户

	usergroup.POST("/getusergrouplist", service.GetUserGroupList)   //获取用户的群聊列表
	usergroup.POST("/getuserfriendlist", service.GetUserFriendList) //获取用户的好友列表

	usergroup.POST("/getuserapplyaddfriendlist", service.GetUserApplyAddFriendList) //获取用户的群聊通知列表
	usergroup.POST("/getuserapplyjoingrouplist", service.GetUserApplyJoinGroupList) //获取用户的好友申请列表

	groupgroup := engine.Group("group", service.IdentityCheck)
	groupgroup.POST("/searchgroup", service.SearchGroup)                 //搜索群聊
	groupgroup.POST("/getgroupmessagelist", service.GetGroupMessageList) //获取指定群聊的消息(限定条数)

	fmt.Println("service run in 3005")
	engine.Run(":3005")
}
