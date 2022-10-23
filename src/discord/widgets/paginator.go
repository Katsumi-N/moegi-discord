// ref https://github.com/Necroforger/dgwidgets/blob/master/paginator.go
package widgets

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// emoji constants
const (
	NavPlus        = "➕"
	NavPlay        = "▶"
	NavPause       = "⏸"
	NavStop        = "⏹"
	NavRight       = "➡"
	NavLeft        = "⬅"
	NavUp          = "⬆"
	NavDown        = "⬇"
	NavEnd         = "⏩"
	NavBeginning   = "⏪"
	NavNumbers     = "🔢"
	NavInformation = "ℹ"
	NavSave        = "💾"
)

type Paginator struct {
	sync.Mutex
	Pages []*discordgo.MessageEmbed
	Index int

	Loop   bool
	Widget *Widget

	Ses *discordgo.Session

	DeleteMessageWhenDone   bool
	DeleteReactionsWhenDone bool
	ColorWhenDone           int

	lockToUser bool
	running    bool
}

func NewPaginator(ses *discordgo.Session, channelID string) *Paginator {
	p := &Paginator{
		Ses:                     ses,
		Pages:                   []*discordgo.MessageEmbed{},
		Index:                   0,
		Loop:                    false,
		DeleteMessageWhenDone:   false,
		DeleteReactionsWhenDone: false,
		ColorWhenDone:           -1,
		Widget:                  NewWidget(ses, channelID, nil),
	}
	p.addHandlers()

	return p
}

func (p *Paginator) addHandlers() {
	p.Widget.Handle(NavBeginning, func(w *Widget, r *discordgo.MessageReaction) {
		if err := p.Goto(0); err == nil {
			p.Update()
		}
	})
	p.Widget.Handle(NavLeft, func(w *Widget, r *discordgo.MessageReaction) {
		if err := p.PreviousPage(); err == nil {
			p.Update()
		}
	})
	p.Widget.Handle(NavEnd, func(w *Widget, r *discordgo.MessageReaction) {
		if err := p.Goto(len(p.Pages) - 1); err == nil {
			p.Update()
		}
	})
}

func (p *Paginator) Spawn() error {
	if p.Running() {
		return errors.New("Already running")
	}
	p.Lock()
	p.running = true
	p.Unlock()

	defer func() {
		p.Lock()
		p.running = false
		p.Unlock()

		if p.DeleteMessageWhenDone && p.Widget.Message != nil {
			p.Ses.ChannelMessageDelete(p.Widget.Message.ChannelID, p.Widget.Message.ID)
		} else if p.ColorWhenDone >= 0 {
			if page, err := p.Page(); err == nil {
				page.Color = p.ColorWhenDone
				p.Update()
			}
		}

		if p.DeleteReactionsWhenDone && p.Widget.Message != nil {
			p.Ses.MessageReactionsRemoveAll(p.Widget.ChannelID, p.Widget.Message.ID)
		}
	}()

	page, err := p.Page()
	if err != nil {
		return err
	}
	p.Widget.Embed = page
	return p.Widget.Spawn()
}
func (p *Paginator) Add(embeds ...*discordgo.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

func (p *Paginator) Page() (*discordgo.MessageEmbed, error) {
	p.Lock()
	defer p.Unlock()

	if p.Index < 0 || p.Index >= len(p.Pages) {
		return nil, errors.New("Index out of bounds")
	}
	return p.Pages[p.Index], nil
}

func (p *Paginator) PreviousPage() error {
	p.Lock()
	defer p.Unlock()

	if p.Index-1 >= 0 && p.Index-1 < len(p.Pages) {
		p.Index--
		return nil
	}

	if p.Loop {
		p.Index = len(p.Pages) - 1
		return nil
	}

	return errors.New("Index out of bounds")
}
func (p *Paginator) Goto(index int) error {
	p.Lock()
	defer p.Unlock()
	if index < 0 || index >= len(p.Pages) {
		return errors.New("Index out of bounds")
	}
	p.Index = index
	return nil
}

func (p *Paginator) Update() error {
	if p.Widget.Message == nil {
		return errors.New("nil message")
	}
	page, err := p.Page()
	if err != nil {
		return err
	}
	_, err = p.Widget.UpdateEmbed(page)
	return err
}

// Running returns the running status of the paginator
func (p *Paginator) Running() bool {
	p.Lock()
	running := p.running
	p.Unlock()
	return running
}

// SetPageFooters sets the footer of each embed to
// Be its page number out of the total length of the embeds.
func (p *Paginator) SetPageFooters() {
	for index, embed := range p.Pages {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("#[%d / %d]", index+1, len(p.Pages)),
		}
	}
}
