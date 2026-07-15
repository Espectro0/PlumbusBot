package discord

import (
	"log"

	"github.com/Espectro0/PlumbusBot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(registry *commands.CommandRegistry) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: "online",
			Activities: []*discordgo.Activity{
				{
					Name: "Rick and Morty 🛸",
					Type: discordgo.ActivityTypeWatching,
				},
			},
		})
		if err != nil {
			return
		}

		log.Printf("Connected as %s", s.State.User.Username)
	}
}
