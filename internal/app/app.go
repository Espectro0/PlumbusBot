package app

import (
	"log"
	"os"
	"os/signal"

	"github.com/Espectro0/PlumbusBot/internal/commands"
	"github.com/Espectro0/PlumbusBot/internal/commands/utility"
	"github.com/Espectro0/PlumbusBot/internal/config"
	"github.com/Espectro0/PlumbusBot/internal/discord"
	"github.com/bwmarrin/discordgo"
)

type App struct {
	Config          *config.Config
	Discord         *discordgo.Session
	CommandRegistry *commands.CommandRegistry
	Router          *commands.Router
}

func New() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	session, err := discord.New(cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	registry := commands.NewRegistry()

	utility.Register(registry)

	router := commands.NewRouter(registry)

	session.AddHandler(router.Handle)
	session.AddHandler(discord.ReadyHandler(registry))

	return &App{
		Config:          cfg,
		Discord:         session,
		CommandRegistry: registry,
		Router:          router,
	}, nil

}

func Run(a *App) error {
	if err := a.Discord.Open(); err != nil {
		return err
	}

	defer a.Discord.Close()

	if err := a.CommandRegistry.Sync(a.Discord, a.Config.GuildId); err != nil {
		log.Printf("Failed to sync commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Shutting down...")

	return nil
}
