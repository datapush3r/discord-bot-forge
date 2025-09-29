package main

import (
	"log"
	"os"
	"time"

	"discord-bot-forge/commands"
	"discord-bot-forge/core"
	"discord-bot-forge/modules"
	"discord-bot-forge/web"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get bot token from environment
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
	}

	// Create bot configuration
	config := &core.Config{
		Token:     token,
		Prefix:    "!",
		OwnerID:   os.Getenv("BOT_OWNER_ID"),
		DebugMode: os.Getenv("DEBUG") == "true",
		Version:   "1.0.0",
	}

	// Create DiscordBotForge instance
	bot, err := core.NewBot(config)
	if err != nil {
		log.Fatal("Error creating DiscordBotForge:", err)
	}

	// Register commands
	bot.RegisterCommand(&commands.PingCommand{})
	bot.RegisterCommand(commands.NewHelpCommand(bot))
	bot.RegisterCommand(commands.NewInfoCommand(bot))

	// Register modules
	bot.RegisterModule(modules.NewLoggingModule("discord-bot-forge.log"))
	bot.RegisterModule(modules.NewStatsModule())

	// Add middleware
	bot.AddMiddleware(core.NewCooldownMiddleware(2 * time.Second))
	bot.AddMiddleware(core.NewLoggingMiddleware())

	// Start web interface
	webServer := web.NewWebServer(bot, "8080")
	go func() {
		log.Println("🌐 Starting web interface on http://localhost:8080")
		if err := webServer.Start(); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()

	// Start the DiscordBotForge bot
	if err := bot.Start(); err != nil {
		log.Fatal("Error starting DiscordBotForge:", err)
	}
}
