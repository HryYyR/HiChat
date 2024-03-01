package rpcserver

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	pb "go-websocket-server/proto"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetUserGroupList var mu sync.Mutex // 创建一个全局的互斥锁
func GetUserGroupList(ID int) ([]*pb.GroupDetail, error) {
	session := adb.Ssql.NewSession()
	defer session.Close()
	session.Begin()
	var usergouplist []*pb.GroupDetail

	// 查询用户加入的群列表(没有详情)
	var gur []models.GroupUserRelative
	if err := session.Table("group_user_relative").Where("user_id=?", ID).Find(&gur); err != nil {
		fmt.Println("查询用户加入的群列表error:", err)
		session.Rollback()
		return []*pb.GroupDetail{}, err
	}

	// 查询用户的所有消息
	var usermessagelist []models.GroupMessage
	if err := session.Table("group_message").Find(&usermessagelist); err != nil {
		fmt.Println("查询所有消息error:", err)
		session.Rollback()
		return []*pb.GroupDetail{}, err
	}

	var pbusermessagelist []*pb.GroupMessage
	for _, m := range usermessagelist {
		pbusermessagelist = append(pbusermessagelist, &pb.GroupMessage{
			Id:          int32(m.ID),
			UserId:      int32(m.UserID),
			UserUuid:    m.UserUUID,
			UserName:    m.UserName,
			UserAvatar:  m.UserAvatar,
			UserCity:    m.UserCity,
			UserAge:     strconv.Itoa(m.UserAge),
			GroupId:     int32(m.GroupID),
			Msg:         m.Msg,
			MsgType:     int32(m.MsgType),
			IsReply:     m.IsReply,
			ReplyUserId: int32(m.ReplyUserID),
			Context:     m.Context,
			CreatedAt:   timestamppb.New(m.CreatedAt),
			DeletedAt:   timestamppb.New(m.DeletedAt),
			UpdatedAt:   timestamppb.New(m.UpdatedAt),
		})
	}

	// 查询用户的未读消息
	var rawunreadmsglist []models.GroupUnreadMessage
	if err := session.Table("group_unread_message").Where("user_id = ?", ID).Find(&rawunreadmsglist); err != nil {
		fmt.Println("查询未读消息error:", err)
		session.Rollback()
		return []*pb.GroupDetail{}, err
	}
	unreadmsglist := make(map[int]int, 0)
	for _, UnreadMessage := range rawunreadmsglist {
		unreadmsglist[UnreadMessage.GroupID] = UnreadMessage.UnreadNumber
	}

	for _, g := range gur {
		// 查询该用户加入的每个群聊的人数
		membercount, err := session.Table("group_user_relative").Where("group_id=?", g.GroupID).Count()
		if err != nil {
			fmt.Println("查询用户加入的每个群聊的人数error:", err)
			session.Rollback()
			return []*pb.GroupDetail{}, err
		}

		var unreadmsg int
		if num, has := unreadmsglist[g.GroupID]; has {
			unreadmsg = num
		}

		var group models.Group //群详情
		//  根据群id查询群的详细信息
		_, err = session.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
		if err != nil {
			fmt.Println("根据群id查询群的详细信息error:", err)
			session.Rollback()
			return []*pb.GroupDetail{}, err
		}
		pbgroup := pb.Group{
			Id:            int32(group.ID),
			Uuid:          group.UUID,
			CreaterId:     int32(group.CreaterID),
			CreaterName:   group.CreaterName,
			GroupName:     group.GroupName,
			Avatar:        group.Avatar,
			Grade:         int32(group.Grade),
			MemberCount:   int32(membercount),
			UnreadMessage: int32(unreadmsg),
			CreatedAt:     timestamppb.New(group.CreatedAt),
			DeletedAt:     timestamppb.New(group.DeletedAt),
			UpdatedAt:     timestamppb.New(group.UpdatedAt),
		}

		var messagelist []*pb.GroupMessage //群消息列表
		// 将该群聊的消息放入消息列表
		for _, m := range pbusermessagelist {
			// fmt.Printf("%+v-----%+v\n", m.GroupID, g.ID)
			if int(m.GroupId) == g.GroupID {
				messagelist = append(messagelist, m)
			}
		}

		var groupitem = pb.GroupDetail{
			GroupInfo:   &pbgroup,
			MessageList: messagelist,
		}
		usergouplist = append(usergouplist, &groupitem)
	}
	session.Commit()
	return usergouplist, nil
}

// GetUserApplyJoinGroupList 获取该用户的通知列表
func GetUserApplyJoinGroupList(ID int) ([]*pb.ApplyJoinGroupMessage, error) {
	// 该用户创建的群聊列表
	var usercreategrouplist []models.Group
	if err := adb.Ssql.Table("group").Where("creater_id=?", ID).Find(&usercreategrouplist); err != nil {
		return nil, err
	}
	userapplylist := make([]*pb.ApplyJoinGroupMessage, 0)
	for _, group := range usercreategrouplist {
		var applyuserlist []models.ApplyJoinGroup
		adb.Ssql.Table("apply_join_group").Where("group_id = ?", group.ID).Find(&applyuserlist)
		for _, applyuser := range applyuserlist {
			userapplylist = append(userapplylist, &pb.ApplyJoinGroupMessage{
				Id:           int32(applyuser.ID),
				ApplyUserId:  int32(applyuser.ApplyUserID),
				AplyUserName: applyuser.ApplyUserName,
				GroupId:      int32(applyuser.GroupID),
				ApplyMsg:     applyuser.ApplyMsg,
				ApplyWay:     int32(applyuser.ApplyWay),
				HandleStatus: int32(applyuser.HandleStatus),
				CreatedAt:    timestamppb.New(applyuser.CreatedAt),
				DeletedAt:    timestamppb.New(applyuser.DeletedAt),
				UpdatedAt:    timestamppb.New(applyuser.UpdatedAt),
			})
		}
	}

	return userapplylist, nil
}

// GetUserApplyAddUserList 获取该用户的好友通知列表
func GetUserApplyAddUserList(ID int) ([]*pb.ApplyAddUserMessage, error) {

	var userapplyadduserlist []models.ApplyAddUser
	if err := adb.Ssql.Table("apply_add_user").Where("pre_apply_user_id=? or apply_user_id=?", ID, ID).Find(&userapplyadduserlist); err != nil {
		return []*pb.ApplyAddUserMessage{}, err
	}

	adduserlist := make([]*pb.ApplyAddUserMessage, 0)
	for _, applyuser := range userapplyadduserlist {
		adduserlist = append(adduserlist, &pb.ApplyAddUserMessage{
			Id:               int32(applyuser.ID),
			PreApplyUserId:   int32(applyuser.PreApplyUserID),
			PreApplyUserName: applyuser.PreApplyUserName,
			ApplyUserId:      int32(applyuser.ApplyUserID),
			ApplyUserName:    applyuser.ApplyUserName,
			ApplyMsg:         applyuser.ApplyMsg,
			ApplyWay:         int32(applyuser.ApplyWay),
			HandleStatus:     int32(applyuser.HandleStatus),
			CreatedAt:        timestamppb.New(applyuser.CreatedAt),
			DeletedAt:        timestamppb.New(applyuser.DeletedAt),
			UpdatedAt:        timestamppb.New(applyuser.UpdatedAt),
		})
	}
	return adduserlist, nil
}

func GetFriendList(ID int) ([]*pb.FriendList, error) {
	var friendrelativelist []models.UserUserRelative
	err := adb.Ssql.Table("user_user_relative").Where("pre_user_id = ? or back_user_id=?", ID, ID).Find(&friendrelativelist)
	if err != nil {
		fmt.Println("获取用户关系失败")
		return []*pb.FriendList{}, err
	}

	var friendlist []*pb.FriendList

	for _, relative := range friendrelativelist {
		var frienddata models.Users
		if relative.BackUserID == ID {
			exit, err := adb.Ssql.Table("users").Where("id = ?", relative.PreUserID).Get(&frienddata)
			if err != nil {
				fmt.Println("获取用户信息失败")
				return []*pb.FriendList{}, err
			}
			if !exit {
				continue
			}
		} else {
			exit, err := adb.Ssql.Table("users").Where("id = ?", relative.BackUserID).Get(&frienddata)
			if err != nil {
				fmt.Println("获取用户信息失败")
				return []*pb.FriendList{}, err
			}
			if !exit {
				continue
			}
		}
		friendlist = append(friendlist, &pb.FriendList{
			Id:        int32(frienddata.ID),
			UserName:  frienddata.UserName,
			NikeName:  frienddata.NikeName,
			Email:     frienddata.Email,
			Avatar:    frienddata.Avatar,
			City:      frienddata.City,
			Age:       strconv.Itoa(frienddata.Age),
			CreatedAt: timestamppb.New(frienddata.CreatedAt),
			DeletedAt: timestamppb.New(frienddata.DeletedAt),
			UpdatedAt: timestamppb.New(frienddata.UpdatedAt),
		})
	}

	return friendlist, nil
}
