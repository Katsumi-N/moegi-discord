// (ref) https://github.com/Necroforger/dgwidgets/blob/master/widget.go
package widgets

import (
	"errors"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type WidgetHandler func(*Widget, *discordgo.MessageReaction)
type Widget struct {
	sync.Mutex
	Embed     *discordgo.MessageEmbed
	Message   *discordgo.Message
	Ses       *discordgo.Session
	ChannelID string
	Timeout   time.Duration
	Close     chan bool

	Handlers map[string]WidgetHandler
	Keys     []string

	DeleteReactions bool

	running bool
}

func NewWidget(ses *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) *Widget {
	return &Widget{
		Embed:           embed,
		Ses:             ses,
		ChannelID:       channelID,
		Close:           make(chan bool),
		Handlers:        map[string]WidgetHandler{},
		Keys:            []string{},
		DeleteReactions: true,
	}
}

func (w *Widget) Running() bool {
	w.Lock()
	running := w.running
	w.Unlock()
	return running
}
func (w *Widget) Spawn() error {
	if w.Running() {
		return errors.New("Not running widget")
	}
	w.running = true
	defer func() {
		w.running = false
	}()
	if w.Embed == nil {
		return errors.New("Embed is nil")
	}

	startTime := time.Now()

	msg, err := w.Ses.ChannelMessageSendEmbed(w.ChannelID, w.Embed)
	if err != nil {
		return err
	}
	w.Message = msg

	// リアクションボタンを追加(投票機能と同じ)
	for _, v := range w.Keys {
		w.Ses.MessageReactionAdd(w.Message.ChannelID, w.Message.ID, v)
	}

	var reaction *discordgo.MessageReaction
	for {

	}
}
