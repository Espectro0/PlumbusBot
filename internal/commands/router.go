package commands

import "github.com/bwmarrin/discordgo"

type Router struct {
	registry *CommandRegistry
}

func NewRouter(registry *CommandRegistry) *Router {
	return &Router{
		registry: registry,
	}
}

func (r *Router) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name

	command, ok := r.registry.Get(name)

	if !ok {
		return
	}

	_ = command.Execute(s, i)

}
