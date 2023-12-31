package rpcserver

import (
	"fmt"
	pb "go-websocket-server/proto"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// var mu sync.Mutex // 创建一个全局的互斥锁

type GetUserGroupListrpcserver struct{}

func (s *GetUserGroupListrpcserver) GetUserGroupList(ctx context.Context, in *pb.UserData) (*pb.UserGroupList, error) {
	// 处理逻辑后返回
	grouplist, err := GetUserGroupList(int(in.Userid))
	if err != nil {
		log.Println(err.Error())
		fmt.Printf("获取群聊列表数据失败!%v", err)
		return &pb.UserGroupList{}, err
	}
	applyjoingrouplist, err := GetUserApplyJoinGroupList(int(in.Userid))
	if err != nil {
		log.Println(err.Error())
		fmt.Printf("获取群聊通知失败!%v", err)
		return &pb.UserGroupList{}, err
	}

	applyadduserlist, err := GetUserApplyAddUserList(int(in.Userid))
	if err != nil {
		log.Println(err.Error())
		fmt.Printf("获取好友通知失败!%v", err)
		return &pb.UserGroupList{}, err
	}

	friendlist, err := GetFriendList(int(in.Userid))
	if err != nil {
		log.Println(err.Error())
		fmt.Printf("获取好友列表失败!%v", err.Error())
		return &pb.UserGroupList{}, err
	}

	return &pb.UserGroupList{
		GroupDetail:        grouplist,
		ApplyJoinGroupList: applyjoingrouplist,
		ApplyAddUserList:   applyadduserlist,
		FriendList:         friendlist,
	}, nil
}

func ListenGetUserGroupListRpcServer() {
	// 监听127.0.0.1:50051地址
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("GetUserGroupListServer listening on: %s\n", lis.Addr().String())

	// 实例化grpc服务端
	s := grpc.NewServer()
	pb.RegisterLoginServer(s, &GetUserGroupListrpcserver{})

	// 往grpc服务端注册反射服务
	reflection.Register(s)

	// 启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
