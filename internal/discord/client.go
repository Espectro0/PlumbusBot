package discord

import "github.com/bwmarrin/discordgo"

func New(token string) (*discordgo.Session, error) {
	return discordgo.New("Bot " + token)
}
