package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Execute     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Registered = make(map[string]*Command)

func Register(cmd Command) {
	Registered[cmd.Name] = &cmd
}
