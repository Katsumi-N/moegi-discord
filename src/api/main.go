package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"grpc-conoha/api/conoha"
	"grpc-conoha/config"
	conohapb "grpc-conoha/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type conohaServer struct {
	conohapb.UnimplementedConohaServiceServer
}

var statusName = map[string]string{
	"SHUTOFF":       "シャットダウンしてるよ",
	"ACTIVE":        "起動してるよ",
	"RESIZE":        "リサイズ中",
	"REBOOT":        "再起動中",
	"VERIFY_RESIZE": "リサイズ承認待ち",
}

// マインクラフト用のConoha VPSサーバーを起動/シャッタダウン/再起動する
func (s *conohaServer) Minecraft(req *conohapb.MinecraftRequest, stream conohapb.ConohaService_MinecraftServer) error {
	token := conoha.GetToken(config.Config.Username, config.Config.Password, config.Config.TenantId)

	if req.GetCommand() == "!conoha server" {
		status, _ := conoha.GetServerStatus(token)

		if err := stream.Send(&conohapb.MinecraftResponse{
			Message:  statusName[string(status)],
			IsNormal: true,
		}); err != nil {
			return err
		}
		return nil
	}
	if req.GetCommand() == "!conoha start" {
		_, statusCode := conoha.StartServer(token, stream)
		is_normal := true
		if statusCode != 202 {
			is_normal = false
		}
		if err := stream.Send(&conohapb.MinecraftResponse{
			Message:  "サーバーを起動しました",
			IsNormal: is_normal,
		}); err != nil {
			return err
		}
		return nil
	}
	if req.GetCommand() == "!conoha stop" {
		_, statusCode := conoha.StopServer(token, stream)
		is_normal := true
		if statusCode != 202 {
			is_normal = false
		}
		if err := stream.Send(&conohapb.MinecraftResponse{
			Message:  "サーバーをシャットダウンしました",
			IsNormal: is_normal,
		}); err != nil {
			return err
		}
	}
	if req.GetCommand() == "!conoha reboot" {
		_, statusCode := conoha.RebootServer(token)
		is_normal := true
		if statusCode != 202 {
			is_normal = false
		}
		if err := stream.Send(&conohapb.MinecraftResponse{
			Message:  "サーバーを再起動しました．",
			IsNormal: is_normal,
		}); err != nil {
			return err
		}
	}
	grpcerr := status.Error(codes.Unimplemented, "登録されていないコマンドです")
	return grpcerr
}

func (s *conohaServer) Vote(ctx context.Context, req *conohapb.VoteRequest) (*conohapb.VoteResponse, error) {
	// いいかんじに投票メッセージを作る
	// Optionsの個数分投票しておく
	retMes := "投票タイトル！: " + req.GetTitle() + "選択肢の個数: " + strconv.Itoa(len(req.GetOptions()))
	return &conohapb.VoteResponse{
		Message: retMes,
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
