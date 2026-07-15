package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandRegistry struct {
	commands map[string]Command
}

func NewRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]

	return cmd, ok
}

func (r *CommandRegistry) Commands() []Command {

	commands := make([]Command, 0, len(r.commands))

	for _, command := range r.commands {
		commands = append(commands, command)
		log.Printf("Registered command: /%s", command.Name())

	}

	return commands
}

func (r *CommandRegistry) Sync(session *discordgo.Session, guildId string) error {
	for _, cmd := range r.Commands() {
		_, err := session.ApplicationCommandCreate(
			session.State.User.ID,
			guildId,
			cmd.SlashCommand(),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
