package models

import "time"

type RedisMessage struct {
	MsgType    int
	Key        string
	Value      string
	Expiration time.Duration //默认分钟
}
