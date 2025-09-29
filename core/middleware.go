package core

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// CooldownMiddleware implements rate limiting for commands
type CooldownMiddleware struct {
	cooldowns map[string]map[string]time.Time
	duration  time.Duration
}

// NewCooldownMiddleware creates a new cooldown middleware
func NewCooldownMiddleware(duration time.Duration) *CooldownMiddleware {
	return &CooldownMiddleware{
		cooldowns: make(map[string]map[string]time.Time),
		duration:  duration,
	}
}

// Name returns the middleware name
func (c *CooldownMiddleware) Name() string {
	return "Cooldown"
}

// Process implements the Middleware interface
func (c *CooldownMiddleware) Process(s *discordgo.Session, m *discordgo.MessageCreate, next func()) error {
	userID := m.Author.ID
	channelID := m.ChannelID
	
	// Initialize user cooldowns if not exists
	if c.cooldowns[userID] == nil {
		c.cooldowns[userID] = make(map[string]time.Time)
	}
	
	// Check if user is on cooldown for this channel
	if lastUsed, exists := c.cooldowns[userID][channelID]; exists {
		if time.Since(lastUsed) < c.duration {
			remaining := c.duration - time.Since(lastUsed)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⏰ Please wait %.1f seconds before using another command.", remaining.Seconds()))
			return nil
		}
	}
	
	// Update cooldown
	c.cooldowns[userID][channelID] = time.Now()
	
	// Execute next middleware/command
	next()
	return nil
}

// PermissionMiddleware checks if user has required permissions
type PermissionMiddleware struct {
	requiredPermissions []string
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(permissions []string) *PermissionMiddleware {
	return &PermissionMiddleware{
		requiredPermissions: permissions,
	}
}

// Name returns the middleware name
func (p *PermissionMiddleware) Name() string {
	return "Permission"
}

// Process implements the Middleware interface
func (p *PermissionMiddleware) Process(s *discordgo.Session, m *discordgo.MessageCreate, next func()) error {
	// Get user permissions
	permissions, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return fmt.Errorf("error getting user permissions: %w", err)
	}
	
	// Check if user has required permissions
	for _, perm := range p.requiredPermissions {
		if !hasPermission(permissions, perm) {
			s.ChannelMessageSend(m.ChannelID, "❌ You don't have permission to use this command.")
			return nil
		}
	}
	
	next()
	return nil
}

// hasPermission checks if user has a specific permission
func hasPermission(permissions int64, permission string) bool {
	switch permission {
	case "ADMINISTRATOR":
		return permissions&discordgo.PermissionAdministrator != 0
	case "MANAGE_MESSAGES":
		return permissions&discordgo.PermissionManageMessages != 0
	case "MANAGE_CHANNELS":
		return permissions&discordgo.PermissionManageChannels != 0
	case "MANAGE_ROLES":
		return permissions&discordgo.PermissionManageRoles != 0
	case "KICK_MEMBERS":
		return permissions&discordgo.PermissionKickMembers != 0
	case "BAN_MEMBERS":
		return permissions&discordgo.PermissionBanMembers != 0
	default:
		return false
	}
}

// LoggingMiddleware logs command usage
type LoggingMiddleware struct{}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// Name returns the middleware name
func (l *LoggingMiddleware) Name() string {
	return "Logging"
}

// Process implements the Middleware interface
func (l *LoggingMiddleware) Process(s *discordgo.Session, m *discordgo.MessageCreate, next func()) error {
	log.Printf("Command executed by %s#%s in channel %s: %s", 
		m.Author.Username, m.Author.Discriminator, m.ChannelID, m.Content)
	next()
	return nil
}

// OwnerOnlyMiddleware restricts commands to bot owner only
type OwnerOnlyMiddleware struct {
	ownerID string
}

// NewOwnerOnlyMiddleware creates a new owner-only middleware
func NewOwnerOnlyMiddleware(ownerID string) *OwnerOnlyMiddleware {
	return &OwnerOnlyMiddleware{
		ownerID: ownerID,
	}
}

// Name returns the middleware name
func (o *OwnerOnlyMiddleware) Name() string {
	return "OwnerOnly"
}

// Process implements the Middleware interface
func (o *OwnerOnlyMiddleware) Process(s *discordgo.Session, m *discordgo.MessageCreate, next func()) error {
	if m.Author.ID != o.ownerID {
		s.ChannelMessageSend(m.ChannelID, "❌ This command is restricted to the bot owner.")
		return nil
	}
	
	next()
	return nil
}
