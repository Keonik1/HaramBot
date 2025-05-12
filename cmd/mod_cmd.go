package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/xhit/go-str2duration/v2"
)

var (
	//////////////// mod
	muteCommand = &discordgo.ApplicationCommandOption{
		Name:        "mute",
		Description: "mute some user",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "username",
				Description:  "Username what will muted",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: false,
			},
			{
				Name:         "time",
				Description:  "Mute time",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
			{
				Name:         "description",
				Description:  "description why user was muted",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     false,
				Autocomplete: false,
			},
		},
	}

	modCommand = &discordgo.ApplicationCommand{
		Name:        "mod",
		Description: "mod commands",
		Options: []*discordgo.ApplicationCommandOption{
			muteCommand,
		},
	}

	modHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData().Options[0]
			var duration, err = str2duration.ParseDuration(data.Options[1].StringValue())
			if err != nil {
				errorMessage := "ERROR: cannot parse time value." +
					"\nWrong value: `" + data.Options[1].StringValue() + "`." +
					"\nValid time units are:" +
					"\n- `s` - seconds" +
					"\n- `m` - minutes" +
					"\n- `h` - hours" +
					"\n- `d` - days" +
					"\n- `w` - weeks" +
					"\nExample: `60s` is `60 seconds`"
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprint(
							errorMessage,
						),
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					panic(err)
				}
				break
			}

			log.Printf("Mute user %s on %s",
				data.Options[0].Value,
				duration,
			)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						"User `%s` was muted on `%s`\nReason:\n%s",
						data.Options[0].StringValue(),
						data.Options[1].StringValue(),
						getNotRequiredOptionValue(data.Options, 2, "No description was provided."),
					),
				},
			})
			if err != nil {
				panic(err)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData().Options[0]
			var choices []*discordgo.ApplicationCommandOptionChoice
			switch {
			case data.Options[1].Focused:
				choices = []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "60s",
						Value: "60s",
					},
					{
						Name:  "5m",
						Value: "5m",
					},
					{
						Name:  "10m",
						Value: "10m",
					},
					{
						Name:  "30m",
						Value: "30m",
					},
					{
						Name:  "1h",
						Value: "1h",
					},
					{
						Name:  "4h",
						Value: "4h",
					},
					{
						Name:  "12h",
						Value: "12h",
					},
					{
						Name:  "1d",
						Value: "1d",
					},
					{
						Name:  "1w",
						Value: "1w",
					},
				}
				if data.Options[1].StringValue() != "" {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  data.Options[1].StringValue(),
						Value: data.Options[1].StringValue(),
					})
				}
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})
			if err != nil {
				panic(err)
			}
		}
	}
)

var (
	modCommands = []*discordgo.ApplicationCommand{
		modCommand,
	}

	modCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mod": modHandler,
	}
)

func GetModCommands() []*discordgo.ApplicationCommand {
	return modCommands
}

func GetModHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return modCommandHandlers
}
