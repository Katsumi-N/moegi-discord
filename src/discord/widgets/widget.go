// (ref) https://github.com/Necroforger/dgwidgets/blob/master/widget.go
package widgets

import (
	"errors"
	"log"
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
	log.Print(w.Embed)
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
		select {
		case k := <-nextMessageReactionAddC(w.Ses):
			reaction = k.MessageReaction
		case <-w.Close:
			return nil
		}

		// botのリアクションは無視する
		if reaction.MessageID != w.Message.ID || w.Ses.State.User.ID == reaction.UserID {
			continue
		}

		// mapではキーによるアクセスをした時，二つ目の引数で有無を取得できる
		if v, ok := w.Handlers[reaction.Emoji.Name]; ok {
			go v(w, reaction)
		}

		if w.DeleteReactions {
			go func() {
				time.Sleep(time.Millisecond * 250)
				w.Ses.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
			}()
		}
	}
}

// emojiの名前のハンドラーを追加
func (w *Widget) Handle(emojiName string, handler WidgetHandler) error {
	if _, ok := w.Handlers[emojiName]; !ok {
		w.Keys = append(w.Keys, emojiName)
		w.Handlers[emojiName] = handler
	}
	// widgetが動いてたらメッセージにemojiを追加
	if w.Running() && w.Message != nil {
		return w.Ses.MessageReactionAdd(w.Message.ChannelID, w.Message.ID, emojiName)
	}
	return nil
}

// よくわからん
func (w *Widget) QueryInput(prompt string, userID string, timeout time.Duration) (*discordgo.Message, error) {
	msg, err := w.Ses.ChannelMessageSend(w.ChannelID, "<@"+userID+">,  "+prompt)
	if err != nil {
		return nil, err
	}
	defer func() {
		w.Ses.ChannelMessageDelete(msg.ChannelID, msg.ID)
	}()

	timeoutChan := make(chan int)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- 0
	}()

	for {
		select {
		case userMsg := <-nextMessageCreateC(w.Ses):
			if userMsg.Author.ID != userID {
				continue
			}
			w.Ses.ChannelMessageDelete(userMsg.ChannelID, userMsg.ID)
			return userMsg.Message, nil
		case <-timeoutChan:
			return nil, errors.New("time out")
		}
	}
}

func (w *Widget) UpdateEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	if w.Message == nil {
		return nil, errors.New("nil")
	}
	return w.Ses.ChannelMessageEditEmbed(w.ChannelID, w.Message.ID, embed)
}
func nextMessageCreateC(s *discordgo.Session) chan *discordgo.MessageCreate {
	out := make(chan *discordgo.MessageCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		out <- e
	})
	return out
}

func nextMessageReactionAddC(s *discordgo.Session) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}
