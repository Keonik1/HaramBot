package mod

import (
	"fmt"
	"haram_bot/tools"

	"github.com/bwmarrin/discordgo"
)

var banCommand = &discordgo.ApplicationCommandOption{
	Name:        "ban",
	Description: "Ban a user",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "username",
			Description: "User to ban",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:        "reason",
			Description: "Reason for ban",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    false,
		},
	},
}

func handleBan(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	reason := tools.GetNotRequiredOptionValue(options, 1, "No reason provided.")

	bannedUser := options[0].UserValue(s)
	tools.LogInfo("Banned user %s (%s). Reason: %s", bannedUser.ID, bannedUser.GlobalName, reason)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("User %s has been banned.\nReason: %s", bannedUser.Mention(), reason),
		},
	})
	tools.CheckInteractionError(err)
}

func GetBanCommand() *discordgo.ApplicationCommandOption {
	return banCommand
}
