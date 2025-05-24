package mod

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/xhit/go-str2duration/v2"

	"haram_bot/tools"
)

var muteCommand = &discordgo.ApplicationCommandOption{
	Name:        "mute",
	Description: "mute some user",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "username",
			Description: "Username to be muted",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:         "time",
			Description:  "Mute duration (e.g. 5m, 1h)",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Name:         "description",
			Description:  "Reason for mute",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     false,
			Autocomplete: false,
		},
	},
}

func handleMute(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	duration, err := str2duration.ParseDuration(options[1].StringValue())
	if err != nil {
		tools.SendTimeParseErrorMessage(s, i, options[1].StringValue())
		return
	}

	mutedUser := options[0].UserValue(s)
	reason := tools.GetNotRequiredOptionValue(options, 2, "No description provided.")
	unmuteTime := tools.FormatTimeWithOffset(duration)

	tools.LogInfo("Mute user %s (%s) for %s", mutedUser.ID, mutedUser.GlobalName, duration) //TODO: replace this to mute logic

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "🔇 User Muted",
					Description: fmt.Sprintf("%s was muted for `%s`.\n⏰ Unmute date: `%s`\nReason: \n%s",
						mutedUser.Mention(), options[1].StringValue(), unmuteTime, reason),
					Color: 0xffa500,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Unmute",
							Style:    discordgo.DangerButton,
							CustomID: fmt.Sprintf("mod:unmute:%s", mutedUser.ID),
						},
					},
				},
			},
		},
	})
	tools.CheckInteractionError(err)
}

func handleComponentUnmute(s *discordgo.Session, i *discordgo.InteractionCreate, id string) {
	// TODO: Реализовать размут пользователя (например, снять роль мута)

	messageWithUnmuteButton := i.Message

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "🔊 User Unmuted",
					Description: fmt.Sprintf("User <@%s> has been manually unmuted by the user %s", id, i.Member.Mention()),
					Color:       0x6aa300,
				},
			},
		},
	})
	tools.CheckInteractionError(err)

	tools.DisableButtonByID(s, i, messageWithUnmuteButton, "mod:unmute:")

}

func autocompleteMute(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// TODO: research to add auto completion for typed username
	if len(options) < 2 || !options[1].Focused {
		return
	}

	choices := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "60s", Value: "60s"},
		{Name: "5m", Value: "5m"},
		{Name: "10m", Value: "10m"},
		{Name: "30m", Value: "30m"},
		{Name: "1h", Value: "1h"},
		{Name: "4h", Value: "4h"},
		{Name: "12h", Value: "12h"},
		{Name: "1d", Value: "1d"},
		{Name: "1w", Value: "1w"},
	}

	if options[1].StringValue() != "" {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  options[1].StringValue(),
			Value: options[1].StringValue(),
		})
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	tools.CheckInteractionError(err)
}

func GetMuteCommand() *discordgo.ApplicationCommandOption {
	return muteCommand
}
