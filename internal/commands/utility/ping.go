package utility

import (
	"fmt"

	"github.com/Espectro0/PlumbusBot/internal/ui"
	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check if the bot is online.",
	}
}

func (c *PingCommand) Execute(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {

	latency := s.HeartbeatLatency().Milliseconds()

	embed := ui.NewEmbed().
		Success().
		Description(fmt.Sprintf("🏓 Pong! \n 📨 Latency:  `%d ms`", latency)).
		Build()

	return ui.RespondEmbed(s, i, embed)
}
