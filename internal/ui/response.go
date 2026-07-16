package ui

import "github.com/bwmarrin/discordgo"

func RespondDeferred(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func EditOriginal(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) error {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	return err
}

func EditOriginalError(s *discordgo.Session, i *discordgo.InteractionCreate, title, description string) error {
	embed := NewEmbed().
		Error().
		Title(title).
		Description(description).
		Build()
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	return err
}

func RespondEmbed(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	embed *discordgo.MessageEmbed,
	isEphemeral ...bool,
) error {

	flags := discordgo.MessageFlags(0)
	if len(isEphemeral) > 0 && isEphemeral[0] {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					embed,
				},
				Flags: flags,
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
		Timestamp().
		Build()

	return RespondEmbed(s, i, embed, true)
}
