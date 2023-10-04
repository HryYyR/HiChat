package rpcserver

import (
	"context"
	"hichat_static_server/models"
	pb "hichat_static_server/proto"
	"time"

	"google.golang.org/grpc"
)

func GetUserGroupList(uid int) ([]models.GroupDetail, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return []models.GroupDetail{}, err
	}
	defer conn.Close()

	c := pb.NewLoginClient(conn) //初始化客户端

	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // 初始化上下文，设置请求超时时间为1秒
	defer cancel()

	r, err := c.GetUserGroupList(ctx, &pb.UserData{Userid: int32(uid)})
	if err != nil {
		return []models.GroupDetail{}, err
	}

	var data []models.GroupDetail
	for _, g := range r.GroupDetail {
		var messageList []models.GroupMessage
		for _, m := range g.MessageList {
			messageList = append(messageList, models.GroupMessage{
				ID:          int(m.Id),
				UserID:      int(m.UserId),
				UserUUID:    m.UserUuid,
				UserName:    m.UserName,
				GroupID:     int(m.GroupId),
				Msg:         m.Msg,
				MsgType:     int(m.MsgType),
				IsReply:     m.IsReply,          //是否是回复消息
				ReplyUserID: int(m.ReplyUserId), //如果是,被回复的用户id
				Context:     m.Context,
				CreatedAt:   m.CreatedAt.AsTime(),
				DeletedAt:   m.DeletedAt.AsTime(),
				UpdatedAt:   m.UpdatedAt.AsTime(),
			})
		}
		data = append(data, models.GroupDetail{
			GroupInfo: models.Group{
				ID:          int(g.GroupInfo.Id),
				UUID:        g.GroupInfo.Uuid,
				CreaterID:   int(g.GroupInfo.CreaterId),
				CreaterName: g.GroupInfo.CreaterName,
				GroupName:   g.GroupInfo.GroupName,
				Avatar:      g.GroupInfo.Avatar,
				Grade:       int(g.GroupInfo.Grade),
				MemberCount: int(g.GroupInfo.MemberCount),
				CreatedAt:   g.GroupInfo.CreatedAt.AsTime(),
				DeletedAt:   g.GroupInfo.DeletedAt.AsTime(),
				UpdatedAt:   g.GroupInfo.UpdatedAt.AsTime(),
			},
			MessageList: messageList,
		})
	}
	// copier.CopyWithOption(&data, &r.GroupDetail, copier.Option{DeepCopy: true})
	// fmt.Printf("%+v\n", data)

	return data, nil
}
