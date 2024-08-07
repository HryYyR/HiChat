package config

import "time"

// var CallNoticeVideoStreamServerIP = "192.168.137.1"
var CallNoticeVideoStreamServerIP = "127.0.0.1"
var CallNoticeVideoStreamServerPort = "50052"

var ServerPort = 3004
var ConsulAddress = "127.0.0.1:8500"

var JwtKey = "Hyyyh1527"

var MysqlMaxIdleConns = 1000
var MysqlMaxOpenConns = 2000

// var MysqlAddress = "host.docker.internal:3306" //docker
var MysqlAddress = "localhost:3306" //localhost

var MysqlUserName = "root"
var MysqlPassword = "root"
var MysqlDatabase = "go_websocket"

var RedisAddr = "localhost:6379"
var RedisPassword = ""
var RedisDB = 0

var RabbitMQAddress = "amqp://guest:guest@localhost:5672/"

var EmailAccount = "2452719312@qq.com"
var EmailPassword = "hdabghzavlyeeajj"

var WriteWait = 10 * time.Second    //socket写入返回超时时间
var ResponseWait = 60 * time.Second //socket反应返回超时时间

var HeartbeatTicker = (ResponseWait * 9) / 10 //心跳检测间隔时间

var MaxMessageSize = int64(1024 * 1024 * 10) //消息最大容量 10m

var FlowControlTime = 1 * time.Minute //接口限流每周期时间
var FlowControlNum = 300              //接口限流每周期最大访问次数

var MsgTypeDefault = 1 //群聊文字消息
var MsgTypeImage = 2   //群聊图片消息
var MsgTypeAudio = 3   //群聊音频消息

var MsgTypeQuitGroup = 201      //退出群聊
var MsgTypeJoinGroup = 202      //加入群聊
var MsgTypeApplyJoinGroup = 203 //申请加入群聊
var MsgTypeDissolveGroup = 204  //解散群聊

var MsgTypeSyncMsg = 400      //同步消息
var MsgTypeClearSyncMsg = 401 //同步消息清零

// 用于通知的消息类型,不用于消息传输

var MsgTypeRefreshGroupAndNotice = 500 //刷新群聊列表和群聊通知列表
var MsgTypeRefreshGroupNotice = 502    //刷新群聊通知列表

var MsgTypeFriendDefault = 1001 //好友文字消息
var MsgTypeFriendImage = 1002   //好友图片消息
var MsgTypeFriendAudio = 1003   //好友音频消息

var MsgTypeRefreshFriendNotice = 1005    //刷新好友通知列表
var MsgTypeRefreshFriend = 1006          //刷新好友列表
var MsgTypeRefreshFriendAndNotice = 1007 //刷新好友列表和通知列表

var MsgTypeSyncFriendMsg = 1400      //同步好友消息
var MsgTypeClearSyncFriendMsg = 1401 //同步好友消息清零

var MsgTypeStartUserToUserVideoCall = 1501 //用户之间发起视频通话
