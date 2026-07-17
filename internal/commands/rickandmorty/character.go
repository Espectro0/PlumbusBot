package rickandmorty

import (
	"fmt"
	"strings"

	api "github.com/Espectro0/PlumbusBot/internal/api/rickandmorty"
	"github.com/Espectro0/PlumbusBot/internal/commands/paginator"
	"github.com/Espectro0/PlumbusBot/internal/ui"
	"github.com/bwmarrin/discordgo"
)

type CharacterCommand struct {
	paginator *paginator.Paginator
}

func NewCharacterCommand(p *paginator.Paginator) *CharacterCommand {
	return &CharacterCommand{paginator: p}
}

func (c *CharacterCommand) Name() string {
	return "character"
}

func (c *CharacterCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "character",
		Description: "Displays a character from «Rick and Morty».",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Search character by name.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "Search character by id.",
				Required:    false,
				MinValue:    func() *float64 { v := float64(1); return &v }(),
			},
		},
	}
}

func (c *CharacterCommand) Execute(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {
	data := i.ApplicationCommandData()

	var name string
	var id int
	for _, opt := range data.Options {
		switch opt.Name {
		case "name":
			name = opt.StringValue()
		case "id":
			id = int(opt.IntValue())
		}
	}

	switch {
	case id > 0 && name != "":
		return ui.RespondError(s, i, "Too Many Arguments", "The `name` and `id` options cannot be used together.")
	case id > 0:
		char, err := api.GetCharacterById(id)
		if err != nil {
			return ui.RespondError(s, i, "Character Not Found", err.Error())
		}
		return ui.RespondEmbed(s, i, characterEmbed(char))

	case name != "":
		chars, err := api.GetCharacterByName(name)
		if err != nil {
			return ui.RespondError(s, i, "No Results Found", err.Error())
		}
		if len(chars) == 1 {
			return ui.RespondEmbed(s, i, characterEmbed(&chars[0]))
		}
		renderer := &CharacterRenderer{characters: chars, perPage: 1}
		return c.paginator.Send(s, i, renderer)

	default:
		char, err := api.RandomCharacter()
		if err != nil {
			return ui.RespondError(s, i, "No Results Found", "I couldn't retrieve the requested character.")
		}
		return ui.RespondEmbed(s, i, characterEmbed(char))
	}
}

func characterEmbed(char *api.Character) *discordgo.MessageEmbed {
	return ui.NewEmbed().
		Color(colorForStatus(char.Status)).
		Title(char.Name).
		Thumbnail(char.Image).
		Description(fmt.Sprintf("Episodes: %d", len(char.Episode))).
		Field("Status", char.Status, true).
		Field("Specie", char.Species, true).
		Field("Gender", char.Gender, true).
		Field("Origin", char.Origin.Name, true).
		Field("Location", char.Location.Name, true).
		Footer(fmt.Sprintf("Character ID: %d", char.ID)).
		Build()
}

func colorForStatus(status string) int {
	switch strings.ToLower(status) {
	case "alive":
		return 0x00FF00
	case "dead":
		return 0xFF0000
	default:
		return 0x808080
	}
}
