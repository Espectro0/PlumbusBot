package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
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

	log.Printf("Connected ad %s", s.State.User.Username)
}
