# QuotoBot 🤖📝

QuotoBot is a Telegram bot for collecting and sharing quotes with your community. It allows users to add quotes, vote on them, search for quotes, and more. The bot includes user authentication through ViaRézo OAuth to ensure only authorized users can interact with it.

## Features ✨

- User authentication via ViaRézo OAuth
- Add new quotes with author attribution
- Vote system for quotes
- Search quotes by content or author
- View random quotes
- Get top voted quotes
- Secure registration process
- Docker-ready for easy deployment

## Architecture 🏗️

QuotoBot consists of two main components:

1. **Telegram Bot**: Handles all user interactions through Telegram commands
2. **Web Server**: Manages user registration and authentication via OAuth

### Directory Structure 📁

```
quotobot/
├── cmd/                # Command-line entry points
│   ├── bot/            # Bot command
│   ├── server/         # Server command
├── internal/           # Application code
│   ├── bot/            # Bot implementation
│   ├── server/         # Web server implementation
├── pkg/                # Shared packages
│   ├── config/         # Configuration handling
│   ├── database/       # Database models and connection
│   ├── logger/         # Logging utilities
```

## Setup and Installation 🚀

### Prerequisites

- Go 1.24+
- SQLite
- A registered Telegram bot token (from @BotFather)
- OAuth credentials from ViaRézo

### Configuration

1. Create a `config.yaml` file based on the provided template:

```yaml
bot:
  token: your_telegram_bot_token
  chat_id: your_telegram_chat_id
  base_url: your_server_domain
  hmac_secret: your_hmac_secret_key

server:
  session_secret: your_session_secret
  hmac_secret: your_hmac_secret_key
  provider_url: https://moncompte.viarezo.fr
  client_id: your_oauth_client_id
  client_secret: your_oauth_client_secret
  redirect_url: https://your_server_domain/oauth/callback
```

### Running Locally

1. Build and run the bot:

```bash
go run main.go bot
```

2. Build and run the server:

```bash
go run main.go server
```

### Docker Deployment 🐳

QuotoBot includes a Docker setup for easy deployment:

```bash
# Build the image
docker-compose build

# Run the services
docker-compose up -d
```

This will start both the bot and server components, along with a Traefik reverse proxy for HTTPS support.

## Usage 📋

### Available Commands

| Command     | Description           | Parameters                                      |
| ----------- | --------------------- | ----------------------------------------------- |
| `/register` | Register with the bot | None                                            |
| `/add`      | Add a new quote       | Format: `/add \| quote content \| \| author \|` |
| `/random`   | Get random quotes     | Optional: number of quotes to retrieve          |
| `/last`     | Get latest quotes     | Optional: number of quotes to retrieve          |
| `/get`      | Get a quote by ID     | Quote ID                                        |
| `/search`   | Search quotes         | Search term, optional: number of results        |
| `/vote`     | Vote for a quote      | Quote ID                                        |
| `/unvote`   | Remove your vote      | Quote ID                                        |
| `/score`    | View vote count       | Quote ID                                        |
| `/top`      | View top quotes       | Optional: number of quotes to retrieve          |

### Registration Process

1. Start a private conversation with the bot
2. Use the `/register` command
3. Click the provided link
4. Authenticate with your ViaRézo account
5. Complete the registration process

## Acknowledgments 🙏

- Built with [go-telegram-bot](https://github.com/go-telegram/bot)
- Uses [GORM](https://gorm.io/) for database access
- Authentication powered by [go-oidc](https://github.com/coreos/go-oidc)

---

Built with ❤️ for ViaRézo
