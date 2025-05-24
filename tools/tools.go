package tools

import (
	"fmt"
	"strings"
	"time"

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
				LogWarning("handler for command %q is being overwritten", k)
			}
			combined[k] = v
		}
	}

	LogInfo("Register handlers...")

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			name := i.ApplicationCommandData().Name
			if h, ok := combined[name]; ok {
				h(s, i)
			} else {
				LogWarning("no handler found for command: %q", name)
			}
		case discordgo.InteractionMessageComponent:
			customID := i.MessageComponentData().CustomID
			prefix := strings.SplitN(customID, ":", 2)[0]
			if h, ok := combined[prefix]; ok {
				h(s, i)
			} else {
				LogWarning("no handler found for component: %q", customID)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			name := i.ApplicationCommandData().Name
			if h, ok := combined[name]; ok {
				h(s, i)
			} else {
				LogWarning("no handler found for autocomplete command: %q", name)
			}
		default:
			LogWarning("unsupported interaction type: %d", i.Type)
		}
	})
}

func FormatTimeWithOffset(d time.Duration) string {
	return time.Now().Add(d).Format("2006-01-02 15:04:05 MST")
}

func DisableButtonByID(s *discordgo.Session, i *discordgo.InteractionCreate, message *discordgo.Message, targetCustomIDPrefix string) {
	if message == nil {
		if i.Message == nil {
			LogError("Message is nil!")
			return
		}
		message = i.Message
	}

	var newComponents []discordgo.MessageComponent

	for _, row := range message.Components {
		actionRow, ok := row.(*discordgo.ActionsRow)
		if ok {
			var newRow discordgo.ActionsRow
			for _, comp := range actionRow.Components {
				if btn, ok := comp.(*discordgo.Button); ok {
					newBtn := btn
					if strings.HasPrefix(btn.CustomID, targetCustomIDPrefix) {
						newBtn.Disabled = true
					}
					newRow.Components = append(newRow.Components, newBtn)
				}
			}
			newComponents = append(newComponents, newRow)
		}
	}

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.ChannelID,
		ID:         message.ID,
		Components: &newComponents,
	})
	CheckInteractionError(err)
}
