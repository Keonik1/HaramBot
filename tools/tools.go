package tools

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetNotRequiredOptionValue(options []*discordgo.ApplicationCommandInteractionDataOption, index int, fallback string) string {
	if len(options) > index {
		return options[index].StringValue()
	}
	return fallback
}

func CheckInteractionError(err error) {
	if err != nil {
		LogInfo("Error responding to interaction: %v", err)
	}
}

func SendTimeParseErrorMessage(s *discordgo.Session, i *discordgo.InteractionCreate, value string) {
	errorMessage := "ERROR: cannot parse time value." +
		"\nWrong value: `" + value + "`." +
		"\nValid time units are:" +
		"\n- `s` - seconds" +
		"\n- `m` - minutes" +
		"\n- `h` - hours" +
		"\n- `d` - days" +
		"\n- `w` - weeks" +
		"\nExample: `60s` is `60 seconds`"
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprint(errorMessage),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	CheckInteractionError(err)
}

func RegisterHandlers(s *discordgo.Session, handlerMaps ...map[string]func(*discordgo.Session, *discordgo.InteractionCreate)) {
	combined := make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate))
	for _, hm := range handlerMaps {
		for k, v := range hm {
			if _, exists := combined[k]; exists {
				LogInfo("WARNING: handler for command %q is being overwritten", k)
			}
			combined[k] = v
		}
	}

	LogInfo("Register handlers...")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := combined[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
