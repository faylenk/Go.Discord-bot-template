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
	commands.Register(commands.Command{
		Name:        "language",
		Description: "Set the bot's language for this server",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "language",
				Description: "Select your preferred language",
				Required:    true,
				Choices:     createLanguageChoices(),
			},
		},
		Execute: executeLanguage,
	})
}

// createLanguageChoices creates the language selection options
func createLanguageChoices() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "English", // Display name
			Value: "en",      // Internal value
		},
		{
			Name:  "Português (Brasil)", // Display name
			Value: "pt-BR",              // Internal value
		},
	}
}

// getLanguageDisplayName returns the display name for a language code
func getLanguageDisplayName(langCode string) string {
	displayNames := map[string]string{
		"en":    "English",
		"pt-BR": "Português (Brasil)",
	}

	if name, ok := displayNames[langCode]; ok {
		return name
	}
	return langCode
}

func executeLanguage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Safety check
	if len(options) == 0 {
		currentLang := config.GetLanguage(i.GuildID)
		errorMsg := fmt.Sprintf(
			"Please select a language from the dropdown menu.\n(Current language: **%s**)",
			getLanguageDisplayName(currentLang),
		)
		respondWithMessage(s, i, errorMsg)
		return
	}

	guildID := i.GuildID
	selectedLang := options[0].StringValue()
	languageName := getLanguageDisplayName(selectedLang)

	// Validate the selected language
	availableLangs := i18n.GetAvailableLocales()
	isValid := false
	for _, lang := range availableLangs {
		if selectedLang == lang {
			isValid = true
			break
		}
	}

	if !isValid {
		currentLang := config.GetLanguage(guildID)

		var availableOptions []string
		for _, langCode := range availableLangs {
			availableOptions = append(availableOptions, getLanguageDisplayName(langCode))
		}

		errorMsg := fmt.Sprintf(
			"**Invalid language selection** ❌\n(Current language: **%s**)\n\nPlease choose one of the following:\n• %s",
			getLanguageDisplayName(currentLang),
			strings.Join(availableOptions, "\n• "),
		)
		respondWithMessage(s, i, errorMsg)
		return
	}

	// Set the new language
	config.SetLanguage(guildID, selectedLang)

	// Create success message in the new language
	successMsg := fmt.Sprintf(
		"**Language Updated** ✅\n\nBot language has been set to: **%s**\n\n%s",
		languageName,
		i18n.T(selectedLang, "language_set"),
	)

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
