package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken     string
	GuildId          string
	LavalinkHost     string
	LavalinkPort     string
	LavalinkPassword string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		DiscordToken:     mustGetenv("DISCORD_TOKEN"),
		GuildId:          os.Getenv("DISCORD_GUILD_ID"),
		LavalinkHost:     mustGetenv("LAVALINK_HOST"),
		LavalinkPort:     mustGetenv("LAVALINK_PORT"),
		LavalinkPassword: mustGetenv("LAVALINK_PASSWORD"),
	}, nil
}

func mustGetenv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Default().Fatalf("Environment variable %s required...", key)
	}

	return value
}
