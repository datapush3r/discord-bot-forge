package core

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot represents the main DiscordBotForge bot instance
type Bot struct {
	Session    *discordgo.Session
	Commands   map[string]Command
	Modules    []Module
	Config     *Config
	Middleware []Middleware
	Version    string
}

// Config holds bot configuration
type Config struct {
	Token     string
	Prefix    string
	OwnerID   string
	DebugMode bool
	Version   string
}

// Command interface defines the structure for bot commands
type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error
	Permissions() []string
	Cooldown() int // seconds
	Category() string
}

// Module interface defines the structure for bot modules
type Module interface {
	Name() string
	Initialize(bot *Bot) error
	Shutdown() error
	Version() string
}

// Middleware interface for request processing
type Middleware interface {
	Process(s *discordgo.Session, m *discordgo.MessageCreate, next func()) error
	Name() string
}

// NewBot creates a new DiscordBotForge bot instance
func NewBot(config *Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	bot := &Bot{
		Session:    session,
		Commands:   make(map[string]Command),
		Modules:    make([]Module, 0),
		Config:     config,
		Middleware: make([]Middleware, 0),
		Version:    config.Version,
	}

	return bot, nil
}

// Start initializes and starts the DiscordBotForge bot
func (b *Bot) Start() error {
	log.Printf("🔥 DiscordBotForge v%s starting up...", b.Version)
	
	// Add message handler
	b.Session.AddHandler(b.messageHandler)

	// Open connection
	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	// Initialize all modules
	for _, module := range b.Modules {
		if err := module.Initialize(b); err != nil {
			log.Printf("Error initializing module %s: %v", module.Name(), err)
		} else {
			log.Printf("✅ Module '%s' v%s initialized", module.Name(), module.Version())
		}
	}

	log.Printf("🚀 DiscordBotForge is now running! Press CTRL+C to exit.")
	
	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return b.Shutdown()
}

// Shutdown gracefully shuts down the DiscordBotForge bot
func (b *Bot) Shutdown() error {
	log.Println("🛑 Shutting down DiscordBotForge...")

	// Shutdown all modules
	for _, module := range b.Modules {
		if err := module.Shutdown(); err != nil {
			log.Printf("Error shutting down module %s: %v", module.Name(), err)
		} else {
			log.Printf("✅ Module '%s' shutdown complete", module.Name())
		}
	}

	// Close Discord session
	return b.Session.Close()
}

// RegisterCommand adds a command to the bot
func (b *Bot) RegisterCommand(cmd Command) {
	b.Commands[cmd.Name()] = cmd
	log.Printf("⚡ Registered command: %s", cmd.Name())
}

// RegisterModule adds a module to the bot
func (b *Bot) RegisterModule(module Module) {
	b.Modules = append(b.Modules, module)
	log.Printf("🔧 Registered module: %s", module.Name())
}

// AddMiddleware adds middleware to the bot
func (b *Bot) AddMiddleware(middleware Middleware) {
	b.Middleware = append(b.Middleware, middleware)
	log.Printf("🛡️ Added middleware: %s", middleware.Name())
}

// messageHandler processes incoming messages
func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Check if message starts with prefix
	if len(m.Content) < len(b.Config.Prefix) || m.Content[:len(b.Config.Prefix)] != b.Config.Prefix {
		return
	}

	// Parse command and arguments
	args := parseArgs(m.Content[len(b.Config.Prefix):])
	if len(args) == 0 {
		return
	}

	commandName := args[0]
	commandArgs := args[1:]

	// Find command
	cmd, exists := b.Commands[commandName]
	if !exists {
		return
	}

	// Process middleware
	for _, middleware := range b.Middleware {
		if err := middleware.Process(s, m, func() {
			// Execute command
			if err := cmd.Execute(s, m, commandArgs); err != nil {
				log.Printf("Error executing command %s: %v", commandName, err)
				s.ChannelMessageSend(m.ChannelID, "❌ An error occurred while executing the command.")
			}
		}); err != nil {
			log.Printf("Middleware error: %v", err)
			return
		}
	}
}

// parseArgs splits a string into arguments, respecting quotes
func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	
	for i, r := range input {
		if r == '"' {
			inQuotes = !inQuotes
		} else if r == ' ' && !inQuotes {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}
	
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	
	return args
}

// GetCommandCategories returns a map of commands grouped by category
func (b *Bot) GetCommandCategories() map[string][]Command {
	categories := make(map[string][]Command)
	
	for _, cmd := range b.Commands {
		category := cmd.Category()
		if category == "" {
			category = "General"
		}
		categories[category] = append(categories[category], cmd)
	}
	
	return categories
}
