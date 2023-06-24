package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"moegi-discord/chatgpt"
	"moegi-discord/config"
	conohapb "moegi-discord/pkg/grpc"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	memberOnlineTimestamps = make(map[string]time.Time)
	memberMutex            = &sync.Mutex{}
)

func main() {
	// discord bot
	discordToken := config.Config.DiscordToken
	dg, err := discordgo.New("Bot " + discordToken)
	defer dg.Close()
	// dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers)
	dg.Identify.Intents = discordgo.IntentsAll

	// slash commands
	appId, guildId := config.Config.DiscordAppId, config.Config.DiscordGuildId
	_, err = dg.ApplicationCommandBulkOverwrite(appId, guildId,
		[]*discordgo.ApplicationCommand{
			{
				Name:        "intro",
				Description: "自己紹介するよ",
			},
			{
				Name:        "conoha",
				Description: "マイクラサーバーを確認/起動/停止します",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "command",
						Description: "使用コマンドを選んでね",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "サーバーの状態を確認",
								Value: "status",
							},
							{
								Name:  "サーバーを起動",
								Value: "start",
							},
							{
								Name:  "サーバーを停止",
								Value: "stop",
							},
						},
					},
				},
			},
			{
				Name:        "moriage",
				Description: "moriage隊長、動きます",
			},
			{
				Name:        "vote",
				Description: "投票つくるよ",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "タイトル",
						Description: "投票のタイトル",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "選択肢",
						Description: "半角空白区切りで選択肢を入力 例 東京 名古屋 沖縄",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "クリロナ",
						Description: "なぜ笑うんだい?",
					},
				},
			},
		})
	if err != nil {
		log.Fatal(err)
		return
	}

	dg.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "intro":
			intro(s, i)
		case "conoha":
			conoha(s, i, data.Options[0].StringValue())
		case "moriage":
			moriage(s, i)
		case "vote":
			vote(s, i, data.Options[0].StringValue(), data.Options[1].StringValue(), data.Options[2].BoolValue())
		}
	})
	dg.AddHandler(ChatGPT)
	dg.AddHandler(checkOnline)
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

func intro(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg, err := ioutil.ReadFile("self-intro.txt")
	if err != nil {
		log.Fatal("can't read self-intro.txt")
		return
	}

	err = s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: string(msg),
			},
		})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func conoha(s *discordgo.Session, i *discordgo.InteractionCreate, cmd string) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "サーバーの" + cmd + " を開始します!",
			},
		})
	conn, err := makeGrpcConnection("8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := conohapb.NewConohaServiceClient(conn)

	req := &conohapb.MinecraftRequest{
		Command: cmd,
	}
	stream, err := client.Minecraft(context.Background(), req)
	if err != nil {
		return
	}
	for {
		res, err := stream.Recv()
		s.ChannelMessageSend(i.ChannelID, res.GetMessage())
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			break
		}
	}

	if err != nil {
		log.Print(err)
	}
}

func vote(s *discordgo.Session, i *discordgo.InteractionCreate, title string, opt string, crirona bool) {
	options := strings.Split(opt, " ")
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "投票を作成するよ",
			},
		})
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := Vote(s, title, options, i.ChannelID)

	if crirona {
		time.AfterFunc(5*time.Minute, func() {
			Remind(s, options, i.ChannelID, msg.ID)
		})
	}
}

func moriage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.ChannelMessageDelete(i.ChannelID, i.ID)
	msg := []string{"みんなのために盛り上げるぜ！", "みんな集まれー", "なぜ集まらないんだい？私は暇だよ？", "あほくさ"}
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(len(msg))
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "@everyone " + msg[randomNum],
			},
		})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func ChatGPT(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, "<@"+config.Config.DiscordAppId+">") || len(strings.Split(m.Content, "\n")) < 2 {
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

func checkOnline(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	memberMutex.Lock()
	defer memberMutex.Unlock()

	if m.User == nil {
		return
	}

	if m.Status == discordgo.StatusOnline {
		lastOnline, exists := memberOnlineTimestamps[m.User.ID]

		// 連続通知を避ける
		if exists && time.Since(lastOnline) < 30*time.Minute {
			return
		}
		user, err := s.User(m.User.ID)
		if err != nil {
			log.Println("error retrieving user:", err)
		}
		msg := fmt.Sprintf("%s がオンラインだよ! 囲めー!!", user.Username)
		s.ChannelMessageSend(config.Config.AttendanceChannelId, msg)
		memberOnlineTimestamps[m.User.ID] = time.Now()
	}
}
