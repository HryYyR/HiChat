package models

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaim struct {
	ID       int
	UUID     string
	UserName string
	jwt.RegisteredClaims
}

type ResponseUserData struct {
	ID          int
	UserName    string
	NikeName    string
	Email       string
	CreatedTime time.Time
	LoginTime   string
	Avatar      string
	GroupList   []GroupDetail
}

type Users struct {
	ID        int    `xorm:"pk autoincr notnull index"`
	UUID      string `xorm:"notnull unique"`
	UserName  string `xorm:"notnull"`
	NikeName  string `xorm:"notnull"`
	Password  string `xorm:"notnull"`
	Email     string `xorm:"notnull"`
	Salt      string `xorm:"notnull"`
	IP        string
	Avatar    string
	City      string
	Age       int
	Grade     int       `xorm:"default(1)"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted"`
	UpdatedAt time.Time `xorm:"updated"`
	LoginTime string    `xorm:"updated"`
	LoginOut  time.Time
}

func (u *Users) TableName() string {
	return "users"
}

// 获取用户的群列表
func (u *Users) GetUserGroupList() ([]GroupDetail, error) {
	session := adb.Ssql.NewSession()
	var usergouplist []GroupDetail

	// 查询用户加入的群列表(没有详情)
	var gur []GroupUserRelative
	if err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).Find(&gur); err != nil {
		fmt.Println("查询用户加入的群列表error:", err)
		session.Rollback()
		return []GroupDetail{}, err
	}

	// 查询用户的所有消息
	var usermessagelist []GroupMessage
	if err := adb.Ssql.Table("group_message").Find(&usermessagelist); err != nil {
		fmt.Println("查询所有消息error:", err)
		session.Rollback()
		return []GroupDetail{}, err
	}

	for _, g := range gur {
		var group Group                //群详情
		var messagelist []GroupMessage //群消息列表

		//  根据群id查询群的详细信息
		_, err := adb.Ssql.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
		if err != nil {
			fmt.Println("根据群id查询群的详细信息error:", err)
			session.Rollback()
			return []GroupDetail{}, err
		}
		// 将该群聊的消息放入消息列表
		for _, m := range usermessagelist {
			// fmt.Printf("%+v-----%+v\n", m.GroupID, g.ID)
			if m.GroupID == g.GroupID {
				messagelist = append(messagelist, m)
			}
		}
		var groupitem = GroupDetail{
			GroupInfo:   group,
			MessageList: messagelist,
		}
		usergouplist = append(usergouplist, groupitem)
	}
	session.Commit()
	return usergouplist, nil
}
