package rpcserver

import (
	"context"
	"hichat_static_server/models"
	pb "hichat_static_server/proto"
	"hichat_static_server/util"
	"time"

	"google.golang.org/grpc"
)

func GetUserGroupList(uid int) (models.UserGroupList, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return models.UserGroupList{}, err
	}
	defer conn.Close()

	c := pb.NewLoginClient(conn) //初始化客户端

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3) // 初始化上下文，设置请求超时时间为3秒
	defer cancel()

	r, err := c.GetUserGroupList(ctx, &pb.UserData{Userid: int32(uid)})
	if err != nil {
		return models.UserGroupList{}, err
	}
	// 处理群相关数据
	groupdata := make([]models.GroupDetail, 0)
	for _, g := range r.GroupDetail {
		var messageList []models.GroupMessage
		for _, m := range g.MessageList {
			messageList = append(messageList, models.GroupMessage{
				ID:          int(m.Id),
				UserID:      int(m.UserId),
				UserUUID:    m.UserUuid,
				UserName:    m.UserName,
				UserAvatar:  m.UserAvatar,
				UserCity:    m.UserCity,
				UserAge:     m.UserAge,
				GroupID:     int(m.GroupId),
				Msg:         m.Msg,
				MsgType:     int(m.MsgType),
				IsReply:     m.IsReply,          //是否是回复消息
				ReplyUserID: int(m.ReplyUserId), //如果是,被回复的用户id
				Context:     m.Context,
				CreatedAt:   util.FormatTampTime(m.CreatedAt),
				DeletedAt:   util.FormatTampTime(m.DeletedAt),
				UpdatedAt:   util.FormatTampTime(m.UpdatedAt),
			})
		}
		groupdata = append(groupdata, models.GroupDetail{
			GroupInfo: models.Group{
				ID:            int(g.GroupInfo.Id),
				UUID:          g.GroupInfo.Uuid,
				CreaterID:     int(g.GroupInfo.CreaterId),
				CreaterName:   g.GroupInfo.CreaterName,
				GroupName:     g.GroupInfo.GroupName,
				Avatar:        g.GroupInfo.Avatar,
				Grade:         int(g.GroupInfo.Grade),
				UnreadMessage: int(g.GroupInfo.UnreadMessage),
				MemberCount:   int(g.GroupInfo.MemberCount),
				CreatedAt:     util.FormatTampTime(g.GroupInfo.CreatedAt),
				DeletedAt:     util.FormatTampTime(g.GroupInfo.DeletedAt),
				UpdatedAt:     util.FormatTampTime(g.GroupInfo.UpdatedAt),
			},
			MessageList: messageList,
		})
	}

	// 处理群聊申请消息相关数据
	applydata := make([]models.ApplyJoinGroup, 0)
	for _, applyuser := range r.ApplyJoinGroupList {
		applydata = append(applydata, models.ApplyJoinGroup{
			ID:            int(applyuser.Id),
			ApplyUserID:   int(applyuser.ApplyUserId),
			ApplyUserName: applyuser.AplyUserName,
			GroupID:       int(applyuser.GroupId),
			ApplyMsg:      applyuser.ApplyMsg,
			ApplyWay:      int(applyuser.ApplyWay),
			HandleStatus:  int(applyuser.HandleStatus),
			CreatedAt:     util.FormatTampTime(applyuser.CreatedAt),
			DeletedAt:     util.FormatTampTime(applyuser.DeletedAt),
			UpdatedAt:     util.FormatTampTime(applyuser.UpdatedAt),
		})
	}
	util.TimeSortApplyJoinGroupList(applydata, "desc")

	// 处理好友申请消息相关数据
	applyadduserdata := make([]models.ApplyAddUser, 0)
	for _, applyuser := range r.ApplyAddUserList {
		applyadduserdata = append(applyadduserdata, models.ApplyAddUser{
			ID:               int(applyuser.Id),
			PreApplyUserID:   int(applyuser.PreApplyUserId),
			PreApplyUserName: applyuser.PreApplyUserName,
			ApplyUserID:      int(applyuser.ApplyUserId),
			ApplyUserName:    applyuser.ApplyUserName,
			ApplyMsg:         applyuser.ApplyMsg,
			ApplyWay:         applyuser.ApplyWay,
			HandleStatus:     int(applyuser.HandleStatus),
			CreatedAt:        util.FormatTampTime(applyuser.CreatedAt),
			DeletedAt:        util.FormatTampTime(applyuser.DeletedAt),
			UpdatedAt:        util.FormatTampTime(applyuser.UpdatedAt),
		})
	}
	util.TimeSortAddUserList(applyadduserdata, "desc")

	// 处理好友列表相关数据
	friendlist := make([]models.Friend, 0)
	for _, Friend := range r.FriendList {
		friendlist = append(friendlist, models.Friend{
			Id:        Friend.Id,
			UserName:  Friend.UserName,
			NikeName:  Friend.NikeName,
			Email:     Friend.Email,
			Avatar:    Friend.Avatar,
			City:      Friend.City,
			Age:       Friend.Age,
			CreatedAt: util.FormatTampTime(Friend.CreatedAt),
			DeletedAt: util.FormatTampTime(Friend.DeletedAt),
			UpdatedAt: util.FormatTampTime(Friend.UpdatedAt),
		})
	}

	UserGroupList := models.UserGroupList{
		GroupDetail:           groupdata,
		ApplyJoinGroupMessage: applydata,
		ApplyAddUserMessage:   applyadduserdata,
		FriendList:            friendlist,
	}
	// fmt.Printf("%+v\n", UserGroupList)

	return UserGroupList, nil
}
