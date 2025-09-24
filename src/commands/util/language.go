package util

import (
	"godiscordbot/src/commands"
	"godiscordbot/src/config"
	"godiscordbot/src/i18n"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.Command{
		Name:        "language",
		Description: "Set bot language",
		Execute: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			if len(options) == 0 {
				msg := i18n.T(config.GetLanguage(i.GuildID), "choose_language")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: msg},
				})
				return
			}

			lang := options[0].StringValue()
			if lang != "en" && lang != "pt-BR" {
				msg := i18n.T(config.GetLanguage(i.GuildID), "choose_language")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: msg},
				})
				return
			}

			config.SetLanguage(i.GuildID, lang)
			msg := i18n.T(lang, "language_set")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: msg},
			})
		},
	})
}
