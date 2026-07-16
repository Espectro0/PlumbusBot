package rickandmorty

import (
	"fmt"

	api "github.com/Espectro0/PlumbusBot/internal/commands/rickandmorty/api"
	"github.com/bwmarrin/discordgo"
)

type CharacterRenderer struct {
	characters []api.Character
	perPage    int
}

func (r *CharacterRenderer) Total() int {
	pages := len(r.characters) / r.perPage
	if len(r.characters)%r.perPage != 0 {
		pages++
	}
	return pages
}

func (r *CharacterRenderer) Embed(page int) *discordgo.MessageEmbed {
	start := page * r.perPage
	end := min(start+r.perPage, len(r.characters))
	slice := r.characters[start:end]

	char := &slice[0]
	embed := characterEmbed(char)

	if r.Total() > 1 {
		embed.Footer.Text = fmt.Sprintf("Character ID: %d • Page %d/%d", char.ID, page+1, r.Total())
	}

	return embed
}
