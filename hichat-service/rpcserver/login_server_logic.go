package rpcserver

import (
	"fmt"
	adb "go-websocket-server/ADB"
	"go-websocket-server/models"
	pb "go-websocket-server/proto"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var mu sync.Mutex // 创建一个全局的互斥锁
func GetUserGroupList(ID int) ([]*pb.GroupDetail, error) {
	mu.Lock()         // 加锁
	defer mu.Unlock() // 在函数结束时解锁，确保释放资源
	session := adb.Ssql.NewSession()
	var usergouplist []*pb.GroupDetail

	// 查询用户加入的群列表(没有详情)
	var gur []models.GroupUserRelative
	if err := adb.Ssql.Table("group_user_relative").Where("user_id=?", ID).Find(&gur); err != nil {
		fmt.Println("查询用户加入的群列表error:", err)
		session.Rollback()
		return []*pb.GroupDetail{}, err
	}

	// 查询用户的所有消息
	var usermessagelist []models.GroupMessage
	if err := adb.Ssql.Table("group_message").Find(&usermessagelist); err != nil {
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

	for _, g := range gur {
		var group models.Group             //群详情
		var messagelist []*pb.GroupMessage //群消息列表

		// 查询该用户加入的每个群聊的人数
		membercount, err := adb.Ssql.Table("group_user_relative").Where("group_id=?", g.GroupID).Count()
		if err != nil {
			fmt.Println("查询用户加入的每个群聊的人数error:", err)
			session.Rollback()
			return []*pb.GroupDetail{}, err
		}

		//  根据群id查询群的详细信息
		_, err = adb.Ssql.Table("group").Where("uuid=?", g.GroupUUID).Get(&group)
		if err != nil {
			fmt.Println("根据群id查询群的详细信息error:", err)
			session.Rollback()
			return []*pb.GroupDetail{}, err
		}
		pbgroup := pb.Group{
			Id:          int32(group.ID),
			Uuid:        group.UUID,
			CreaterId:   int32(group.CreaterID),
			CreaterName: group.CreaterName,
			GroupName:   group.GroupName,
			Avatar:      group.Avatar,
			Grade:       int32(group.Grade),
			CreatedAt:   timestamppb.New(group.CreatedAt),
			DeletedAt:   timestamppb.New(group.DeletedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
			MemberCount: int32(membercount),
		}

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

// func GetUserGroupList(uid int) ([]models.GroupDetail, error) {
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	if err != nil {
// 		return []models.GroupDetail{}, err
// 	}
// 	defer conn.Close()

// 	c := pb.NewLoginClient(conn) //初始化客户端

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // 初始化上下文，设置请求超时时间为1秒
// 	defer cancel()

// 	r, err := c.GetUserGroupList(ctx, &pb.UserData{Userid: int32(uid)})
// 	if err != nil {
// 		return []models.GroupDetail{}, err
// 	}

// 	var data []models.GroupDetail
// 	for _, g := range r.GroupDetail {

// 		var messageList []models.GroupMessage
// 		for _, m := range g.MessageList {
// 			messageList = append(messageList, models.GroupMessage{
// 				ID:          int(m.Id),
// 				UserID:      int(m.UserId),
// 				UserUUID:    m.UserUuid,
// 				UserName:    m.UserName,
// 				GroupID:     int(m.GroupId),
// 				Msg:         m.Msg,
// 				MsgType:     int(m.MsgType),
// 				IsReply:     m.IsReply,          //是否是回复消息
// 				ReplyUserID: int(m.ReplyUserId), //如果是,被回复的用户id
// 				Context:     m.Context,
// 				CreatedAt:   m.CreatedAt.AsTime(),
// 				DeletedAt:   m.DeletedAt.AsTime(),
// 				UpdatedAt:   m.UpdatedAt.AsTime(),
// 			})
// 		}
// 		data = append(data, models.GroupDetail{
// 			GroupInfo: models.Group{
// 				ID:          int(g.GroupInfo.Id),
// 				UUID:        g.GroupInfo.Uuid,
// 				CreaterID:   int(g.GroupInfo.CreaterId),
// 				CreaterName: g.GroupInfo.CreaterName,
// 				GroupName:   g.GroupInfo.GroupName,
// 				Avatar:      g.GroupInfo.Avatar,
// 				Grade:       int(g.GroupInfo.Grade),
// 				CreatedAt:   g.GroupInfo.CreatedAt.AsTime(),
// 				DeletedAt:   g.GroupInfo.DeletedAt.AsTime(),
// 				UpdatedAt:   g.GroupInfo.UpdatedAt.AsTime(),
// 			},
// 			MessageList: messageList,
// 		})
// 	}
// 	// copier.CopyWithOption(&data, &r.GroupDetail, copier.Option{DeepCopy: true})
// 	fmt.Printf("%+v\n", data)

// 	return data, nil
// }
