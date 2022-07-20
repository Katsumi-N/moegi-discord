package main

import (
	conohapb "grpc-conoha/pkg/grpc"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client conohapb.ConohaServiceClient
)

func main() {

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed")
		return
	}
	defer conn.Close()

	// gRPCクライアントを生成
	client = conohapb.NewConohaServiceClient(conn)

}
