package util

import (
	"fmt"
	"godiscordbot/src/commands"
	"godiscordbot/src/config"
	"godiscordbot/src/i18n"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.Command{
		Name:        "ping",
		Description: "Checks the bot's latency.",
		Execute:     executePing,
	})
}

func executePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lang := config.GetLanguage(i.GuildID)

	latency := s.HeartbeatLatency().Round(time.Millisecond)

	start := time.Now()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral, // 🔹 deixa a resposta ephemeral
		},
	})

	if err != nil {
		log.Printf("Error responding to ping command: %v", err)
		return
	}

	apiPing := time.Since(start).Round(time.Millisecond)

	responseText := fmt.Sprintf(
		"%s\n🏓 **Heartbeat**: %s\n⚡ **API Response**: %s",
		i18n.T(lang, "ping_response"),
		latency,
		apiPing,
	)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &responseText,
	})

	if err != nil {
		log.Printf("Error editing ping response: %v", err)
	}
}
