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

package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

var messageActionEventCreate = &api.MessageAction{
	CustomID: "event_create",
	Handler:  handleEventCreation,
}

var messageActionEventCreateModal = &api.MessageAction{
	CustomID: "event_create_data_category",
	Handler:  handleEventCreationModal,
}

var messageActionEventCreateModalAfter = &api.MessageAction{
	CustomID: "event_creation_modal",
	Handler:  handleEventCreationModalAfter,
}

var messageActionEventEdit = &api.MessageAction{
	CustomID: "event_edit",
	Handler:  handleEventEditing,
}

var messageActionEventDelete = &api.MessageAction{
	CustomID: "event_delete",
	Handler:  handleEventDeletion,
}

func handleEventCreation(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: buildCreationCategorySelectionEmbed(s, i),
	})
}

func handleEventCreationModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: buildCreationModal(s, i),
	})
}

func handleEventCreationModalAfter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Create Channel with Event
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: buildCreationCategorySelectionEmbed(s, i),
	})
}

func handleEventEditing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: buildEditingEmbed(),
	})
}

func handleEventDeletion(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: buildDeletionEmbed(),
	})
}

func buildCreationCategorySelectionEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(
		":calendar_spiral: Create an Event",
		"Choose the category to create you Event in :point_down:")

	cudButtons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					// Select menu, as other components, must have a customID, so we set it to this value.
					CustomID:    "event_create_data_category",
					Placeholder: "Categories",
					Options:     buildSelectMenuChannelOptions(s, i),
				},
			},
		},
	}

	resp.Embeds[0].Color = embedColor
	resp.Components = cudButtons

	return resp
}

func buildSelectMenuChannelOptions(s *discordgo.Session, i *discordgo.InteractionCreate) []discordgo.SelectMenuOption {
	var resp []discordgo.SelectMenuOption

	channels, _ := s.GuildChannels(i.GuildID)

	for _, c := range channels {
		// Only category channel get processed
		if c.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		categoryChannelOption := discordgo.SelectMenuOption{
			Label: c.Name,
			Value: c.ID,
		}
		resp = append(resp, categoryChannelOption)
	}

	return resp
}

func buildCreationModal(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	resp := &discordgo.InteractionResponseData{
		CustomID: "event_creation_modal",
		Title:    "Create an Event",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "event_create_data_name",
						Label:       "Event Name",
						Style:       discordgo.TextInputShort,
						Placeholder: "A short Event Name",
						Required:    true,
						MaxLength:   24,
						MinLength:   1,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "event_create_data_description",
						Label:       "What is this event about?",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "An optional description about your event",
						Required:    false,
						MaxLength:   500,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "event_create_data_date",
						Label:       "When does the Event start? (DD/MM/YYYY HH:mm)",
						Style:       discordgo.TextInputShort,
						Placeholder: "DD/MM/YYYY HH:mm",
						Required:    true,
						MaxLength:   16,
						MinLength:   16,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "event_create_data_group_number",
						Label:       "How many groups participate? ('1' - '99')",
						Style:       discordgo.TextInputShort,
						Placeholder: "Enter a number from '1' - '99'",
						Value:       "1",
						Required:    false,
						MaxLength:   2,
						MinLength:   1,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "event_create_data_group_size",
						Label:       "How big are the groups? ('1' - '99')",
						Style:       discordgo.TextInputShort,
						Placeholder: "Enter a number from '1' - '99'",
						Value:       "10",
						Required:    false,
						MaxLength:   2,
						MinLength:   1,
					},
				},
			},
		},
	}

	return resp
}

func buildEditingEmbed() *discordgo.InteractionResponseData {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(
		":calendar_spiral: Which Event do you want to edit? :wrench:",
		"Here are all the Events listed you can edit.\n")

	eventsField := []*discordgo.MessageEmbedField{
		{
			Name: "Events",
			Value: "- <#1025840920404435045>: <t:1664634300:F>\n" +
				"- <#1025840920404435045>: <t:1664634300:F>\n" +
				"- <#1025840920404435045>: <t:1664634300:F>\n" +
				"- <#1025840920404435045>: <t:1664634300:F>\n",
		},
	}

	resp.Embeds[0].Fields = eventsField
	resp.Embeds[0].Color = embedColor

	return resp
}

func buildDeletionEmbed() *discordgo.InteractionResponseData {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(
		":calendar_spiral: Which Event do you want to delete? :wastebasket:",
		"Here are all the Events listed you can delete.\n")

	eventsField := []*discordgo.MessageEmbedField{
		{
			Name: "Events",
			Value: "- <#1025840920404435045>: <t:1664634300:F>\n" +
				"- <#1025840920404435045>: <t:1664634300:F>\n" +
				"- <#1025840920404435045>: <t:1664634300:F>\n" +
				"- <#1025840920404435045>: <t:1664634300:F>\n",
		},
	}

	resp.Embeds[0].Fields = eventsField
	resp.Embeds[0].Color = embedColor

	return resp
}
