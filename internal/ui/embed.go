package ui

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	ColorSuccess = 0x57F287 // Green
	ColorError   = 0xED4245 // Red
	ColorWarning = 0xFEE75C // Yellow
	ColorInfo    = 0x5865F2 // Blurple
)

type Embed struct {
	embed *discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{
		embed: &discordgo.MessageEmbed{
			Color: ColorInfo,
		},
	}
}

func (e *Embed) Success() *Embed {
	e.embed.Color = ColorSuccess
	return e
}

func (e *Embed) Error() *Embed {
	e.embed.Color = ColorError
	return e
}

func (e *Embed) Warning() *Embed {
	e.embed.Color = ColorWarning
	return e
}

func (e *Embed) Info() *Embed {
	e.embed.Color = ColorInfo
	return e
}

func (e *Embed) Color(color int) *Embed {
	e.embed.Color = color
	return e
}

func (e *Embed) Timestamp() *Embed {
	e.embed.Timestamp = time.Now().Format(time.RFC3339)
	return e
}

func (e *Embed) Title(title string) *Embed {
	e.embed.Title = title
	return e
}

func (e *Embed) Description(description string) *Embed {
	e.embed.Description = description
	return e
}

func (e *Embed) Field(name, value string, inline bool) *Embed {
	e.embed.Fields = append(e.embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})

	return e
}

func (e *Embed) Thumbnail(url string) *Embed {
	e.embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}

	return e
}

func (e *Embed) Image(url string) *Embed {
	e.embed.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}

	return e
}

func (e *Embed) Footer(text string) *Embed {
	e.embed.Footer = &discordgo.MessageEmbedFooter{
		Text: text,
	}

	return e
}

func (e *Embed) Author(name, iconURL string) *Embed {
	e.embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconURL,
	}

	return e
}

func (e *Embed) Build() *discordgo.MessageEmbed {
	return e.embed
}
