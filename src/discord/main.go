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
	"time"

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

// コマンド例 !vote 旅行先 北海道 東京 沖縄 --Crirona
func Vote(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "!vote") {
		return
	}

	voteArr := strings.Split(m.Content, " ")
	isCrirona := false
	// --time引数を検知したい
	for i, v := range voteArr {
		if v == "--crirona" {
			isCrirona = true
			voteArr = voteArr[:i]
			break
		}
	}

	if len(voteArr) <= 2 {
		s.ChannelMessageSend(m.ChannelID, "選択肢を入力してね！")
	}
	if len(voteArr) >= 10 {
		s.ChannelMessageSend(m.ChannelID, "選択肢は7個以下でお願いします！")
	}
	voteOptions := voteArr[2:]

	voteMsg := m.Message.Author.Username + "が作った投票だよ！\n"

	voteEmoji := []string{Eone, Etwo, Ethree, Efour, Efive, Esix, Eseven}
	for i, v := range voteOptions {
		voteMsg += voteEmoji[i] + v + "\n"
	}

	sendMsg := &discordgo.MessageEmbed{
		Title:       voteArr[1],
		Description: voteMsg,
		Color:       1752220,
	}
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, sendMsg)
	if err != nil {
		log.Fatal(err)
	}
	for i := range voteOptions {
		s.MessageReactionAdd(m.ChannelID, msg.ID, voteEmoji[i])
	}

	if isCrirona {
		time.AfterFunc(3*time.Minute, func() {
			// TODO: sendMsgにリアクションをしていない人にメンションする
			s.ChannelMessageSend(m.ChannelID, "なぜ回答しないんだい？みんなは回答しているよ")
		})
	}
}
