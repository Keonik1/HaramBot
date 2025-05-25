package conf

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"haram_bot/tools"
)

var nsfwChannelsCommand = &discordgo.ApplicationCommandOption{
	Name:        "nsfw-channels",
	Description: "Haram Bot NSFW channels management",
	Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "get",
			Description: "Get list of NSFW Channels",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Required:    false,
		},
		{
			Name:        "add",
			Description: "Add channel to list of NSFW Channels",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel",
					Description: "Channel that will be added to NSFW channel list",
					Type:        discordgo.ApplicationCommandOptionChannel,
					Required:    true,
				},
			},
		},
		{
			Name:        "del",
			Description: "Delete channel from NSFW channels",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel",
					Description: "Channel that will be deleted from NSFW channel list",
					Type:        discordgo.ApplicationCommandOptionChannel, //replace tp string and add own autocompletion (get from db current list)
					Required:    true,
				},
			},
		},
	},
}

func handleNSFWChannels(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	command := options[0].Name

	switch command {
	case "get":
		sendNSFWChannelsList(s, i)
	case "add":
		addChannelToNSFWChannelsList(s, i, options[0].Options)
	case "del":
		deleteChannelFromNSFWChannelsList(s, i, options[0].Options)
	default:
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Unknown subcommand '%s'", command),
			},
		})
		tools.CheckInteractionError(err)
	}
}

func getNSFWChannelsList(serverID string) []string {
	channelsIDs := []string{
		"1373275236782313512",
	} //replace to sql request from DB
	return channelsIDs
}

func sendNSFWChannelsList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelsIDs := getNSFWChannelsList(i.GuildID)

	message := "List of current NSFW Channels:"
	for _, channelID := range channelsIDs {
		message = fmt.Sprintf("%s\n- <#%s>", message, channelID)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	tools.CheckInteractionError(err)
}

func addChannelToNSFWChannelsList(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	channelID := options[0].Value

	// add to DB
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Channel <#%s> was added to NSFW Channels", channelID),
		},
	})
	tools.CheckInteractionError(err)
}

func deleteChannelFromNSFWChannelsList(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	channelID := options[0].Value

	// delete from DB
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Channel <#%s> was deleted from NSFW Channels", channelID),
		},
	})
	tools.CheckInteractionError(err)
}

func GetNSFWChannelsCommand() *discordgo.ApplicationCommandOption {
	return nsfwChannelsCommand
}
