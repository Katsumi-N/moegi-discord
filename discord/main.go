package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"grpc-conoha/config"
	conohapb "grpc-conoha/pkg/grpc"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// discord bot
	discordToken := config.Config.DiscordToken
	dg, err := discordgo.New("Bot " + discordToken)
	defer dg.Close()
	// dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages)
	dg.AddHandler(Minecraft)
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection,", err)
		return
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

func Minecraft(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		Command: m.Content,
	}
	res, err := client.Minecraft(context.Background(), req)
	if err != nil {
		log.Fatal(err)
		return
	}
	testMessage := "@here gRPCこのはちゃんどりゃあああああああああああああ"
	s.ChannelMessageSend(m.ChannelID, testMessage)
	fmt.Println(res.GetIsNormal(), res.GetMessage())
}
