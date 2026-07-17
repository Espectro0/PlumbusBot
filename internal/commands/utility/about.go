package utility

import (
	"context"
	"math/rand"
	"strconv"
	"sync"

	api "github.com/Espectro0/PlumbusBot/internal/commands/rickandmorty/api"
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
	totalCharacters, totalLocations, err := getStats(context.Background())
	if err != nil {
		return ui.RespondError(s, i, "API Error", "Failed to communicate with the Rick and Morty API")
	}

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
		Field("Total Characters", strconv.Itoa(totalCharacters), true).
		Field("Total Locations", strconv.Itoa(totalLocations), true).
		Image(url).
		Footer("\"I always wondered how Plumbuses were made.\" — Rick Sanchez").
		Build()

	return ui.RespondEmbed(s, i, embed)
}

func getStats(ctx context.Context) (int, int, error) {
	var (
		totalCharacters int
		totalLocations  int
		err1            error
		err2            error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		totalCharacters, err1 = api.GetTotalCharacter()
	}()

	go func() {
		defer wg.Done()
		totalLocations, err2 = api.GetTotalLocation()
	}()

	wg.Wait()

	if err1 != nil {
		return 0, 0, err1
	}

	if err2 != nil {
		return 0, 0, err2
	}

	return totalCharacters, totalLocations, nil
}
