package main

import (
	"godiscordbot/src/commands"
	_ "godiscordbot/src/commands/fun"
	_ "godiscordbot/src/commands/util"
	"godiscordbot/src/config" // Importe o config

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Não foi possível carregar o arquivo .env, continuando com variáveis do sistema...")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Token do bot não encontrado.")
	}

	// Inicializa o banco de dados
	config.InitDB()
	defer config.CloseDB()

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	// Handler de interação
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := commands.Registered[i.ApplicationCommandData().Name]; ok {
			cmd.Execute(s, i)
		}
	})

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()

	// Registra slash commands
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

	log.Println("Bot está online! CTRL+C para parar.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
