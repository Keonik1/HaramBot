package cmd

import "github.com/bwmarrin/discordgo"

func getNotRequiredOptionValue(options []*discordgo.ApplicationCommandInteractionDataOption, index int, fallback string) string {
	if len(options) > index {
		return options[index].StringValue()
	}
	return fallback
}
