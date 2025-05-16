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
			Name:         "username",
			Description:  "Username to be muted",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: false,
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

	tools.LogInfo("Mute user %s for %s", options[0].Value, duration) //TODO: replace this to mute logic

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("User `%s` was muted for `%s`\nReason:\n%s",
				options[0].StringValue(),
				options[1].StringValue(),
				tools.GetNotRequiredOptionValue(options, 2, "No description provided."),
			),
		},
	})
	tools.CheckInteractionError(err)
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
