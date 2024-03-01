package config

import (
	"time"
)

var ServerPort = 3005
var ConsulAddress = "127.0.0.1:8500"

var LoginRpcAddr = ":8080"
var LoginHttpAddr = ":8090"

var JwtKey = "Hyyyh1527"

var MysqlAddress = "localhost:3306" //localhost
//var MysqlAddress = os.Getenv("DB_HOST") + ":3306" //localhost
// var MysqlAddress = "host.docker.internal:3306" //docker

var MysqlUserName = "root"
var MysqlPassword = "root"
var MysqlDatabase = "go_websocket"

// var RedisAddr = "host.docker.internal:6379"
var RedisAddr = "localhost:6379"

// var RedisAddr = os.Getenv("REDIS_HOST") + ":6379"
var RedisPassword = ""
var RedisDB = 0

var RabbitMQAddress = "amqp://guest:guest@host.docker.internal:5672/"

//var RabbitMQAddress = "amqp://guest:guest@" + os.Getenv("MQ_HOST") + ":5672/"

var EmailAccount = "2452719312@qq.com"
var EmailPassword = "hdabghzavlyeeajj"

var WriteWait = 10 * time.Second    //socket写入返回超时时间
var ResponseWait = 60 * time.Second //socket反应返回超时时间

var HeartbeatTicker = (ResponseWait * 9) / 10 //心跳检测间隔时间

var MaxMessageSize = int64(1024 * 1024) //消息最大容量

var MsgTypeDefault = 1        //群聊文字消息
var MsgTypeRefreshGroup = 200 //刷新群聊
var MsgTypeQuitGroup = 201    //退出群聊
var MsgTypeSyncMsg = 400      //同步消息
var MsgTypeClearSyncMsg = 401 //同步消息清零

var MsgTypeRedisDelKey = 1601    //redis删除key
var MsgTypeRedisSetString = 1602 //redis设置string
var MsgTypeRedisRpushList = 1603 //redis向List添加元素
