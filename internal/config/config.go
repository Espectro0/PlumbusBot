package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	GuildId      string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		DiscordToken: mustGetenv("DISCORD_TOKEN"),
		GuildId:      os.Getenv("DISCORD_GUILD_ID"),
	}, nil
}

func mustGetenv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Default().Fatalf("Environment variable %s required...", key)
	}

	return value
}
