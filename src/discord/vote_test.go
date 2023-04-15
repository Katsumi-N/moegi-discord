package main

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

type MockSession struct{}

func (m *MockSession) ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	// return &discordgo.Message{Embeds: []*discordgo.MessageEmbed{embed}}, nil
	return nil, nil
}

func (m *MockSession) MessageReactionAdd(channelID string, messageID string, emojiID string, options ...discordgo.RequestOption) error {
	return nil
}

func TestVote(t *testing.T) {
	mockSession := &MockSession{}

	type args struct {
		s       MockSession
		title   string
		options []string
		cid     string
	}

	tests := []struct {
		name string
		args args
		want *discordgo.Message
	}{
		// TODO: Add test cases.
		{
			name: "Correct message",
			args: args{*mockSession, "test", []string{"op1", "op2", "op3"}, "channel"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Vote(&tt.args.s, tt.args.title, tt.args.options, tt.args.cid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vote() = %v, want %v", got, tt.want)
			}
		})
	}
}
