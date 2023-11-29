package models

import (
	"fmt"
	"github.com/goinggo/mapstructure"
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
	Introduce     string
	GroupList     []GroupDetail
	ApplyList     []ApplyJoinGroupResponse
	ApplyUserList []ApplyAddUser
	FriendList    []FriendResponse
}

func (u *Users) SaveToRedis() error {
	_, err := adb.Rediss.HMSet(strconv.Itoa(u.ID), map[string]interface{}{
		"ID":          u.ID,
		"UserName":    u.UserName,
		"NikeName":    u.NikeName,
		"Email":       u.Email,
		"CreatedTime": u.CreatedAt.Local().String(),
		"LoginTime":   u.LoginTime,
		"Avatar":      u.Avatar,
		"Age":         u.Age,
		"City":        u.City,
		"Introduce":   u.Introduce,
		//"GroupList":   ru.GroupList,
		//"ApplyList"     []ApplyJoinGroupResponse
		//"ApplyUserList" []ApplyAddUser
		//"FriendList"    []FriendResponse
	}).Result()
	if err != nil {
		return err
	}
	adb.Rediss.Expire(strconv.Itoa(u.ID), time.Hour*360)
	return nil
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
func (u *Users) GetUserData(userdata *Users) error {
	var userinfo Users

	//从redis获取数据
	var udata = adb.Rediss.HGetAll(strconv.Itoa(u.ID)).Val()
	if len(udata) != 0 {
		_ = mapstructure.Decode(udata, &userinfo)
		userinfo.ID, _ = strconv.Atoi(udata["ID"])
		userinfo.Age, _ = strconv.Atoi(udata["Age"])
		//fmt.Printf("%+v", userinfo)
		*userdata = userinfo
		return nil
	}
	//fmt.Println("走数据库")

	exit, err := adb.Ssql.Omit("Password,Salt,Grade,IP").Table("users").Where("id =?", u.ID).Get(&userinfo)
	if !exit {
		return fmt.Errorf("用户不存在")
	}
	if err != nil {
		return err
	}

	err = u.SaveToRedis()
	if err != nil {
		return err
	}

	*userdata = userinfo
	return nil
}

// GetUserGroupList 获取用户的群列表
func (u *Users) GetUserGroupList(grouplist *[]GroupDetail) error {

	//未读数量
	var unreadmsglist []GroupUnreadMessage
	if err := adb.Ssql.Table("group_unread_message").Where("user_id = ?", u.ID).Find(&unreadmsglist); err != nil {
		return err
	}
	unreadmsgmap := make(map[int]int, len(unreadmsglist)+1) //k:groupid   v:unreadnum
	for _, UnreadMessage := range unreadmsglist {
		unreadmsgmap[UnreadMessage.GroupID] = UnreadMessage.UnreadNumber
	}

	var gur GroupUserRelative
	err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).
		Iterate(&gur, func(i int, bean interface{}) error {
			gur := bean.(*GroupUserRelative)

			//群的信息
			var group Group
			_, err := adb.Ssql.Table("group").Where("id=?", gur.GroupID).Get(&group)
			if err != nil {
				fmt.Println("根据群id查询群的详细信息error:", err)
				return err
			}
			group.UnreadMessage = unreadmsgmap[group.ID] //设置未读数量

			//群的消息
			tempdata := make([]GroupMessage, 0)
			if err := adb.Ssql.Table("group_message").Where("group_id=?", gur.GroupID).
				Desc("id").Limit(20).Find(&tempdata); err != nil {
				fmt.Println("查询群消息失败error:", err)
				return err
			}
			sort.Slice(tempdata, func(i, j int) bool {
				return tempdata[i].ID < tempdata[j].ID
			})

			var groupitem = GroupDetail{
				GroupInfo:   group,
				MessageList: tempdata,
			}
			*grouplist = append(*grouplist, groupitem)
			return nil
		})
	if err != nil {
		return err
	}
	return nil

	//
	//// 查询用户加入的群列表(没有详情)
	//var gur []GroupUserRelative
	//if err := adb.Ssql.Table("group_user_relative").Where("user_id=?", u.ID).Find(&gur); err != nil {
	//	fmt.Println("查询用户加入的群列表error:", err)
	//	return err
	//}
	//
	//// 查询用户的所有消息
	//usermessagelist := make([]GroupMessage, 0)
	//if err := adb.Ssql.Table("group_message").Desc("id").Find(&usermessagelist); err != nil {
	//	fmt.Println("查询所有消息error:", err)
	//	return err
	//}
	////MessageTimeSort(usermessagelist, "desc")
	//
	//var unreadmsglist []GroupUnreadMessage
	//if err := adb.Ssql.Table("group_unread_message").Where("user_id = ?", u.ID).Find(&unreadmsglist); err != nil {
	//	return err
	//}
	//unreadmsgmap := make(map[int]int, len(unreadmsglist)+1) //k:groupid   v:unreadnum
	//for _, UnreadMessage := range unreadmsglist {
	//	unreadmsgmap[UnreadMessage.GroupID] = UnreadMessage.UnreadNumber
	//}
	//
	//for _, g := range gur {
	//	var group Group                //群详情
	//	var messagelist []GroupMessage //群消息列表
	//
	//	//  根据群id查询群的详细信息
	//	exit, err := adb.Ssql.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
	//	if !exit {
	//		continue
	//	}
	//	if err != nil {
	//		fmt.Println("根据群id查询群的详细信息error:", err)
	//		return err
	//	}
	//
	//	group.UnreadMessage = unreadmsgmap[g.GroupID] //放入未读消息数量
	//
	//	// 将该群聊的消息放入消息列表
	//	for _, m := range usermessagelist {
	//		// fmt.Printf("%+v-----%+v\n", m.GroupID, g.ID)
	//		if len(messagelist) >= 10 {
	//			break
	//		}
	//		if m.GroupID == g.GroupID {
	//			messagelist = append(messagelist, m)
	//		}
	//	}
	//	sort.Slice(messagelist, func(i, j int) bool { return messagelist[i].ID < (messagelist[j].ID) })
	//	//MessageTimeSort(messagelist, "asc")
	//
	//	var groupitem = GroupDetail{
	//		GroupInfo:   group,
	//		MessageList: messagelist,
	//	}
	//	*grouplist = append(*grouplist, groupitem)
	//}
}

// GetApplyMsgList 获取用户的群聊通知列表
func (u *Users) GetApplyMsgList(applylist *[]ApplyJoinGroupResponse) error {
	var grouplist []Group
	if err := adb.Ssql.Table("group").Where("creater_id=?", u.ID).Find(&grouplist); err != nil {
		return err
	}
	//获取群主为此用户的群聊的通知
	for _, g := range grouplist {
		var applyjoingrouplist []ApplyJoinGroupResponse
		if err := adb.Ssql.Table("apply_join_group").Cols("apply_join_group.*", "group.group_name").
			Where("group_id=?", g.ID).Join("INNER", "group", "group.id=apply_join_group.group_id").Find(&applyjoingrouplist); err != nil {
			return err
		}
		*applylist = append(*applylist, applyjoingrouplist...)
	}
	//获取此用户发起申请的群聊通知
	userapplylist := make([]ApplyJoinGroupResponse, 0)
	if err := adb.Ssql.Table("apply_join_group").Cols("apply_join_group.*", "group.group_name").
		Where("apply_user_id =?", u.ID).Join("INNER", "group", "group.id=apply_join_group.group_id").
		Find(&userapplylist); err != nil {
		return err
	}
	*applylist = append(*applylist, userapplylist...)

	return nil
}

// GetApplyAddUserList 获取用户的好友申请列表
func (u *Users) GetApplyAddUserList(applylist *[]ApplyAddUser) error {
	var userlist []ApplyAddUser
	if err := adb.Ssql.Table("apply_add_user").Where("pre_apply_user_id=? or apply_user_id=?", u.ID, u.ID).Desc(" created_at").Find(&userlist); err != nil {
		return err
	}
	*applylist = userlist
	return nil
}

// GetFriendList 获取好友列表
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
	//fmt.Printf("%+v\n", unreadmessagelist)

	// 映射map,方便操作
	unreadmessagemap := make(map[int]int, len(unreadmessagelist)+1) //key:friend_id  v:unread_number
	for _, m := range unreadmessagelist {
		unreadmessagemap[m.UserID] = m.UnreadNumber
	}
	//fmt.Printf("%+v\n", unreadmessagemap)
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
	applygrouplist := make([]ApplyJoinGroupResponse, 0)
	err := u.GetApplyMsgList(&applygrouplist)
	if err != nil {
		return err
	}
	sort.Slice(applygrouplist, func(i, j int) bool {
		return applygrouplist[i].CreatedAt.After(applygrouplist[j].CreatedAt)
	})

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
	} //根据未读消息排序
	sort.Slice(friendlist, func(i, j int) bool {
		return friendlist[i].UnreadMessage > friendlist[j].UnreadMessage
	})

	// 群聊列表
	grouplist := make([]GroupDetail, 0)
	err = u.GetUserGroupList(&grouplist)
	if err != nil {
		return err
	} //根据未读消息排序
	sort.Slice(grouplist, func(i, j int) bool {
		return grouplist[i].GroupInfo.UnreadMessage > friendlist[j].UnreadMessage
	})

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
		Introduce:     userdata.Introduce,
		GroupList:     grouplist,
		ApplyList:     applygrouplist,
		ApplyUserList: applyuserlist,
		FriendList:    friendlist,
	}

	return nil
}
