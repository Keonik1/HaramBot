package mod

import (
	"github.com/bwmarrin/discordgo"
)

var modCommands = &discordgo.ApplicationCommand{
	Name:        "mod",
	Description: "Moderation commands",
	Options: []*discordgo.ApplicationCommandOption{
		GetMuteCommand(), // Add other subcommands here
	},
}

func ModHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData().Options[0]
		switch data.Name {
		case "mute":
			handleMute(s, i, data.Options)
		}
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		data := i.ApplicationCommandData().Options[0]
		switch data.Name {
		case "mute":
			autocompleteMute(s, i, data.Options)
		}
	}
}

func GetModCommands() *discordgo.ApplicationCommand {
	return modCommands
}
