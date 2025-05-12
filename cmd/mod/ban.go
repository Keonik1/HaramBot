package mod

import (
	"fmt"
	"haram_bot/tools"
	"log"

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
			Type:        discordgo.ApplicationCommandOptionString,
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
	username := options[0].StringValue()
	reason := tools.GetNotRequiredOptionValue(options, 1, "No reason provided.")

	log.Printf("Banned user %s. Reason: %s", username, reason)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("User `%s` has been banned.\nReason: %s", username, reason),
		},
	})
	tools.CheckInteractionError(err)
}

func GetBanCommand() *discordgo.ApplicationCommandOption {
	return banCommand
}
