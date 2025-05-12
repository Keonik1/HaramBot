package cmd

import (
	"github.com/bwmarrin/discordgo"

	"haram_bot/cmd/mod"
)

func GetModCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{mod.GetModCommands()}
}

func GetModHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mod": mod.ModHandler,
	}
}
