package utility

import (
	"math/rand"

	"github.com/Espectro0/PlumbusBot/internal/ui"
	"github.com/bwmarrin/discordgo"
)

type AboutCommand struct{}

func NewAboutCommand() *AboutCommand {
	return &AboutCommand{}
}

func (c *AboutCommand) Name() string {
	return "about"
}

func (c *AboutCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "about",
		Description: "Check out the info available about me.",
	}
}

func (c *AboutCommand) Execute(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {
	url := "https://static0.srcdn.com/wordpress/wp-content/uploads/2019/09/Rick-and-Morty-Plumbus-SR.jpg?w=1600&h=900&fit=crop"

	var plumbusQuotes = []string{
		"Everybody has a Plumbus in their home.",
		"Plumbus. Everyone has one. Nobody knows why.",
		"If your Plumbus starts screaming... that's normal.",
		"Schleem certified since C-137.",
		"Approved by absolutely no one.",
		"Warning: Do not over-rub the fleeb.",
		"Keep away from Mr. Meeseeks.",
		"Now compatible with Portal Fluid®.",
		"Tested on Jerry. Results inconclusive.",
		"Jerry-proof. Mostly.",
		"Every home deserves a Plumbus.",
		"Ask your nearest Rick for maintenance.",
		"Works better than expected. Nobody knows why.",
	}

	quote := plumbusQuotes[rand.Intn(len(plumbusQuotes))]

	embed := ui.NewEmbed().
		Info().
		Title("🤖 About Plumbus").
		Description(
			"> *Everybody has a Plumbus in their home.*\n\n"+
				"PlumbusBot is a Discord bot inspired by **Rick and Morty**, "+
				"built with **Go** for learning, experimentation, and a bit of interdimensional chaos.",
		).
		Field(
			"✨ Features",
			"🎵 **Music**\n"+
				"🛸 **Rick & Morty**\n"+
				"⚡ **Fast**\n"+
				"❤️ **Open Source**\n",
			false,
		).
		Field(
			"📺 Interdimensional Cable",
			"> *"+quote+"*",
			false,
		).
		Image(url).
		Footer("\"I always wondered how Plumbuses were made.\" — Rick Sanchez").
		Build()

	return ui.RespondEmbed(s, i, embed)
}
