package main

import (
	"log"
	"moegi-discord/config"

	"github.com/bwmarrin/discordgo"
)

type Session interface {
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	MessageReactionAdd(channelID string, messageID string, emojiID string, options ...discordgo.RequestOption) error
}

var VoteEmoji = []string{Eone, Etwo, Ethree, Efour, Efive, Esix, Eseven}

// コマンド例 !vote 旅行先 北海道 東京 沖縄 --Crirona
func Vote(s Session, title string, options []string, cid string) *discordgo.Message {
	voteMsg := ""
	for i, v := range options {
		voteMsg += VoteEmoji[i] + v + "\n"
	}

	sendMsg := &discordgo.MessageEmbed{
		Title:       title,
		Description: voteMsg,
		Color:       1752220,
	}
	msg, err := s.ChannelMessageSendEmbed(cid, sendMsg)
	if err != nil {
		log.Fatal(err)
	}

	for i := range options {
		s.MessageReactionAdd(cid, msg.ID, VoteEmoji[i])
	}

	return msg
}

func Remind(s *discordgo.Session, options []string, cid string, mid string) {
	g, err := s.State.Guild(config.Config.DiscordGuildId)
	if err != nil {
		log.Fatal(err)
	}
	responed := make(map[string]bool, len(g.Members))
	for _, mem := range g.Members {
		responed[mem.User.ID] = false
	}

	// sendMsgにリアクションをしていない人にメンションする
	for i := range options {
		reacters, err := s.MessageReactions(cid, mid, VoteEmoji[i], 100, "", "")
		if err != nil {
			log.Fatal(err)
		}
		for _, u := range reacters {
			responed[u.ID] = true
		}
	}
	msg := ""
	for k, v := range responed {
		if !v {
			msg += "<@" + k + "> "
		}
	}
	s.ChannelMessageSend(cid, msg)
	s.ChannelMessageSend(cid, "なぜ回答しないんだい？みんなは回答しているよ")
}
