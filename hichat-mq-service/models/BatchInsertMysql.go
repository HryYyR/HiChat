package models

import (
	adb "HiChat/hichat-mq-service/ADB"
	"fmt"
	"log"
	"sync"
	"time"
)

type InsertMysqlBatch struct {
	UserMsgList  []UserMessage
	GroupMsgList []Message
	Mu           sync.Mutex
	Interval     time.Duration
	BatchSize    int //100
}

// AddUserMsg 添加用户消息
func (b *InsertMysqlBatch) AddUserMsg(msg *UserMessage) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	b.UserMsgList = append(b.UserMsgList, *msg)
}

// AddGroupMsg 添加群聊消息
func (b *InsertMysqlBatch) AddGroupMsg(msg *Message) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	b.GroupMsgList = append(b.GroupMsgList, *msg)
}

func (b *InsertMysqlBatch) InsertMessages() {
	ticker := time.NewTicker(b.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.Mu.Lock()
			if len(b.UserMsgList) != 0 {
				//批量插入用户消息
				Userbatch := b.UserMsgList
				if len(Userbatch) > b.BatchSize {
					Userbatch = Userbatch[:b.BatchSize]
				}
				b.UserMsgList = b.UserMsgList[len(Userbatch):]

				if _, err := adb.Ssql.Table("user_message").Insert(&Userbatch); err != nil {
					log.Printf("Failed to insert messages: %v", err)
				} else {
					log.Printf("用户消息插入 %d 条消息,剩余%d", len(Userbatch), len(b.UserMsgList))
				}
			}
			if len(b.GroupMsgList) != 0 {
				//批量插入群聊消息
				Groupbatch := b.GroupMsgList
				if len(Groupbatch) > b.BatchSize {
					Groupbatch = Groupbatch[:b.BatchSize]
				}
				b.GroupMsgList = b.GroupMsgList[len(Groupbatch):]

				if _, err := adb.Ssql.Table("group_message").Insert(&Groupbatch); err != nil {
					fmt.Println(err.Error())
				} else {
					log.Printf("群聊消息插入 %d 条消息,剩余%d", len(Groupbatch), len(b.GroupMsgList))
				}
			}
			b.Mu.Unlock()

		}
	}
}
