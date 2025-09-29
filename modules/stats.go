package modules

import (
	"fmt"
	"log"
	"sync"
	"time"

	"discord-bot-forge/core"
	"github.com/bwmarrin/discordgo"
)

// StatsModule tracks DiscordBotForge statistics
type StatsModule struct {
	mu           sync.RWMutex
	startTime    time.Time
	messageCount int64
	commandCount int64
	userCount    int64
	version      string
}

// NewStatsModule creates a new stats module
func NewStatsModule() *StatsModule {
	return &StatsModule{
		startTime: time.Now(),
		version:   "1.0.0",
	}
}

func (s *StatsModule) Name() string {
	return "Statistics"
}

func (s *StatsModule) Version() string {
	return s.version
}

func (s *StatsModule) Initialize(bot *core.Bot) error {
	// Add message handler to track messages
	bot.Session.AddHandler(s.messageHandler)
	
	// Add ready handler to track users
	bot.Session.AddHandler(s.readyHandler)
	
	log.Println("Statistics module initialized")
	return nil
}

func (s *StatsModule) Shutdown() error {
	log.Println("Statistics module shutdown")
	return nil
}

func (s *StatsModule) messageHandler(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	
	s.mu.Lock()
	s.messageCount++
	s.mu.Unlock()
}

func (s *StatsModule) readyHandler(session *discordgo.Session, r *discordgo.Ready) {
	s.mu.Lock()
	s.userCount = int64(len(r.Guilds))
	s.mu.Unlock()
}

// GetStats returns current statistics
func (s *StatsModule) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	uptime := time.Since(s.startTime)
	
	return map[string]interface{}{
		"uptime":        uptime.String(),
		"messages":      s.messageCount,
		"commands":      s.commandCount,
		"servers":       s.userCount,
		"start_time":    s.startTime.Format("2006-01-02 15:04:05"),
	}
}

// IncrementCommandCount increments the command counter
func (s *StatsModule) IncrementCommandCount() {
	s.mu.Lock()
	s.commandCount++
	s.mu.Unlock()
}
