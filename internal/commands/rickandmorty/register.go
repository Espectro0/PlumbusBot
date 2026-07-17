package rickandmorty

import (
	"github.com/Espectro0/PlumbusBot/internal/commands"
	"github.com/Espectro0/PlumbusBot/internal/commands/paginator"
)

func Register(registry *commands.CommandRegistry, p *paginator.Paginator) {
	registry.Register(NewCharacterCommand(p))
	registry.Register(NewLocationCommand(p))
}
