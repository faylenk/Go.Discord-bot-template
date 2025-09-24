package util

import (
	"godiscordbot/src/commands"
	"godiscordbot/src/config"
	"godiscordbot/src/i18n"
	"log"

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
	responseText := i18n.T(lang, "ping_response")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseText,
		},
	})

	if err != nil {
		log.Printf("Error responding to ping command: %v", err)
	}
}
