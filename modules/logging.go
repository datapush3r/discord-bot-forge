package modules

import (
	"fmt"
	"log"
	"os"
	"time"

	"discord-bot-forge/core"
	"github.com/bwmarrin/discordgo"
)

// LoggingModule provides logging functionality for DiscordBotForge
type LoggingModule struct {
	logFile *os.File
	logger  *log.Logger
	version string
}

// NewLoggingModule creates a new logging module
func NewLoggingModule(logFile string) *LoggingModule {
	return &LoggingModule{
		version: "1.0.0",
	}
}

func (l *LoggingModule) Name() string {
	return "Logging"
}

func (l *LoggingModule) Version() string {
	return l.version
}

func (l *LoggingModule) Initialize(bot *core.Bot) error {
	// Create log file
	file, err := os.OpenFile("discord-bot-forge.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	
	l.logFile = file
	l.logger = log.New(file, "[DiscordBotForge] ", log.LstdFlags)
	
	// Add message logging handler
	bot.Session.AddHandler(l.messageLogger)
	
	l.logger.Println("Logging module initialized")
	return nil
}

func (l *LoggingModule) Shutdown() error {
	if l.logFile != nil {
		l.logger.Println("Logging module shutting down")
		return l.logFile.Close()
	}
	return nil
}

func (l *LoggingModule) messageLogger(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s#%s in #%s: %s",
		timestamp,
		m.Author.Username,
		m.Author.Discriminator,
		m.ChannelID,
		m.Content,
	)
	
	l.logger.Println(logMessage)
	
	if s.State.User.ID == m.Author.ID {
		return
	}
	
	log.Printf("Message: %s", logMessage)
}

// Log logs a message to the log file
func (l *LoggingModule) Log(message string) {
	if l.logger != nil {
		l.logger.Println(message)
	}
}
