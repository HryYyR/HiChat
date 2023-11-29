package adb

import (
	"HiChat/hichat-mq-service/config"

	"github.com/go-redis/redis"
)

var Rediss *redis.Client

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword, // no password set
		DB:       config.RedisDB,       // use default DB
	})
	Rediss = rdb
}

// func Init
