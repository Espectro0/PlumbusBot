package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string

	SlashCommand() *discordgo.ApplicationCommand

	Execute(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) error
}
