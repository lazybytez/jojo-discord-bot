/*
 * JOJO Discord Bot - An advanced multi-purpose discord bot
 * Copyright (C) 2022 Lazy Bytez (Elias Knodel, Pascal Zarrad)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package slash_commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

const (
	GenericErrorResponseEmbedName  = ":x: Damn, something went wrong!"
	GenericErrorResponseEmbedValue = "Something unexpected happened while processing the command!"
)

// GenerateInteractionResponseTemplate creates a prefilled discordgo.InteractionResponseData
// that is prepared to be filled with specific data or errors.
//
// The template will get an empty embed with the specified name and description.
// The default color of the auto-generated embed equals the one defined under api.DefaultEmbedColor
func GenerateInteractionResponseTemplate(
	name string,
	description string,
) *discordgo.InteractionResponseData {
	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       name,
				Description: description,
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}
	return resp
}

// GenerateEphemeralInteractionResponseTemplate creates a prefilled discordgo.InteractionResponseData
// that is prepared to be filled with specific data or errors.
// In addition, the discordgo.InteractionResponseData will be ephemeral and therefore only
// show to the user that triggered the interaction.
//
// The template will get an empty embed with the specified name and description.
func GenerateEphemeralInteractionResponseTemplate(
	name string,
	description string,
) *discordgo.InteractionResponseData {
	resp := &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       name,
				Description: description,
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}
	return resp
}

// Respond to the target interaction with the passed
// discordgo.InteractionResponseData as a message in the channel
// where the interaction has been triggered.
func Respond(
	c *api.Component,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})

	if nil != err {
		c.Logger().Err(err, "Failed to deliver interaction response on slash-command!")
	}
}

// RespondWithSimpleEmbedMessage fills the passed discordgo.InteractionResponseData
// with an embed that contains a single field with a name and value.
//
// The prepared interaction response will be sent.
func RespondWithSimpleEmbedMessage(
	c *api.Component,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	header string,
	message string,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  header,
			Value: message,
		},
	}

	resp.Embeds[0].Fields = embeds

	Respond(c, s, i, resp)
}

// RespondWithGenericErrorMessage fills the passed discordgo.InteractionResponseData
// with a generic error message as content.
//
// The prepared interaction response will be sent.
func RespondWithGenericErrorMessage(
	c *api.Component,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  GenericErrorResponseEmbedName,
			Value: GenericErrorResponseEmbedValue,
		},
	}

	resp.Embeds[0].Fields = embeds

	Respond(c, s, i, resp)
}

// RespondWithCommandIsGuildOnly responds with a message stating that the executed
// command is only available on guilds.
func RespondWithCommandIsGuildOnly(
	c *api.Component,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	commandName string,
) {
	resp := GenerateInteractionResponseTemplate("Stop!",
		fmt.Sprintf("The `%s` command and its subcommand cannot be executed in DMs!", commandName))

	Respond(c, s, i, resp)
}

// EditResponse edits the passed interactions original response with
// the passed discordgo.WebhooKParams.
// This allows things like sending a notice that something will take some time
// and then edit the message to tell the user that the action has been done.
func EditResponse(
	c *api.Component,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	editData *discordgo.WebhookEdit,
) *discordgo.Message {
	message, err := s.InteractionResponseEdit(i.Interaction, editData)

	if nil != err {
		c.Logger().Err(err, "Failed to deliver interaction response on slash-command!")
	}

	return message
}
