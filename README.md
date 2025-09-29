# 🔥 DiscordBotForge

A powerful, modular framework for forging Discord bots in Go. Built with extensibility and ease of use in mind, featuring a comprehensive web interface for bot management.

## ✨ Features

- **🔥 Modular Architecture**: Easy to extend with custom commands and modules
- **⚡ Command System**: Simple interface for creating commands with categories
- **🛡️ Middleware Support**: Built-in cooldown, permission, and logging middleware
- **🔧 Module System**: Pluggable modules for logging, statistics, and more
- **🌐 Web Interface**: Beautiful, responsive web dashboard for bot management
- **📡 Real-time Updates**: WebSocket-powered live status updates
- **⚙️ Configuration**: Environment-based configuration with `.env` support
- **🛑 Graceful Shutdown**: Proper cleanup on bot shutdown
- **📊 Statistics**: Built-in usage tracking and statistics
- **📝 Logging**: Comprehensive logging system

## 🚀 Quick Start

### Option 1: Docker (Recommended)

1. **Prerequisites**:
   - Docker and Docker Compose installed
   - Discord bot token from [Discord Developer Portal](https://discord.com/developers/applications)

2. **Setup DiscordBotForge**:
   ```bash
   cd discord-bot-forge
   cp env.docker .env
   # Edit .env with your Discord bot token
   ```

3. **Start with Docker**:
   ```bash
   # Quick start script
   ./docker-start.sh
   
   # Or manually
   make docker-prod    # Production stack
   make docker-dev     # Development stack
   ```

4. **Access the web interface**:
   - **Web Interface**: http://localhost:8080

### Option 2: Local Development

1. **Setup DiscordBotForge**:
   ```bash
   cd discord-bot-forge
   go mod tidy
   ```

2. **Create a Discord Application**:
   - Go to [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application
   - Go to "Bot" section and create a bot
   - Copy the bot token

3. **Configure your bot**:
   ```bash
   cp env.example .env
   # Edit .env with your bot token
   ```

4. **Forge your first bot**:
   ```bash
   cd examples/simple_bot
   go run main.go
   ```

5. **Access the web interface**:
   - Open your browser to `http://localhost:8080`
   - Monitor your bot, manage commands, and view logs in real-time!

## 🌐 Web Interface

DiscordBotForge includes a comprehensive web interface that provides:

### Dashboard
- **Real-time Status**: Live bot status, uptime, and statistics
- **Quick Actions**: Restart, stop, and refresh bot functionality
- **Activity Feed**: Recent bot activity and events
- **Statistics Overview**: Commands, modules, and middleware counts

### Commands Management
- **Command List**: View all registered commands with details
- **Add Commands**: Create new commands through the web interface
- **Command Details**: View usage, permissions, and cooldown settings
- **Category Organization**: Commands organized by categories

### Modules Management
- **Module Status**: View all loaded modules and their status
- **Module Control**: Start, stop, and restart individual modules
- **Available Modules**: Browse and install additional modules
- **Version Management**: Track module versions and updates

### Live Logs
- **Real-time Logging**: View bot logs as they happen
- **Log Filtering**: Filter by level, search terms, and time
- **Log Statistics**: Count of different log levels
- **Export Functionality**: Download logs for analysis

### Settings
- **Bot Configuration**: Modify prefix, owner ID, and debug settings
- **Security Settings**: Configure rate limiting and permissions
- **Database Settings**: Database connection and configuration
- **Web Interface Settings**: Port, WebSocket, and interface options

## 🏗️ Architecture

### Core Components

- **Bot**: Main DiscordBotForge instance that manages commands, modules, and middleware
- **Command**: Interface for bot commands with categories and metadata
- **Module**: Interface for bot modules (logging, stats, etc.)
- **Middleware**: Request processing pipeline with built-in middleware
- **Web Server**: HTTP server with REST API and WebSocket support
- **Web Interface**: Responsive HTML/CSS/JS dashboard

### Built-in Commands

- **ping**: Check bot latency
- **help**: Show available commands with categories
- **info**: Display DiscordBotForge information and statistics

### Built-in Modules

- **Logging**: Logs all messages to file and console
- **Statistics**: Tracks usage stats (messages, commands, uptime)

### Built-in Middleware

- **Cooldown**: Rate limiting for commands
- **Permissions**: Check user permissions before command execution
- **Logging**: Log command usage
- **OwnerOnly**: Restrict commands to bot owner only

## 🔨 Creating Commands

```go
type MyCommand struct{}

func (c *MyCommand) Name() string {
    return "mycommand"
}

func (c *MyCommand) Description() string {
    return "My awesome command"
}

func (c *MyCommand) Usage() string {
    return "mycommand <arg1> <arg2>"
}

func (c *MyCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
    s.ChannelMessageSend(m.ChannelID, "Hello from DiscordBotForge!")
    return nil
}

func (c *MyCommand) Permissions() []string {
    return []string{} // No special permissions required
}

func (c *MyCommand) Cooldown() int {
    return 5 // 5 second cooldown
}

func (c *MyCommand) Category() string {
    return "Fun" // Command category
}
```

## 🔧 Creating Modules

```go
type MyModule struct{}

func (m *MyModule) Name() string {
    return "MyModule"
}

func (m *MyModule) Version() string {
    return "1.0.0"
}

func (m *MyModule) Initialize(bot *core.Bot) error {
    // Initialize your module
    return nil
}

func (m *MyModule) Shutdown() error {
    // Cleanup resources
    return nil
}
```

## 🛡️ Using Middleware

```go
// Add cooldown middleware (2 seconds)
bot.AddMiddleware(core.NewCooldownMiddleware(2 * time.Second))

// Add permission middleware
bot.AddMiddleware(core.NewPermissionMiddleware([]string{"ADMINISTRATOR"}))

// Add logging middleware
bot.AddMiddleware(core.NewLoggingMiddleware())

// Add owner-only middleware
bot.AddMiddleware(core.NewOwnerOnlyMiddleware(bot.Config.OwnerID))
```

## 🌐 Web Interface Integration

```go
// Start web interface
webServer := web.NewWebServer(bot, "8080")
go func() {
    log.Println("🌐 Starting web interface on http://localhost:8080")
    if err := webServer.Start(); err != nil {
        log.Printf("Web server error: %v", err)
    }
}()
```

## 📁 Project Structure

```
discord-bot-forge/
├── core/                 # Core framework components
│   ├── bot.go           # Main bot implementation
│   └── middleware.go    # Built-in middleware
├── commands/            # Built-in commands
│   └── basic.go        # Basic commands (ping, help, info)
├── modules/            # Built-in modules
│   ├── logging.go      # Logging module
│   └── stats.go        # Statistics module
├── web/                # Web interface
│   ├── server.go       # Web server and API
│   ├── templates/      # HTML templates
│   │   ├── dashboard.html
│   │   ├── commands.html
│   │   ├── modules.html
│   │   ├── logs.html
│   │   └── settings.html
│   └── static/         # Static assets
│       ├── css/style.css
│       └── js/app.js
├── examples/           # Example bots
│   ├── simple_bot/     # Simple example bot
│   └── advanced_bot/   # Advanced example bot
├── go.mod              # Go module file
├── env.example         # Environment configuration example
├── Makefile           # Build and run commands
└── README.md          # This file
```

## 🎯 Examples

- **Simple Bot**: Basic bot with ping, help, and info commands + web interface on port 8080
- **Advanced Bot**: Bot with additional features and web interface on port 8081

## 🐳 Docker Deployment

DiscordBotForge includes comprehensive Docker support for easy deployment:

### Production Stack
- **DiscordBotForge**: Main bot application

### Quick Docker Commands
```bash
# Start production stack
make docker-prod

# Start development stack
make docker-dev

# View logs
make docker-logs

# Stop services
make docker-stop

# Clean up
make docker-clean
```

### Docker Services
| Service | Port | Description |
|---------|------|-------------|
| DiscordBotForge | 8080 | Web interface |

For detailed Docker documentation, see [DOCKER.md](DOCKER.md).

## 🔌 API Endpoints

The web interface provides REST API endpoints:

- `GET /api/status` - Get bot status and statistics
- `GET /api/commands` - Get all registered commands
- `GET /api/modules` - Get all loaded modules
- `GET /api/logs` - Get recent logs
- `POST /api/restart` - Restart the bot
- `POST /api/stop` - Stop the bot
- `WebSocket /ws` - Real-time updates

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

---

**🔥 Forge your Discord bots with DiscordBotForge!**
