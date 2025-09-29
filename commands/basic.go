package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"discord-bot-forge/core"
)

// PingCommand implements a simple ping command
type PingCommand struct{}

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Description() string {
	return "Pong! Check bot latency"
}

func (c *PingCommand) Usage() string {
	return "ping"
}

func (c *PingCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	message, err := s.ChannelMessageSend(m.ChannelID, "🏓 Pong!")
	if err != nil {
		return err
	}
	
	// Edit message with latency info
	latency := s.HeartbeatLatency()
	content := fmt.Sprintf("🏓 Pong! Latency: %v", latency)
	
	_, err = s.ChannelMessageEdit(m.ChannelID, message.ID, content)
	return err
}

func (c *PingCommand) Permissions() []string {
	return []string{}
}

func (c *PingCommand) Cooldown() int {
	return 0
}

func (c *PingCommand) Category() string {
	return "General"
}

// HelpCommand shows available commands
type HelpCommand struct {
	bot *core.Bot
}

func NewHelpCommand(bot *core.Bot) *HelpCommand {
	return &HelpCommand{bot: bot}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Show available commands"
}

func (c *HelpCommand) Usage() string {
	return "help [command]"
}

func (c *HelpCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) > 0 {
		// Show help for specific command
		cmdName := args[0]
		if cmd, exists := c.bot.Commands[cmdName]; exists {
			embed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Command: %s", cmd.Name()),
				Description: cmd.Description(),
				Color:       0x00ff00,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Usage",
						Value:  fmt.Sprintf("`%s%s`", c.bot.Config.Prefix, cmd.Usage()),
						Inline: false,
					},
					{
						Name:   "Category",
						Value:  cmd.Category(),
						Inline: true,
					},
					{
						Name:   "Cooldown",
						Value:  fmt.Sprintf("%d seconds", cmd.Cooldown()),
						Inline: true,
					},
				},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		} else {
			s.ChannelMessageSend(m.ChannelID, "❌ Command not found.")
		}
	} else {
		// Show all commands grouped by category
		categories := c.bot.GetCommandCategories()
		
		var fields []*discordgo.MessageEmbedField
		for category, commands := range categories {
			var commandList strings.Builder
			for _, cmd := range commands {
				commandList.WriteString(fmt.Sprintf("**%s** - %s\n", cmd.Name(), cmd.Description()))
			}
			
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   category,
				Value:  commandList.String(),
				Inline: false,
			})
		}
		
		embed := &discordgo.MessageEmbed{
			Title:       "🔥 DiscordBotForge Commands",
			Description: fmt.Sprintf("Use `%shelp <command>` for detailed information", c.bot.Config.Prefix),
			Color:       0xff6b35,
			Fields:      fields,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("DiscordBotForge v%s", c.bot.Version),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
	return nil
}

func (c *HelpCommand) Permissions() []string {
	return []string{}
}

func (c *HelpCommand) Cooldown() int {
	return 0
}

func (c *HelpCommand) Category() string {
	return "General"
}

// InfoCommand shows bot information
type InfoCommand struct {
	bot *core.Bot
}

func NewInfoCommand(bot *core.Bot) *InfoCommand {
	return &InfoCommand{bot: bot}
}

func (c *InfoCommand) Name() string {
	return "info"
}

func (c *InfoCommand) Description() string {
	return "Show DiscordBotForge information"
}

func (c *InfoCommand) Usage() string {
	return "info"
}

func (c *InfoCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	embed := &discordgo.MessageEmbed{
		Title:       "🔥 DiscordBotForge",
		Description: "A modular framework for forging Discord bots",
		Color:       0xff6b35,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Version",
				Value:  c.bot.Version,
				Inline: true,
			},
			{
				Name:   "Commands",
				Value:  fmt.Sprintf("%d", len(c.bot.Commands)),
				Inline: true,
			},
			{
				Name:   "Modules",
				Value:  fmt.Sprintf("%d", len(c.bot.Modules)),
				Inline: true,
			},
			{
				Name:   "Middleware",
				Value:  fmt.Sprintf("%d", len(c.bot.Middleware)),
				Inline: true,
			},
			{
				Name:   "Prefix",
				Value:  c.bot.Config.Prefix,
				Inline: true,
			},
			{
				Name:   "Debug Mode",
				Value:  fmt.Sprintf("%t", c.bot.Config.DebugMode),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Built with Go and discordgo",
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return nil
}

func (c *InfoCommand) Permissions() []string {
	return []string{}
}

func (c *InfoCommand) Cooldown() int {
	return 0
}

func (c *InfoCommand) Category() string {
	return "General"
}
