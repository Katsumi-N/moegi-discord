package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	TenantId       string
	ServerEndpoint string
	ServerId       string
	Username       string
	Password       string
	DiscordToken   string
	Flavor1gb      string
	Flavor4gb      string
	DiscordGuildId string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		TenantId:       cfg.Section("conoha").Key("tenantId").String(),
		ServerEndpoint: cfg.Section("conoha").Key("server_endpoint").String(),
		ServerId:       cfg.Section("conoha").Key("serverId").String(),
		Username:       cfg.Section("conoha").Key("username").String(),
		Password:       cfg.Section("conoha").Key("password").String(),
		Flavor1gb:      cfg.Section("conoha").Key("flavor_1gb").String(),
		Flavor4gb:      cfg.Section("conoha").Key("flavor_4gb").String(),
		DiscordToken:   cfg.Section("discord").Key("token").String(),
		DiscordGuildId: cfg.Section("discord").Key("guildid").String(),
	}

	// deploy to EC2

	// Config = ConfigList{
	// 	TenantId:       os.Getenv("conohatenantId"),
	// 	ServerEndpoint: os.Getenv(""),
	// 	ServerId:       os.Getenv("conohaserverId"),
	// 	Username:       os.Getenv("conohausername"),
	// 	Password:       os.Getenv("conohapassword"),
	// 	Flavor1gb:      os.Getenv("conohaflavor1gb"),
	// 	Flavor4gb:      os.Getenv("conohaflavor4gb"),
	// 	DiscordToken:   os.Getenv("discordtoken"),
	// }
}
