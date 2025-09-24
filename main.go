package main

import (
	"godiscordbot/src/commands"
	_ "godiscordbot/src/commands/fun"
	_ "godiscordbot/src/commands/util"
	"godiscordbot/src/config"
	"godiscordbot/src/i18n"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Loads .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("It was not possible to load .env file, proceeding with system environment variables.")
	}

	// Loads Translations
	err = i18n.LoadLocales()
	if err != nil {
		log.Printf("Error loading locales: %v", err)
	} else {
		log.Println("Translations loaded successfully")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("No bot token provided in DISCORD_TOKEN environment variable.")
	}

	// Initialize the database
	config.InitDB()
	defer config.CloseDB()

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	// Interaction handler
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Ensure the guild exists in the database
		if i.GuildID != "" {
			config.EnsureGuildExists(i.GuildID)
		}

		if cmd, ok := commands.Registered[i.ApplicationCommandData().Name]; ok {
			cmd.Execute(s, i)
		}
	})

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	// Register commands
	var discordCmds []*discordgo.ApplicationCommand
	for _, cmd := range commands.Registered {
		discordCmds = append(discordCmds, &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		})
	}
	_, err = bot.ApplicationCommandBulkOverwrite(bot.State.User.ID, "", discordCmds)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is now running.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
