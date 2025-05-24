package mod

import (
	"haram_bot/tools"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var modCommands = &discordgo.ApplicationCommand{
	Name:        "mod",
	Description: "Moderation commands",
	Options: []*discordgo.ApplicationCommandOption{
		GetMuteCommand(), // Add other subcommands here
		GetBanCommand(),
	},
}

func ModHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData().Options[0]
		switch data.Name {
		case "mute":
			handleMute(s, i, data.Options)
		case "ban":
			handleBan(s, i, data.Options)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData().Options[0]
		switch data.Name {
		case "mute":
			autocompleteMute(s, i, data.Options)
		}
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		tools.LogTrace("%v", data.CustomID)
		parts := strings.Split(data.CustomID, ":")
		if len(parts) < 3 {
			return
		}
		action := parts[1]
		userID := parts[2]
		switch action {
		case "unmute":
			tools.LogTrace("unmute handler")
			handleComponentUnmute(s, i, userID)
		}
	}
}

func GetModCommands() *discordgo.ApplicationCommand {
	return modCommands
}
