package cmd

import (
	"github.com/bwmarrin/discordgo"

	"haram_bot/cmd/help"
	"haram_bot/cmd/mod"
)

const (
	CmdMod  string = "mod"
	CmdHelp string = "help"
)

type CommandModule struct {
	Command func() *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commandModules = map[string]CommandModule{
	CmdMod: {
		Command: mod.GetModCommands,
		Handler: mod.ModHandler,
	},
	CmdHelp: {
		Command: help.GetHelpCommands,
		Handler: help.HelpHandler,
	},
}

func GetCommandsByName(name string) []*discordgo.ApplicationCommand {
	if module, ok := commandModules[name]; ok {
		return []*discordgo.ApplicationCommand{module.Command()}
	}
	return nil
}

func GetHandlersByName(name string) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if module, ok := commandModules[name]; ok {
		return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			name: module.Handler,
		}
	}
	return nil
}

func GetModCommands() []*discordgo.ApplicationCommand {
	return GetCommandsByName(CmdMod)
}

func GetModHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return GetHandlersByName(CmdMod)
}

func GetHelpCommands() []*discordgo.ApplicationCommand {
	return GetCommandsByName(CmdHelp)
}

func GetHelpHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return GetHandlersByName(CmdHelp)
}
