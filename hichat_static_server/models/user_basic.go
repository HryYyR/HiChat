package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	adb "hichat_static_server/ADB"
	"hichat_static_server/proto"
	"hichat_static_server/tool"
	"sort"
	"strconv"
	"sync"
	"time"
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

func (r *ResponseUserData) ResponseUserDataToProto() *proto.ResponseUserData {
	GroupList := make([]*proto.GroupList, 0)
	for _, d := range r.GroupList {
		ginfo := proto.Group{
			Id:            int32(d.GroupInfo.ID),
			Uuid:          d.GroupInfo.UUID,
			CreaterId:     int32(d.GroupInfo.CreaterID),
			CreaterName:   d.GroupInfo.CreaterName,
			GroupName:     d.GroupInfo.GroupName,
			Avatar:        d.GroupInfo.Avatar,
			Grade:         int32(d.GroupInfo.Grade),
			MemberCount:   int32(d.GroupInfo.MemberCount),
			UnreadMessage: int32(d.GroupInfo.UnreadMessage),
			CreatedAt:     timestamppb.New(d.GroupInfo.CreatedAt),
			DeletedAt:     timestamppb.New(d.GroupInfo.DeletedAt),
			UpdatedAt:     timestamppb.New(d.GroupInfo.UpdatedAt),
		}
		mlist := make([]*proto.GroupMessage, 0)
		for _, m := range d.MessageList {
			gmsg := proto.GroupMessage{
				Id:          int32(m.ID),
				UserId:      int32(m.UserID),
				UserUuid:    m.UserUUID,
				UserName:    m.UserName,
				UserAvatar:  m.UserAvatar,
				UserCity:    m.UserCity,
				UserAge:     m.UserAge,
				GroupId:     int32(m.GroupID),
				Msg:         m.Msg,
				MsgType:     int32(m.MsgType),
				IsReply:     m.IsReply,
				ReplyUserId: int32(m.ReplyUserID),
				Context:     m.Context,
				CreatedAt:   timestamppb.New(m.CreatedAt),
				DeletedAt:   timestamppb.New(m.DeletedAt),
				UpdatedAt:   timestamppb.New(m.UpdatedAt),
			}
			mlist = append(mlist, &gmsg)
		}
		GroupList = append(GroupList, &proto.GroupList{
			GroupInfo:   &ginfo,
			MessageList: mlist,
		})
	}

	ApplyList := make([]*proto.ApplyGroupList, 0)
	for _, r := range r.ApplyList {
		ApplyList = append(ApplyList, &proto.ApplyGroupList{
			Id:            int32(r.ID),
			ApplyUserId:   int32(r.ApplyUserID),
			ApplyUserName: r.ApplyUserName,
			GroupId:       int32(r.GroupID),
			GroupName:     r.GroupName,
			ApplyMsg:      r.ApplyMsg,
			ApplyWay:      int32(r.ApplyWay),
			HandleStatus:  int32(r.HandleStatus),
			CreatedAt:     timestamppb.New(r.CreatedAt),
			DeletedAt:     timestamppb.New(r.DeletedAt),
			UpdatedAt:     timestamppb.New(r.UpdatedAt),
		})
	}

	ApplyUserList := make([]*proto.ApplyUserList, 0)
	for _, user := range r.ApplyUserList {
		ApplyUserList = append(ApplyUserList, &proto.ApplyUserList{
			Id:               int32(user.ID),
			PreApplyUserId:   int32(user.PreApplyUserID),
			PreApplyUserName: user.PreApplyUserName,
			ApplyUserId:      int32(user.ApplyUserID),
			ApplyUserName:    user.ApplyUserName,
			ApplyMsg:         user.ApplyMsg,
			ApplyWay:         user.ApplyWay,
			HandleStatus:     int32(user.HandleStatus),
			CreatedAt:        timestamppb.New(user.CreatedAt),
			DeletedAt:        timestamppb.New(user.DeletedAt),
			UpdatedAt:        timestamppb.New(user.UpdatedAt),
		})
	}

	Friendlist := make([]*proto.FriendList, 0)
	for _, f := range r.FriendList {
		msglist := make([]*proto.UserMessage, 0)
		for _, msg := range f.MessageList {
			msglist = append(msglist, &proto.UserMessage{
				ID:                int32(msg.ID),
				UserID:            int32(msg.UserID),
				UserName:          msg.UserName,
				UserAvatar:        msg.UserAvatar,
				ReceiveUserID:     int32(msg.ReceiveUserID),
				ReceiveUserName:   msg.ReceiveUserName,
				ReceiveUserAvatar: msg.ReceiveUserAvatar,
				Msg:               msg.Msg,
				MsgType:           int32(msg.MsgType),
				IsReply:           msg.IsReply,
				ReplyUserID:       int32(msg.ReplyUserID),
				CreatedAt:         timestamppb.New(msg.CreatedAt),
				DeletedAt:         timestamppb.New(msg.DeletedAt),
				UpdatedAt:         timestamppb.New(msg.UpdatedAt),
			})
		}

		Friendlist = append(Friendlist, &proto.FriendList{
			Id:            f.Id,
			UserName:      f.UserName,
			NikeName:      f.NikeName,
			Email:         f.Email,
			Avatar:        f.Avatar,
			City:          f.City,
			Age:           f.Age,
			UnreadMessage: int32(f.UnreadMessage),
			MessageList:   msglist,
			CreatedAt:     timestamppb.New(f.CreatedAt),
			DeletedAt:     timestamppb.New(f.DeletedAt),
			UpdatedAt:     timestamppb.New(f.UpdatedAt),
		})
	}

	return &proto.ResponseUserData{
		ID:            int32(r.ID),
		UserName:      r.UserName,
		NikeName:      r.NikeName,
		Email:         r.Email,
		CreatedTime:   timestamppb.New(r.CreatedTime),
		LoginTime:     r.LoginTime,
		Avatar:        r.Avatar,
		Age:           int32(r.Age),
		City:          r.City,
		Introduce:     r.Introduce,
		GroupList:     GroupList,
		ApplyList:     ApplyList,
		ApplyUserList: ApplyUserList,
		FriendList:    Friendlist,
	}

}

func (u *Users) SaveToRedis() error {
	_, err := adb.Rediss.HMSet(strconv.Itoa(u.ID), map[string]interface{}{
		"ID":        u.ID,
		"UserName":  u.UserName,
		"NikeName":  u.NikeName,
		"Email":     u.Email,
		"CreatedAt": tool.FormatTime(u.CreatedAt),
		"LoginTime": u.LoginTime,
		"Avatar":    u.Avatar,
		"Age":       u.Age,
		"City":      u.City,
		"Introduce": u.Introduce,
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
		fmt.Println("走redis")
		_ = mapstructure.Decode(udata, &userinfo)
		userinfo.ID, _ = strconv.Atoi(udata["ID"])
		userinfo.Age, _ = strconv.Atoi(udata["Age"])
		fmt.Printf("%+v", userinfo)
		*userdata = userinfo
		return nil
	}
	fmt.Println("走mysql")

	exit, err := adb.Ssql.Omit("Password,Salt,Grade,IP").Table("users").Where("id =?", u.ID).Get(&userinfo)
	if !exit {
		return fmt.Errorf("用户不存在")
	}
	if err != nil {
		fmt.Println("mysql查询失败", err)
		return err
	}

	err = u.SaveToRedis()
	if err != nil {
		fmt.Println("保存到redis失败", err)
		//return err
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cctx, ccancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(5)

	// 群聊通知列表
	applygrouplist := make([]ApplyJoinGroupResponse, 0)
	go func(c context.Context) {
		err := u.GetApplyMsgList(&applygrouplist)
		if err != nil {
			fmt.Println("群聊通知列表", err)
			ccancel()
		}
		sort.Slice(applygrouplist, func(i, j int) bool {
			return applygrouplist[i].CreatedAt.After(applygrouplist[j].CreatedAt)
		})
		wg.Done()
	}(cctx)

	// 好友申请列表
	applyuserlist := make([]ApplyAddUser, 0)
	go func(c context.Context) {
		err := u.GetApplyAddUserList(&applyuserlist)
		if err != nil {
			fmt.Println("好友申请列表", err)
			ccancel()
		}
		wg.Done()
	}(cctx)

	// 好友列表
	friendlist := make([]FriendResponse, 0)
	go func(c context.Context) {
		err := u.GetFriendListAndMEssage(&friendlist)
		if err != nil {
			fmt.Println("好友列表", err)
			ccancel()
		}
		//根据未读消息排序
		sort.Slice(friendlist, func(i, j int) bool {
			return friendlist[i].UnreadMessage > friendlist[j].UnreadMessage
		})
		wg.Done()
	}(cctx)

	// 群聊列表
	grouplist := make([]GroupDetail, 0)
	go func(c context.Context) {
		err := u.GetUserGroupList(&grouplist)
		if err != nil {
			fmt.Println("群聊列表", err)
			ccancel()
		}
		//根据未读消息排序
		sort.Slice(grouplist, func(i, j int) bool {
			return grouplist[i].GroupInfo.UnreadMessage > grouplist[j].GroupInfo.UnreadMessage
		})
		wg.Done()
	}(cctx)

	var userdata Users
	go func(c context.Context) {
		err := u.GetUserData(&userdata)
		if err != nil {
			fmt.Println("获取用户数据error:", err)
			ccancel()
		}
		wg.Done()
	}(cctx)

	wg.Wait()

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

	var reserr error
	select {
	case <-cctx.Done():
		err := cctx.Err()
		fmt.Println(err)
		if err == context.Canceled {
			reserr = errors.New("请求失败(请求被取消)")
		} else if err == context.DeadlineExceeded {
			reserr = errors.New("请求超时")
		}
	default:
		fmt.Println("登录成功")
	}

	return reserr
}
