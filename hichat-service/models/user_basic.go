package models

import (
	"github.com/golang-jwt/jwt/v4"
	adb "go-websocket-server/ADB"
	"time"
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
	Introduce string
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

// CheckUserExit 检查用户是否存在
func (u *Users) CheckUserExit() (Users, bool, error) {
	var applyuserdata Users
	has, err := adb.SqlStruct.Conn.Table("users").Where("id=?", u.ID).Get(&applyuserdata)
	if err != nil {
		return applyuserdata, false, err
	}
	if !has {
		return applyuserdata, false, nil
	}
	return applyuserdata, true, nil
}

type UserClaim struct {
	ID       int
	UUID     string
	UserName string
	jwt.RegisteredClaims
}

// GetUserGroupList 获取用户的群列表及群消息
//func (u *Users) GetUserGroupList() ([]GroupDetail, error) {
//	var usergouplist []GroupDetail
//
//	// 查询用户加入的群列表(没有详情)
//	var gur []GroupUserRelative
//	if err := adb.SqlStruct.Conn.Table("group_user_relative").Where("user_id=?", u.ID).Find(&gur); err != nil {
//		fmt.Println("查询用户加入的群列表error:", err)
//		return []GroupDetail{}, err
//	}
//
//	// 查询用户的所有消息
//	var usermessagelist []GroupMessage
//	if err := adb.SqlStruct.Conn.Table("group_message").Find(&usermessagelist); err != nil {
//		fmt.Println("查询所有消息error:", err)
//		return []GroupDetail{}, err
//	}
//
//	for _, g := range gur {
//		var group Group                //群详情
//		var messagelist []GroupMessage //群消息列表
//
//		//  根据群id查询群的详细信息
//		_, err := adb.SqlStruct.Conn.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
//		if err != nil {
//			fmt.Println("根据群id查询群的详细信息error:", err)
//			return []GroupDetail{}, err
//		}
//
//		// 查询该用户加入的每个群聊的人数
//		membercount, err := adb.SqlStruct.Conn.Table("group_user_relative").Where("group_id=?", g.GroupID).Count()
//		if err != nil {
//			fmt.Println("查询用户加入的每个群聊的人数error:", err)
//			return []GroupDetail{}, err
//		}
//		group.MemberCount = int(membercount)
//
//		var unreadmsgdata GroupUnreadMessage
//		has, err := adb.SqlStruct.Conn.Table("group_unread_message").Where("group_id=? and user_id=?", g.GroupID, g.UserID).Get(&unreadmsgdata)
//		if err != nil {
//			fmt.Println("查询用户的未读数error:", err)
//			return []GroupDetail{}, err
//		}
//		if has {
//			group.UnreadMessage = unreadmsgdata.UnreadNumber
//		}
//
//		// 将该群聊的消息放入消息列表
//		for _, m := range usermessagelist {
//			// fmt.Printf("%+v-----%+v\n", m.GroupID, g.ID)
//			if m.GroupID == g.GroupID {
//				messagelist = append(messagelist, m)
//			}
//		}
//		var groupitem = GroupDetail{
//			GroupInfo:   group,
//			MessageList: messagelist,
//		}
//		usergouplist = append(usergouplist, groupitem)
//	}
//	return usergouplist, nil
//}

// GetApplyMsgList 获取用户的申请消息通知
//func (u *Users) GetApplyMsgList(applylist *[]ApplyJoinGroup) error {
//	var grouplist []Group
//	if err := adb.SqlStruct.Conn.Table("group").Where("creater_id=?", u.ID).Find(&grouplist); err != nil {
//		return err
//	}
//
//	for _, g := range grouplist {
//		var applyjoingrouplist []ApplyJoinGroup
//		if err := adb.SqlStruct.Conn.Table("apply_join_group").Where("group_id=?", g.ID).Find(&applyjoingrouplist); err != nil {
//			return err
//		}
//		*applylist = append(*applylist, applyjoingrouplist...)
//	}
//
//	return nil
//}

// GetApplyAddUserList 获取申请列表
//func (u *Users) GetApplyAddUserList(applylist *[]ApplyAddUser) error {
//	var userlist []ApplyAddUser
//	if err := adb.SqlStruct.Conn.Table("apply_add_user").Where("pre_apply_user_id=? or apply_user_id=?", u.ID, u.ID).Find(&userlist); err != nil {
//		return err
//	}
//	*applylist = userlist
//	return nil
//}

// GetFriendList 获取好友列表
//func (u *Users) GetFriendList(friendlist *[]Friend) error {
//	var friendrelativelist []UserUserRelative
//	err := adb.SqlStruct.Conn.Table("user_user_relative").Where("pre_user_id = ? or back_user_id=?", u.ID, u.ID).Find(&friendrelativelist)
//	if err != nil {
//		return err
//	}
//	for _, relative := range friendrelativelist {
//		var frienddata Users
//		if relative.BackUserID == u.ID {
//			exit, err := adb.SqlStruct.Conn.Table("users").Where("id = ?", relative.PreUserID).Get(&frienddata)
//			if err != nil {
//				fmt.Println("获取用户信息失败")
//				return err
//			}
//			if !exit {
//				continue
//			}
//		} else {
//			exit, err := adb.SqlStruct.Conn.Table("users").Where("id = ?", relative.BackUserID).Get(&frienddata)
//			if err != nil {
//				fmt.Println("获取用户信息失败")
//				return err
//			}
//			if !exit {
//				continue
//			}
//		}
//		data := Friend{
//			Id:        int32(frienddata.ID),
//			UserName:  frienddata.UserName,
//			NikeName:  frienddata.NikeName,
//			Email:     frienddata.Email,
//			Avatar:    frienddata.Avatar,
//			City:      frienddata.City,
//			Age:       strconv.Itoa(frienddata.Age),
//			CreatedAt: frienddata.CreatedAt,
//			DeletedAt: frienddata.DeletedAt,
//			UpdatedAt: frienddata.UpdatedAt,
//		}
//		*friendlist = append(*friendlist, data)
//	}
//	return nil
//}

// type ResponseUserData struct {
// 	ID          int
// 	UserName    string
// 	NikeName    string
// 	Email       string
// 	CreatedTime time.Time
// 	LoginTime   string
// 	GroupList   []GroupDetail
// }
