# PlumbusBot

A Discord bot built with [discordgo](https://github.com/bwmarrin/discordgo).

## Features

- **`/ping`** — Check if the bot is online.
- **`/about`** — Check out the info available about the bot.
- **`/character [name] [id]`** — Displays a character from *Rick and Morty* via the [Rick and Morty API](https://rickandmortyapi.com). Supports searching by name or ID, random character, and paginated results for multiple matches.
- **Lavalink** — Ready for music playback (configurable via environment variables).

## Requirements

- Go 1.26+
- Lavalink server (for music features)

## Setup

1. Clone the repository:

```bash
git clone https://github.com/Espectro0/PlumbusBot.git
cd PlumbusBot
```

2. Copy the environment file and fill in your credentials:

```bash
cp .env.example .env
```

3. Edit `.env`:

| Variable | Description |
|----------|-------------|
| `DISCORD_TOKEN` | Your Discord bot token |
| `DISCORD_GUILD_ID` | (Optional) Guild ID for guild-only commands |
| `LAVALINK_HOST` | Lavalink server host |
| `LAVALINK_PORT` | Lavalink server port |
| `LAVALINK_PASSWORD` | Lavalink server password |

4. Build and run:

```bash
go build -o plumbusbot ./cmd/bot
./plumbusbot
```

## Structure

```
cmd/bot/          — Entry point
internal/
  app/            — Application bootstrap and lifecycle
  config/         — Environment configuration
  discord/        — Discord session setup
  ui/             — Embed builder and response helpers
  commands/
    paginator/    — Reusable paginated embed system with button navigation
    utility/      — /ping and /about commands
    rickandmorty/ — Rick and Morty API integration and /character command
      api/        — HTTP client with caching and timeout
```

## Caching

The `/character` command caches API responses in memory for 5 minutes, making repeated searches instant.
