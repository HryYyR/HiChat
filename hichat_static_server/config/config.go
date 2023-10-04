package config

import "time"

var JwtKey = "Hyyyh1527"

var MysqlUserName = "root"
var MysqlPassword = "root"
var MysqlDatabase = "go_websocket"

var RedisAddr = "localhost:6379"
var RedisPassword = ""
var RedisDB = 0

var EmailAccount = "2452719312@qq.com"
var EmailPassword = "hdabghzavlyeeajj"

var WriteWait = 10 * time.Second    //socket写入返回超时时间
var ResponseWait = 60 * time.Second //socket反应返回超时时间

var HeartbeatTicker = (ResponseWait * 9) / 10 //心跳检测间隔时间

var MaxMessageSize = int64(1024 * 1024) //消息最大容量
