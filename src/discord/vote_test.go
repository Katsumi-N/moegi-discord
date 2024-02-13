package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

type MockSession struct {
	MockChannelMessageSendEmbed func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	MockMessageReactionAdd      func(channelID string, messageID string, emojiID string, options ...discordgo.RequestOption) error
	MockChannelMessageSend      func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

func (m *MockSession) ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	// return &discordgo.Message{Embeds: []*discordgo.MessageEmbed{embed}}, nil
	return m.MockChannelMessageSendEmbed(channelID, embed)
}

func (m *MockSession) MessageReactionAdd(channelID string, messageID string, emojiID string, options ...discordgo.RequestOption) error {
	return m.MockMessageReactionAdd(channelID, messageID, emojiID)
}

func (m *MockSession) ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return m.ChannelMessageSend(channelID, content)
}

func TestVote(t *testing.T) {
	mockSession := &MockSession{}

	type args struct {
		s       MockSession
		title   string
		options []string
		cid     string
	}

	// Table(テストケース)
	tests := []struct {
		name string
		args args
		want *discordgo.MessageEmbed
	}{
		// TODO: Add test cases.
		{
			name: "Vote created with correct Message",
			args: args{*mockSession, "test", []string{"fish", "beef", "chicken"}, "channelId"},
			want: &discordgo.MessageEmbed{
				Title:       "vote_title",
				Description: "vote description",
				Color:       1752220,
			},
		},
	}

	for _, tt := range tests {
		s := &MockSession{
			MockChannelMessageSendEmbed: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
				if channelID != tt.args.cid || embed != tt.want {
					t.Errorf("ChannelMessageSendEmbed was called with incorrect arguments. got: %v, %v, want: %v, %v", channelID, embed, tt.args.cid, tt.want)
				}
				return nil, nil
			},
		}
		// Failed to read file: open config.ini: no such file or directory
		t.Run(tt.name, func(t *testing.T) {
			Vote(s, tt.args.title, tt.args.options, tt.args.cid)
		})
	}
}
