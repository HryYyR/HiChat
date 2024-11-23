package service_registry

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hichat_static_server/proto"
	"hichat_static_server/rpcserver"
	"log"
	"net"
	"net/http"
)

type LoginRegistryServiceConfig struct {
	RpcAddr  string
	HttpAddr string
}

func LoginRegistryService(conf LoginRegistryServiceConfig) {
	lis, err := net.Listen("tcp", conf.RpcAddr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer()                             // 创建一个gRPC server对象
	proto.RegisterLoginServer(s, &rpcserver.Server{}) //注册
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.DialContext( // 创建一个连接到我们刚刚启动的 gRPC 服务器的客户端连接
		context.Background(), // gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
		"0.0.0.0"+conf.RpcAddr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// 注册Greeter
	err = proto.RegisterLoginHandler(context.Background(), gwmux, conn)
	if err != nil {
		fmt.Println("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    conf.HttpAddr,
		Handler: gwmux,
	}
	// 8090端口提供gRPC-Gateway服务
	err = gwServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
