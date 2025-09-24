package util

import (
	"godiscordbot/src/commands"
	"godiscordbot/src/config"
	"godiscordbot/src/i18n"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.Command{
		Name:        "ping",
		Description: "Ping command",
		Execute: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			lang := config.GetLanguage(i.GuildID)
			msg := i18n.T(lang, "ping_response")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: msg},
			})
		},
	})
}
