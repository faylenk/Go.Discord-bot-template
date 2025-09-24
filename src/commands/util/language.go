package util

import (
	"fmt"
	"godiscordbot/src/commands"
	"godiscordbot/src/config"
	"godiscordbot/src/i18n"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	langChoices := createLanguageChoices()

	commands.Register(commands.Command{
		Name:        "language",
		Description: "Set the bot's language for this server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "lang",
				Description: "The language you want to set.",
				Required:    true,
				Choices:     langChoices,
			},
		},
		Execute: executeLanguage,
	})
}

// createLanguageChoices creates a list of language choices for the command option.
func createLanguageChoices() []*discordgo.ApplicationCommandOptionChoice {
	availableLangs := i18n.GetAvailableLocales()
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(availableLangs))

	// Maps language codes to their display names
	langNames := map[string]string{
		"en":    "English",
		"pt-BR": "Português (Brasil)",
		// Add more languages here as needed
	}

	for _, langCode := range availableLangs {
		name, ok := langNames[langCode]
		if !ok {
			name = langCode
		}
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: langCode,
		})
	}
	return choices
}

func executeLanguage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	guildID := i.GuildID
	selectedLang := options[0].StringValue()

	currentLang := config.GetLanguage(guildID)
	availableLangs := i18n.GetAvailableLocales()
	isValid := false
	for _, lang := range availableLangs {
		if selectedLang == lang {
			isValid = true
			break
		}
	}

	if !isValid {
		errorMsg := fmt.Sprintf(
			"%s\n%s: %s",
			i18n.T(currentLang, "invalid_language"),
			i18n.T(currentLang, "available_languages"),
			strings.Join(availableLangs, ", "),
		)
		respondWithMessage(s, i, errorMsg)
		return
	}

	config.SetLanguage(guildID, selectedLang)

	successMsg := i18n.T(selectedLang, "language_set")
	respondWithMessage(s, i, successMsg)
}

func respondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}
