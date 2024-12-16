package rpcserver

import (
	"context"
	"fmt"
	"go-websocket-server/config"
	"go-websocket-server/models"
	pb "go-websocket-server/proto"
	"google.golang.org/grpc"
	"log"
)

func CallNoticeVideoStreamServer(info models.UserMessage) (*pb.Noticevideostreamserverres, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.CallNoticeVideoStreamServerIP, config.CallNoticeVideoStreamServerPort), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		log.Printf("did not connect: %v", err)
		fmt.Printf("CallNoticeVideoStreamServer did not connect: %v", err)
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	// 创建 gRPC 客户端
	client := pb.NewVideostreamserviceClient(conn)

	// 构造请求
	req := &pb.Noticevideostreamserverreq{
		StartUserid:     int32(info.UserID),
		StartUsername:   info.UserName,
		ReceiveUserid:   int32(info.ReceiveUserID),
		ReceiveUsername: info.ReceiveUserName,
	}

	// 调用 gRPC 服务
	res, err := client.Noticevideostreamserver(context.Background(), req)
	if err != nil {
		log.Printf("调用 音视频RPC失败: %v", err)
		return nil, err
	}

	// 返回结果
	return res, nil

}
