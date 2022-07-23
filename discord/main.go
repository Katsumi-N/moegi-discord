package main

import (
	"context"
	"fmt"
	"log"

	conohapb "grpc-conoha/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:8080"
	// gRPCサーバーとのコネクションを確立する
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("grpc cliant connection failed")
		return
	}
	defer conn.Close()

	client := conohapb.NewConohaServiceClient(conn)

	req := &conohapb.MinecraftRequest{
		Command: "!server",
	}
	res, err := client.Minecraft(context.Background(), req)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(res.GetIsNormal(), res.GetMessage())
}
