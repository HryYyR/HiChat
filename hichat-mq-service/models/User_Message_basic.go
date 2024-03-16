package models

import (
	adb "HiChat/hichat-mq-service/ADB"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type UserMessage struct {
	ID                int    `xorm:"pk autoincr index"`
	UUID              string `xorm:"notnull"`
	UserID            int    `xorm:"notnull"`
	UserName          string `xorm:"notnull"`
	UserAvatar        string
	ReceiveUserID     int    `xorm:"notnull"`
	ReceiveUserName   string `xorm:"notnull"`
	ReceiveUserAvatar string
	Msg               string `xorm:"notnull"`
	MsgType           int
	IsReply           bool //是否是回复消息
	ReplyUserID       int  //如果是,被回复的用户id
	Context           []byte
	CreatedAt         time.Time `xorm:"created"`
	DeletedAt         time.Time `xorm:"deleted"`
	UpdatedAt         time.Time `xorm:"updated"`
}

func (u *UserMessage) SaveFriendMsgToDb() error {
	if _, err := adb.Ssql.Table("user_message").Insert(&u); err != nil {
		return err
	}
	var key string

	fmt.Println(u.UserID, u.ReceiveUserID)
	fmt.Println(u.UserID > u.ReceiveUserID)
	if u.UserID > u.ReceiveUserID {
		key = fmt.Sprintf("%d%d", u.UserID, u.ReceiveUserID)
	} else {
		key = fmt.Sprintf("%d%d", u.ReceiveUserID, u.UserID)
	}
	jsondata, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = adb.Rediss.RPush(key, string(jsondata)).Err()
	if err != nil {
		fmt.Println(errors.New("保存好友消息失败"))
	}

	return nil
}

func (u *UserMessage) SyncFriendMsgToDb() error {
	var data UserUnreadMessage
	exit, err := adb.Ssql.Table("user_unread_message").Where("user_id = ? and friend_id= ?", u.UserID, u.ReceiveUserID).Get(&data)
	if err != nil {
		return err
	}
	if exit {
		_, err := adb.Ssql.Table("user_unread_message").Where("user_id = ? and friend_id= ?", u.UserID, u.ReceiveUserID).Update(UserUnreadMessage{
			UnreadNumber: data.UnreadNumber + 1,
		})
		if err != nil {
			return err
		}
	} else {
		newmsg := UserUnreadMessage{
			UserName:     u.UserName,
			UserID:       u.UserID,
			FriendID:     u.ReceiveUserID,
			UnreadNumber: 1,
		}
		_, err := adb.Ssql.Table("user_unread_message").Insert(&newmsg)
		fmt.Println("insert ok")
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UserMessage) ClearFriendMsgNum() error {
	_, err := adb.Ssql.Table("user_unread_message").Cols("unread_number").Where("user_id = ? and friend_id=?", u.ReceiveUserID, u.UserID).Update(&UserUnreadMessage{
		UnreadNumber: 0,
	})
	fmt.Println("update ok")
	if err != nil {
		return err
	}

	return nil
}
