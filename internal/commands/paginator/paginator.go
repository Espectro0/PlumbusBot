package paginator

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	ButtonFirst    = "pag:first"
	ButtonPrevious = "pag:prev"
	ButtonPage     = "pag:page"
	ButtonNext     = "pag:next"
	ButtonLast     = "pag:last"
	ButtonClose    = "pag:close"
)

type Renderer interface {
	Total() int
	Embed(page int) *discordgo.MessageEmbed
}

type session struct {
	MessageID string
	UserID    string
	ChannelID string
	Page      int
	Renderer  Renderer
	Total     int
	ExpiresAt time.Time
}

func (s *session) HasPrevious() bool { return s.Page > 0 }
func (s *session) HasNext() bool     { return s.Page < s.Total-1 }

type Paginator struct {
	s        *discordgo.Session
	sessions map[string]*session
	timeout  time.Duration
	mu       sync.Mutex
}

func New(s *discordgo.Session, timeout time.Duration) *Paginator {
	return &Paginator{
		s:        s,
		sessions: make(map[string]*session),
		timeout:  timeout,
	}
}

func (p *Paginator) Send(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	renderer Renderer,
) error {
	ses := &session{
		UserID:    i.Member.User.ID,
		ChannelID: i.ChannelID,
		Renderer:  renderer,
		Total:     renderer.Total(),
		ExpiresAt: time.Now().Add(p.timeout),
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{ses.Renderer.Embed(ses.Page)},
			Components: components(ses),
		},
	})
	if err != nil {
		return err
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return err
	}

	ses.MessageID = msg.ID

	p.mu.Lock()
	p.sessions[msg.ID] = ses
	p.mu.Unlock()

	return nil
}

func (p *Paginator) SendFollowup(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	renderer Renderer,
) error {
	ses := &session{
		UserID:    i.Member.User.ID,
		ChannelID: i.ChannelID,
		Renderer:  renderer,
		Total:     renderer.Total(),
		ExpiresAt: time.Now().Add(p.timeout),
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{ses.Renderer.Embed(ses.Page)},
		Components: components(ses),
	})
	if err != nil {
		return err
	}

	ses.MessageID = msg.ID

	p.mu.Lock()
	p.sessions[msg.ID] = ses
	p.mu.Unlock()

	return nil
}

func (p *Paginator) Edit(
	s *discordgo.Session,
	channelID string,
	messageID string,
	userID string,
	renderer Renderer,
) error {
	p.mu.Lock()
	ses, ok := p.sessions[messageID]
	if !ok {
		ses = &session{
			UserID:    userID,
			ChannelID: channelID,
			Renderer:  renderer,
			Total:     renderer.Total(),
			ExpiresAt: time.Now().Add(p.timeout),
		}
		ses.MessageID = messageID
		p.sessions[messageID] = ses
	}
	p.mu.Unlock()

	comp := components(ses)
	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    channelID,
		ID:         messageID,
		Embeds:     &[]*discordgo.MessageEmbed{ses.Renderer.Embed(ses.Page)},
		Components: &comp,
	})

	return err
}

func (p *Paginator) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	data := i.MessageComponentData()
	if data.ComponentType != discordgo.ButtonComponent || data.CustomID == ButtonPage {
		return
	}

	p.mu.Lock()
	ses, ok := p.sessions[i.Message.ID]
	if !ok || ses.UserID != i.Member.User.ID || time.Now().After(ses.ExpiresAt) {
		if ok {
			delete(p.sessions, i.Message.ID)
		}
		p.mu.Unlock()
		return
	}

	if data.CustomID == ButtonClose {
		delete(p.sessions, i.Message.ID)
		p.mu.Unlock()

		comp := disableComponents(components(ses))
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{ses.Renderer.Embed(ses.Page)},
				Components: comp,
			},
		})
		return
	}

	switch data.CustomID {
	case ButtonFirst:
		ses.Page = 0
	case ButtonPrevious:
		if ses.Page > 0 {
			ses.Page--
		}
	case ButtonNext:
		if ses.Page < ses.Total-1 {
			ses.Page++
		}
	case ButtonLast:
		ses.Page = ses.Total - 1
	}

	ses.ExpiresAt = time.Now().Add(p.timeout)
	p.mu.Unlock()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{ses.Renderer.Embed(ses.Page)},
			Components: components(ses),
		},
	})
}

func (p *Paginator) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for id, ses := range p.sessions {
		comp := disableComponents(components(ses))
		p.s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    ses.ChannelID,
			ID:         id,
			Components: &comp,
		})
	}

	clear(p.sessions)
}

func components(ses *session) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				navButton(ButtonFirst, "⏮", !ses.HasPrevious()),
				navButton(ButtonPrevious, "◀", !ses.HasPrevious()),
				navButton(ButtonNext, "▶", !ses.HasNext()),
				navButton(ButtonLast, "⏭", !ses.HasNext()),
			},
		},
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				closeButton(),
			},
		},
	}
}

func navButton(customID, label string, disabled bool) discordgo.MessageComponent {
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.SecondaryButton,
		CustomID: customID,
		Disabled: disabled,
	}
}

func closeButton() discordgo.MessageComponent {
	return &discordgo.Button{
		Style:    discordgo.DangerButton,
		CustomID: ButtonClose,
		Emoji:    &discordgo.ComponentEmoji{Name: "❌"},
	}
}

func disableComponents(comp []discordgo.MessageComponent) []discordgo.MessageComponent {
	rows := make([]discordgo.MessageComponent, len(comp))
	for i, c := range comp {
		row := c.(*discordgo.ActionsRow)
		newRow := &discordgo.ActionsRow{
			Components: make([]discordgo.MessageComponent, len(row.Components)),
		}
		for j, btn := range row.Components {
			b := *(btn.(*discordgo.Button))
			b.Disabled = true
			newRow.Components[j] = &b
		}
		rows[i] = newRow
	}
	return rows
}
