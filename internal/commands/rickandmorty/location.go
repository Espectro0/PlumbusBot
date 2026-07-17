package rickandmorty

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"github.com/Espectro0/PlumbusBot/internal/commands/paginator"
	api "github.com/Espectro0/PlumbusBot/internal/commands/rickandmorty/api"
	"github.com/Espectro0/PlumbusBot/internal/ui"
	"github.com/bwmarrin/discordgo"
)

type LocationCommand struct {
	paginator *paginator.Paginator
}

func NewLocationCommand(p *paginator.Paginator) *LocationCommand {
	return &LocationCommand{paginator: p}
}

func (c *LocationCommand) Name() string {
	return "location"
}

func (c *LocationCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "location",
		Description: "Displays a location from «Rick and Morty».",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Search location by name.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "Search location by id.",
				Required:    false,
				MinValue:    func() *float64 { v := float64(1); return &v }(),
			},
		},
	}
}

func (l *LocationCommand) Execute(
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
		loc, err := api.GetLocationById(id)
		if err != nil {
			return ui.RespondError(s, i, "Location Not Found", err.Error())
		}
		return ui.RespondEmbed(s, i, locationEmbed(loc))

	case name != "":
		locs, err := api.GetLocationByName(name)
		if err != nil {
			return ui.RespondError(s, i, "No Results Found", err.Error())
		}
		if len(locs) == 1 {
			return ui.RespondEmbed(s, i, locationEmbed(&locs[0]))
		}
		renderer := &LocationRenderer{locations: locs, perPage: 1}
		return l.paginator.Send(s, i, renderer)

	default:
		loc, err := api.RandomLocation()
		if err != nil {
			return ui.RespondError(s, i, "No Results Found", "I couldn't retrieve the requested character.")
		}
		return ui.RespondEmbed(s, i, locationEmbed(loc))
	}
}

func locationEmbed(loc *api.Location) *discordgo.MessageEmbed {
	dimension := loc.Dimension
	if dimension == "" {
		dimension = "Not Available"
	}

	return ui.NewEmbed().
		Color(colorForDimension(loc.Dimension)).
		Title(loc.Name).
		Description(fmt.Sprintf("**Characters Seen:** %d", len(loc.Residents))).
		Field("Location Type", loc.Type, true).
		Field("Dimension", dimension, true).
		Footer(fmt.Sprintf("Location ID: %d", loc.ID)).
		Build()
}

func colorForDimension(dim string) int {
	if strings.ToLower(dim) == "unknown" {
		return 0x808080
	}

	h := fnv.New32a()
	h.Write([]byte(dim))
	hue := float64(h.Sum32() % 360)

	return hslToRGB(hue, 0.7, 0.5)
}

func hslToRGB(h, s, l float64) int {
	c := (1 - abs(2*l-1)) * s
	x := c * (1 - abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return (int((r+m)*255) << 16) | (int((g+m)*255) << 8) | int((b+m)*255)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
