package rpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hichat-streammedia-service/config"
	"hichat-streammedia-service/models"
	pb "hichat-streammedia-service/proto"
	"hichat-streammedia-service/util"
	"time"

	"log"
	"net"
)

type noticevideostreamrpcserver struct{}

func (n *noticevideostreamrpcserver) Noticevideostreamserver(ctx context.Context, in *pb.Noticevideostreamserverreq) (*pb.Noticevideostreamserverres, error) {
	roomuuid := util.GenerateUUID()
	room := models.UserToUserRoom{
		RoomUUID:         roomuuid,
		RoomName:         "",
		RoomType:         1,
		StartUserID:      int(in.StartUserid),
		StartUserName:    in.StartUsername,
		ReceivedUserID:   int(in.ReceiveUserid),
		ReceivedUserName: in.ReceiveUsername,
		CreateTime:       time.Now(),
	}
	models.ServiceCenter.Room[roomuuid] = room
	fmt.Println("Room created")
	go room.CheckUserLive()

	fmt.Println(models.ServiceCenter.Room[roomuuid].StartUserID, models.ServiceCenter.Room[roomuuid].ReceivedUserID)

	response := &pb.Noticevideostreamserverres{
		Status: 200,
		Msg:    "创建成功",
	}
	return response, nil
}

func ListenNoticeVideoStreamRpcServer() {
	// 监听127.0.0.1:50052地址
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.RPCServerIP, config.RPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("ListenNoticeVideoStreamRpcServer listening on: %s\n", lis.Addr().String())

	// 实例化grpc服务端
	s := grpc.NewServer()
	pb.RegisterVideostreamserviceServer(s, &noticevideostreamrpcserver{})
	// 往grpc服务端注册反射服务
	reflection.Register(s)

	// 启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
