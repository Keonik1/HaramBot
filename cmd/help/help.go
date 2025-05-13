package help

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"haram_bot/tools"
)

var helpCommand = &discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Print help",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "topic",
			Description:  "Print extended help about specified topic",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     false,
			Autocomplete: true,
		},
	},
}

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate, helpType string) {
	var helpMessage string
	switch helpType {
	case "common":
		helpMessage = "some common help"
	case "mod":
		helpMessage = "some mod help"
	case "conf":
		helpMessage = "some conf help"
	case "cmd":
		helpMessage = "some cmd help"
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprint(helpMessage),
		},
	})
	tools.CheckInteractionError(err)
}

func autocompleteHelp(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	choices := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "common", Value: "common"},
		{Name: "mod", Value: "mod"},
		{Name: "conf", Value: "conf"},
		{Name: "cmd", Value: "cmd"},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	tools.CheckInteractionError(err)
}

func HelpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		helpType := tools.GetNotRequiredOptionValue(data.Options, 0, "common")
		handleHelp(s, i, helpType)
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		data := i.ApplicationCommandData()
		autocompleteHelp(s, i, data.Options)
	}
}

func GetHelpCommands() *discordgo.ApplicationCommand {
	return helpCommand
}
