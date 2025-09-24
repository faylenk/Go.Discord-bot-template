package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Execute     func(s *discordgo.Session, i *discordgo.InteractionCreate)
} // Type that defines a command structure.

var Registered = make(map[string]*Command)

func Register(cmd Command) {
	Registered[cmd.Name] = &cmd
}
