package fun

import (
	"godiscordbot/src/commands"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Register(commands.Command{
		Name:        "joke",
		Description: "Conte uma piada",
		Execute: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "kkkkkkkkkkkkkkkkkkkkkkkkkk",
				},
			})
		},
	})
}
