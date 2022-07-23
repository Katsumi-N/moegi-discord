package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"grpc-conoha/api/conoha"
	"grpc-conoha/config"
	conohapb "grpc-conoha/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type conohaServer struct {
	conohapb.UnimplementedConohaServiceServer
}

// サービスメソッドのサンプル
func (s *conohaServer) Minecraft(ctx context.Context, req *conohapb.MinecraftRequest) (*conohapb.MinecraftResponse, error) {
	token := conoha.GetToken(config.Config.Username, config.Config.Password, config.Config.TenantId)

	if req.GetCommand() == "!server" {
		status, _ := conoha.GetServerStatus(token)
		return &conohapb.MinecraftResponse{
			Message:  string(status),
			IsNormal: true,
		}, nil
	}
	if req.GetCommand() == "!start" {
		status, statusCode := conoha.StartServer(token)
		is_normal := true
		if statusCode != 202 {
			is_normal = false
		}
		return &conohapb.MinecraftResponse{
			Message:  string(status),
			IsNormal: is_normal,
		}, nil
	}
	if req.GetCommand() == "!stop" {
		status, statusCode := conoha.StopServer(token)
		is_normal := true
		if statusCode != 202 {
			is_normal = false
		}
		return &conohapb.MinecraftResponse{
			Message:  string(status),
			IsNormal: is_normal,
		}, nil
	}
	if req.GetCommand() == "!reboot" {
		status, statusCode := conoha.RebootServer(token)
		is_normal := true
		if statusCode != 202 {
			is_normal = false
		}
		return &conohapb.MinecraftResponse{
			Message:  string(status),
			IsNormal: is_normal,
		}, nil
	}
	return nil, nil
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
