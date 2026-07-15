package ui

import "github.com/bwmarrin/discordgo"

func RespondEmbed(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	embed *discordgo.MessageEmbed,
) error {

	return s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					embed,
				},
			},
		},
	)
}

func RespondSuccess(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title,
	description string,
) error {

	embed := NewEmbed().
		Success().
		Title(title).
		Description(description).
		Build()

	return RespondEmbed(s, i, embed)
}

func RespondInfo(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title,
	description string,
) error {

	embed := NewEmbed().
		Info().
		Title(title).
		Description(description).
		Build()

	return RespondEmbed(s, i, embed)
}

func RespondWarning(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title,
	description string,
) error {

	embed := NewEmbed().
		Warning().
		Title(title).
		Description(description).
		Build()

	return RespondEmbed(s, i, embed)
}

func RespondError(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	title,
	description string,
) error {

	embed := NewEmbed().
		Error().
		Title(title).
		Description(description).
		Build()

	return RespondEmbed(s, i, embed)
}
