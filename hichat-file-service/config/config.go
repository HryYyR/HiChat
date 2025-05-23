package config

import (
	"os"
	"time"
)

var ENV = "dev"
var ServerName = "hichat-file-server"

var ServerPort = 3006
var ConsulAddress = "127.0.0.1:8500"

var JwtKey = "Hyyyh1527"

// var MysqlAddress = "host.docker.internal:3306" //docker
var MysqlAddress = "127.0.0.1:3306" //127.0.0.1
var MysqlUserName = "root"
var MysqlPassword = "root"
var MysqlDatabase = "go_websocket"

var RedisAddr = "127.0.0.1:6379"
var RedisPassword = ""
var RedisDB = 0

var EmailAccount = "2452719312@qq.com"
var EmailPassword = "hdabghzavlyeeajj"

var WriteWait = 10 * time.Second    //socket写入返回超时时间
var ResponseWait = 60 * time.Second //socket反应返回超时时间

var HeartbeatTicker = (ResponseWait * 9) / 10 //心跳检测间隔时间

var MaxMessageSize = int64(1024 * 1024) //消息最大容量

var MsgTypeDefault = 1 //群聊文字消息
var MsgTypeImage = 2   //群聊图片消息
var MsgTypeAudio = 3   //群聊音频消息

var MsgTypeRefreshGroup = 200 //刷新群聊
var MsgTypeQuitGroup = 201    //退出群聊
var MsgTypeSyncMsg = 400      //同步消息
var MsgTypeClearSyncMsg = 401 //同步消息清零

func SetEnvironment(env string) {
	envJwtKey, exists := os.LookupEnv("JwtKey")
	if exists {
		JwtKey = envJwtKey
	}
	envEmailAccount, exists := os.LookupEnv("EmailAccount")
	if exists {
		EmailAccount = envEmailAccount
	}
	envEmailPassword, exists := os.LookupEnv("EmailPassword")
	if exists {
		EmailPassword = envEmailPassword
	}

	switch env {
	case "docker":
		MysqlAddress = "hichat-mysql:3306"
		ConsulAddress = "hichat-consul:8500"
		RedisAddr = "hichat-redis:6379"
	default:
	}
}
