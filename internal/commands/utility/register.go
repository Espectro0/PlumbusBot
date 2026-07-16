package utility

import "github.com/Espectro0/PlumbusBot/internal/commands"

func Register(registry *commands.CommandRegistry) {
	registry.Register(NewPingCommand())
	registry.Register(NewAboutCommand())
}
