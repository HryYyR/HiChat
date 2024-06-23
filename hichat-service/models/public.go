package models

import (
	"time"
)

const (
	// WriteWait 定义了写操作的超时时间，为10秒
	WriteWait = 10 * time.Second
	// PongWait 定义了读操作的超时时间，为60秒
	PongWait = 60 * time.Second
	// PingPeriod 定义了发送ping消息的周期，为PongWait的9/10，用于保持连接活跃
	PingPeriod = (PongWait * 9) / 10
	// MaxMessageSize 定义了最大消息大小，为1MB
	MaxMessageSize = int64(1024 * 1024) //消息最大容量
)

var (
	// Newline 表示换行的字节序列
	Newline = []byte{'\n'}
	// Space 表示空格的字节序列
	Space = []byte{' '}
)
