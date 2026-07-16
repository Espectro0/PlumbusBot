package app

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Espectro0/PlumbusBot/internal/commands"
	"github.com/Espectro0/PlumbusBot/internal/commands/paginator"
	"github.com/Espectro0/PlumbusBot/internal/commands/rickandmorty"
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
	Paginator       *paginator.Paginator
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

	router := commands.NewRouter(registry)

	p := paginator.New(session, 5*time.Minute)

	utility.Register(registry)
	rickandmorty.Register(registry, p)

	session.AddHandler(router.Handle)
	session.AddHandler(p.Handle)
	session.AddHandler(discord.ReadyHandler(registry))

	return &App{
		Config:          cfg,
		Discord:         session,
		CommandRegistry: registry,
		Router:          router,
		Paginator:       p,
	}, nil

}

func Run(a *App) error {

	printBanner()

	if err := a.Discord.Open(); err != nil {
		return err
	}

	defer a.Discord.Close()
	defer a.Paginator.Close()

	if err := a.CommandRegistry.Sync(a.Discord, a.Config.GuildId); err != nil {
		log.Printf("Failed to sync commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Shutting down...")

	return nil
}
