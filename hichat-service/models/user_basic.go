package models

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

// 获取用户的群列表及群消息
func (u *Users) GetUserGroupList() ([]GroupDetail, error) {
	var usergouplist []GroupDetail

	// 查询用户加入的群列表(没有详情)
	var gur []GroupUserRelative
	if err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).Find(&gur); err != nil {
		fmt.Println("查询用户加入的群列表error:", err)
		return []GroupDetail{}, err
	}

	// 查询用户的所有消息
	var usermessagelist []GroupMessage
	if err := adb.Ssql.Table("group_message").Find(&usermessagelist); err != nil {
		fmt.Println("查询所有消息error:", err)
		return []GroupDetail{}, err
	}

	for _, g := range gur {
		var group Group                //群详情
		var messagelist []GroupMessage //群消息列表

		//  根据群id查询群的详细信息
		_, err := adb.Ssql.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
		if err != nil {
			fmt.Println("根据群id查询群的详细信息error:", err)
			return []GroupDetail{}, err
		}

		// 查询该用户加入的每个群聊的人数
		membercount, err := adb.Ssql.Table("group_user_relative").Where("group_id=?", g.GroupID).Count()
		if err != nil {
			fmt.Println("查询用户加入的每个群聊的人数error:", err)
			return []GroupDetail{}, err
		}
		group.MemberCount = int(membercount)

		var unreadmsgdata GroupUnreadMessage
		has, err := adb.Ssql.Table("group_unread_message").Where("group_id=? and user_id=?", g.GroupID, g.UserID).Get(&unreadmsgdata)
		if err != nil {
			fmt.Println("查询用户的未读数error:", err)
			return []GroupDetail{}, err
		}
		if has {
			group.UnreadMessage = unreadmsgdata.UnreadNumber
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
	return usergouplist, nil
}

// 获取用户的申请消息通知
func (u *Users) GetApplyMsgList(applylist *[]ApplyJoinGroup) error {
	var grouplist []Group
	if err := adb.Ssql.Table("group").Where("creater_id=?", u.ID).Find(&grouplist); err != nil {
		return err
	}

	for _, g := range grouplist {
		var applyjoingrouplist []ApplyJoinGroup
		if err := adb.Ssql.Table("apply_join_group").Where("group_id=?", g.ID).Find(&applyjoingrouplist); err != nil {
			return err
		}
		*applylist = append(*applylist, applyjoingrouplist...)
	}

	return nil
}
func (u *Users) GetApplyAddUserList(applylist *[]ApplyAddUser) error {
	var userlist []ApplyAddUser
	if err := adb.Ssql.Table("apply_add_user").Where("pre_apply_user_id=? or apply_user_id=?", u.ID, u.ID).Find(&userlist); err != nil {
		return err
	}
	*applylist = userlist
	return nil
}

// 获取好友列表
func (u *Users) GetFriendList(friendlist *[]Friend) error {
	var friendrelativelist []UserUserRelative
	err := adb.Ssql.Table("user_user_relative").Where("pre_user_id = ? or back_user_id=?", u.ID, u.ID).Find(&friendrelativelist)
	if err != nil {
		return err
	}

	for _, relative := range friendrelativelist {
		var frienddata Users
		if relative.BackUserID == u.ID {
			exit, err := adb.Ssql.Table("users").Where("id = ?", relative.PreUserID).Get(&frienddata)
			if err != nil {
				fmt.Println("获取用户信息失败")
				return err
			}
			if !exit {
				continue
			}
		} else {
			exit, err := adb.Ssql.Table("users").Where("id = ?", relative.BackUserID).Get(&frienddata)
			if err != nil {
				fmt.Println("获取用户信息失败")
				return err
			}
			if !exit {
				continue
			}
		}
		data := Friend{
			Id:        int32(frienddata.ID),
			UserName:  frienddata.UserName,
			NikeName:  frienddata.NikeName,
			Email:     frienddata.Email,
			Avatar:    frienddata.Avatar,
			City:      frienddata.City,
			Age:       strconv.Itoa(frienddata.Age),
			CreatedAt: frienddata.CreatedAt,
			DeletedAt: frienddata.DeletedAt,
			UpdatedAt: frienddata.UpdatedAt,
		}
		*friendlist = append(*friendlist, data)

	}
	return nil
}

// type ResponseUserData struct {
// 	ID          int
// 	UserName    string
// 	NikeName    string
// 	Email       string
// 	CreatedTime time.Time
// 	LoginTime   string
// 	GroupList   []GroupDetail
// }

type GroupDetail struct {
	GroupInfo   Group
	MessageList []GroupMessage
}

type UserClaim struct {
	ID       int
	UUID     string
	UserName string
	jwt.RegisteredClaims
}
