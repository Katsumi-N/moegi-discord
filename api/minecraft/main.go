package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	conohapb "grpc-conoha/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type conohaServer struct {
	conohapb.UnimplementedConohaServiceServer
}

// サービスメソッドのサンプル
func (s *conohaServer) Minecraft(ctx context.Context, req *conohapb.MinecraftRequest) (*conohapb.MinecraftResponse, error) {
	return &conohapb.MinecraftResponse{
		Message: fmt.Sprintf("Hello, command is %s", req.GetCommand()),
	}, nil
}

// 自作サービス構造体のコンストラクタ
func NewConohaServer() *conohaServer {
	return &conohaServer{}
}

func main() {
	port := 8080
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーを作成
	s := grpc.NewServer()

	// gRPCサーバーにserviceを登録
	conohapb.RegisterConohaServiceServer(s, NewConohaServer())

	// grpcURL用にサーバーリフレクションを設定する
	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port->%v", port)
		s.Serve(listner)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server")
	s.GracefulStop()
}
