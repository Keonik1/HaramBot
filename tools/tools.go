package tools

import (
	"fmt"
	"log"

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
		log.Println("Error responding to interaction:", err)
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
