package conf

import (
	"github.com/bwmarrin/discordgo"
)

var confCommands = &discordgo.ApplicationCommand{
	Name:        "conf",
	Description: "Configuration commands",
	Options: []*discordgo.ApplicationCommandOption{
		GetNSFWChannelsCommand(), // Add other subcommands here
		// GetMuteRolesCommand(),
	},
}

func ConfHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData().Options[0]
		switch data.Name {
		case "nsfw-channels":
			handleNSFWChannels(s, i, data.Options)
			// case "mute-roles":
			// 	handleMuteRoles(s, i, data.Options)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData().Options[0]
		switch data.Name {
		// case "nsfw-channels":
		// 	autocompleteNSFWChannels(s, i, data.Options)
		// case "mute-roles":
		// 	autocompleteMuteRoles(s, i, data.Options)
		}
	}
}

func GetConfCommands() *discordgo.ApplicationCommand {
	return confCommands
}
