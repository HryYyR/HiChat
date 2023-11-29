package models

//import (
//	adb "HiChat/hichat-mq-service/ADB"
//	"fmt"
//	"time"
//)
//
//type RedisMessage struct {
//	MsgType    int
//	Key        string
//	Value      string
//	Expiration time.Duration //默认分钟
//}
//
//func (r *RedisMessage) RedisDelKey() error {
//	err := adb.Rediss.Del(r.Key).Err()
//	if err != nil {
//		fmt.Println(err.Error())
//		return err
//	}
//	return nil
//}
//
//func (r *RedisMessage) RedisSetString() error {
//	err := adb.Rediss.SetNX(r.Key, r.Value, r.Expiration).Err()
//	if err != nil {
//		fmt.Println(err.Error())
//		return err
//	}
//	return nil
//}
//
//func (r *RedisMessage) RedisRpushList() error {
//	err := adb.Rediss.RPush(r.Key, r.Value).Err()
//	if err != nil {
//		fmt.Println(err.Error())
//		return err
//	}
//	return nil
//}
