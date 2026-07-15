package app

import (
	"os"
	"os/signal"

	"github.com/Espectro0/PlumbusBot/internal/config"
	"github.com/Espectro0/PlumbusBot/internal/discord"
	"github.com/bwmarrin/discordgo"
)

type App struct {
	Config  *config.Config
	Discord *discordgo.Session
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

	session.AddHandler(discord.Ready)

	return &App{
		Config:  cfg,
		Discord: session,
	}, nil

}

func Run(a *App) error {
	if err := a.Discord.Open(); err != nil {
		return err
	}

	defer a.Discord.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	return nil
}
