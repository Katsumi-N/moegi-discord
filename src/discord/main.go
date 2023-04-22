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

	"moegi-discord/chatgpt"
	"moegi-discord/config"
	"moegi-discord/discord/widgets"
	conohapb "moegi-discord/pkg/grpc"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MOEGI_ID = "818384700540321802"

func main() {
	// discord bot
	discordToken := config.Config.DiscordToken
	dg, err := discordgo.New("Bot " + discordToken)
	defer dg.Close()
	// dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers)
	dg.Identify.Intents = discordgo.IntentsAll
	dg.AddHandler(Minecraft)
	dg.AddHandler(Introduction)
	dg.AddHandler(vote)
	dg.AddHandler(Widget)
	dg.AddHandler(ChatGPT)
	dg.AddHandler(moriage)
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
	if !strings.Contains(m.Content, "!intro") {
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
func vote(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "!vote") {
		return
	}
	received := strings.Split(m.Content, " ")
	isCrirona := false
	for i, v := range received {
		if v == "--crirona" {
			isCrirona = true
			received = received[:i]
			break
		}
	}
	if len(received) <= 2 {
		return
	} else if len(received) >= 10 {
		return
	}
	title, options := received[1], received[2:]
	msg := Vote(s, title, options, m.ChannelID)

	if isCrirona {
		time.AfterFunc(5*time.Minute, func() {
			Remind(s, options, m.ChannelID, msg.ID)
		})
	}
}

func Widget(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "!maple") {
		return
	}
	s.ChannelMessageSend(m.ChannelID, "イベントを取得中...")

	region := strings.Split(m.Content, " ")[1]
	var (
		maple    *[]MapleInfo
		eventNum int
		err      error
	)
	if region == "jms" {
		maple, eventNum, err = ScrapingEventInfo()
	}

	if err != nil {
		log.Fatal(err)
	}
	s.ChannelMessageSend(m.ChannelID, "ウィジェットを表示するよ！")
	p := widgets.NewPaginator(s, m.ChannelID)
	ma := *maple
	for i := 0; i < eventNum; i++ {
		p.Add(&discordgo.MessageEmbed{
			Title:       ma[i].Title,
			Description: ma[i].Url + "\n\n" + ma[i].Date + "\n\n" + ma[i].Description,
		})
	}

	p.SetPageFooters()

	p.ColorWhenDone = 0xffff

	p.Spawn()
}

func ChatGPT(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "<@"+MOEGI_ID+">") || len(strings.Split(m.Content, "\n")) < 2 {
		return
	}

	msg, err := chatgpt.Chat(strings.Split(m.Content, "\n")[1:])
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("ChatGPT returned:\n", msg)
	s.ChannelMessageSend(m.ChannelID, strings.Join(msg, "\n"))
}

func moriage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "!moriage") {
		return
	}
	s.ChannelMessageDelete(m.ChannelID, m.ID)
	msg := "@everyone \nみんなのために盛り上げるぜ！"
	s.ChannelMessageSend(m.ChannelID, msg)
}
