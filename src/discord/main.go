package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
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
	dg.AddHandler(Introduction)
	dg.AddHandler(Vote)

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

func makeGrpcConnection(port string) (conn *grpc.ClientConn, err error) {
	address := "server:" + port
	// gRPCサーバーとのコネクションを確立する
	conn, err = grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("grpc cliant connection failed")
		return nil, err
	}

	return conn, nil
}

func Introduction(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := m.Content
	if command != "!intro" {
		return
	}
	s.ChannelMessageSend(m.ChannelID, "自己紹介します！")
	introMessage, err := ioutil.ReadFile("self-intro.txt")

	if err != nil {
		log.Println("can't read self-intro.txt")
		return
	}
	s.ChannelMessageSend(m.ChannelID, string(introMessage))
}

func Minecraft(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "!conoha") {
		return
	}
	conn, err := makeGrpcConnection("8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := conohapb.NewConohaServiceClient(conn)

	req := &conohapb.MinecraftRequest{
		Command: m.Content,
	}
	stream, err := client.Minecraft(context.Background(), req)
	if err != nil {
		return
	}
	for {
		res, err := stream.Recv()

		s.ChannelMessageSend(m.ChannelID, res.GetMessage())
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			break
		}
	}

	if err != nil {
		log.Print(err)
	}
}

// コマンド例 !vote 旅行先 北海道 東京 沖縄
func Vote(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "!vote") {
		return
	}
	conn, err := makeGrpcConnection("8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	voteArr := strings.Split(m.Content, " ")
	voteTitle := voteArr[1]
	if len(voteArr) <= 2 {
		log.Fatal("input vote options")
	}
	if len(voteArr) >= 12 {
		log.Fatal("vote options must be less than 10")
	}
	voteArr = voteArr[2:]
	client := conohapb.NewConohaServiceClient(conn)

	req := &conohapb.VoteRequest{
		Command: m.Content,
		Title:   voteTitle,
		Options: voteArr,
	}
	res, err := client.Vote(context.Background(), req)
	if err != nil {
		return
	}

	s.ChannelMessageSend(m.ChannelID, res.GetMessage())

	if err != nil {
		log.Print(err)
	}

}
