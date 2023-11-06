package models

import (
	"fmt"
	adb "hichat_static_server/ADB"
	"sort"
	"strconv"
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
	ID            int
	UserName      string
	NikeName      string
	Email         string
	CreatedTime   time.Time
	LoginTime     string
	Avatar        string
	Age           int
	City          string
	GroupList     []GroupDetail
	ApplyList     []ApplyJoinGroup
	ApplyUserList []ApplyAddUser
	FriendList    []FriendResponse
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
func (u *Users) GetUserData(userdata *Users) error {
	var userinfo Users
	exit, err := adb.Ssql.Omit("Password,Salt,Grade,IP").Table("users").Where("id =?", u.ID).Get(&userinfo)
	if !exit {
		return fmt.Errorf("用户不存在")
	}
	if err != nil {
		return err
	}

	*userdata = userinfo
	return nil
}

// 获取用户的群列表
func (u *Users) GetUserGroupList(grouplist *[]GroupDetail) error {

	// 查询用户加入的群列表(没有详情)
	var gur []GroupUserRelative
	if err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).Find(&gur); err != nil {
		fmt.Println("查询用户加入的群列表error:", err)
		return err
	}

	// 查询用户的所有消息
	var usermessagelist []GroupMessage
	if err := adb.Ssql.Table("group_message").Desc("id").Find(&usermessagelist); err != nil {
		fmt.Println("查询所有消息error:", err)
		return err
	}
	//MessageTimeSort(usermessagelist, "desc")

	var unreadmsglist []GroupUnreadMessage
	if err := adb.Ssql.Table("group_unread_message").Where("user_id = ?", u.ID).Find(&unreadmsglist); err != nil {
		return err
	}
	unreadmsgmap := make(map[int]int, len(unreadmsglist)+1) //k:groupid   v:unreadnum
	for _, UnreadMessage := range unreadmsglist {
		unreadmsgmap[UnreadMessage.GroupID] = UnreadMessage.UnreadNumber
	}

	for _, g := range gur {
		var group Group                //群详情
		var messagelist []GroupMessage //群消息列表

		//  根据群id查询群的详细信息
		exit, err := adb.Ssql.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
		if !exit {
			continue
		}
		if err != nil {
			fmt.Println("根据群id查询群的详细信息error:", err)
			return err
		}

		group.UnreadMessage = unreadmsgmap[g.GroupID] //放入未读消息数量

		// 将该群聊的消息放入消息列表
		for _, m := range usermessagelist {
			// fmt.Printf("%+v-----%+v\n", m.GroupID, g.ID)
			if len(messagelist) >= 10 {
				break
			}
			if m.GroupID == g.GroupID {
				messagelist = append(messagelist, m)
			}
		}
		sort.Slice(messagelist, func(i, j int) bool { return messagelist[i].ID < (messagelist[j].ID) })
		//MessageTimeSort(messagelist, "asc")

		var groupitem = GroupDetail{
			GroupInfo:   group,
			MessageList: messagelist,
		}
		*grouplist = append(*grouplist, groupitem)
	}
	return nil
}

// 获取用户的群聊通知列表
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

// 获取用户的好友申请列表
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

func (u *Users) GetFriendListAndMEssage(friendresponselist *[]FriendResponse) error {
	var friendlist []Friend
	err := u.GetFriendList(&friendlist) //不带消息的好友列表
	if err != nil {
		return err
	}

	// 消息列表
	var messagelist []UserMessageItem
	err = adb.Ssql.Table("user_message").Omit("uuid,context").Where("user_id=? or receive_user_id=?", u.ID, u.ID).Find(&messagelist)
	if err != nil {
		return err
	}

	// 每个好友的未读数量
	var unreadmessagelist []UserUnreadMessage
	err = adb.Ssql.Table("user_unread_message").Cols("user_id,friend_id,unread_number").Where("friend_id=?", u.ID).Find(&unreadmessagelist)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", unreadmessagelist)

	// 映射map,方便操作
	unreadmessagemap := make(map[int]int, len(unreadmessagelist)+1) //key:friend_id  v:unread_number
	for _, m := range unreadmessagelist {
		unreadmessagemap[m.UserID] = m.UnreadNumber
	}
	fmt.Printf("%+v\n", unreadmessagemap)
	for _, resfriend := range friendlist {
		msglist := make([]UserMessageItem, 0)
		for _, msg := range messagelist {
			if msg.ReceiveUserID == int(resfriend.Id) || msg.UserID == int(resfriend.Id) {
				msglist = append(msglist, msg)
			}
		}

		resmsgitem := FriendResponse{
			Id:            resfriend.Id,
			UserName:      resfriend.UserName,
			NikeName:      resfriend.NikeName,
			Email:         resfriend.Email,
			Avatar:        resfriend.Avatar,
			City:          resfriend.City,
			Age:           resfriend.Age,
			UnreadMessage: unreadmessagemap[int(resfriend.Id)],
			MessageList:   msglist,
			CreatedAt:     resfriend.CreatedAt,
			DeletedAt:     resfriend.DeletedAt,
			UpdatedAt:     resfriend.UpdatedAt,
		}

		*friendresponselist = append(*friendresponselist, resmsgitem)

	}

	return nil
}

func (u *Users) Login(logindata *ResponseUserData) error {
	// 群聊通知列表
	applygrouplist := make([]ApplyJoinGroup, 0)
	err := u.GetApplyMsgList(&applygrouplist)
	if err != nil {
		return err
	}
	// 好友申请列表
	applyuserlist := make([]ApplyAddUser, 0)
	err = u.GetApplyAddUserList(&applyuserlist)
	if err != nil {
		return err
	}

	// 好友列表
	friendlist := make([]FriendResponse, 0)
	err = u.GetFriendListAndMEssage(&friendlist)
	if err != nil {
		return err
	}

	// 群聊列表
	grouplist := make([]GroupDetail, 0)
	err = u.GetUserGroupList(&grouplist)
	if err != nil {
		return err
	}

	var userdata Users
	err = u.GetUserData(&userdata)
	if err != nil {
		return err
	}

	*logindata = ResponseUserData{
		ID:            userdata.ID,
		UserName:      userdata.UserName,
		NikeName:      userdata.NikeName,
		Email:         userdata.Email,
		CreatedTime:   userdata.CreatedAt,
		LoginTime:     userdata.LoginTime,
		Avatar:        userdata.Avatar,
		Age:           userdata.Age,
		City:          userdata.City,
		GroupList:     grouplist,
		ApplyList:     applygrouplist,
		ApplyUserList: applyuserlist,
		FriendList:    friendlist,
	}

	return nil
}

func MessageTimeSort(arr []GroupMessage, s string) {
	if s == "desc" {
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].CreatedAt.After(arr[j].CreatedAt)
		})
	} else {
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].CreatedAt.Before(arr[j].CreatedAt)
		})
	}
}
