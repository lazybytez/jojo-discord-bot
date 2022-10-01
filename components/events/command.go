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

var embedColor = 0xF9E2AF

var eventsCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "events",
		Description: "Show and manage your custom Events!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "list",
				Description: "List all the Events you signed up for!",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "manage",
				Description: "Manage (Create, Edit, Delete) the Events on Discords you have permission on!",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
	Handler: handleEventsCommand,
}

func handleEventsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	subCommand := map[string]func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
		option *discordgo.ApplicationCommandInteractionDataOption,
	){
		"list":   handleList,
		"manage": handleManage,
	}

	api.ProcessSubCommands(s, i, nil, subCommand)
}

func handleList(s *discordgo.Session, i *discordgo.InteractionCreate, _ *discordgo.ApplicationCommandInteractionDataOption) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: buildListEmbed(s),
	})
}

func handleManage(s *discordgo.Session, i *discordgo.InteractionCreate, _ *discordgo.ApplicationCommandInteractionDataOption) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: buildManageEmbed(s),
	})
}

func buildListEmbed(s *discordgo.Session) *discordgo.InteractionResponseData {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(
		":calendar_spiral: Your upcoming events",
		"Here you can see all the events you signed up for the next 14 days.\n")

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

func buildManageEmbed(s *discordgo.Session) *discordgo.InteractionResponseData {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(
		":calendar_spiral: Event Manager",
		"Do you want to Create a new Event or Edit / Delete an existing one?")

	cudButtons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "📜",
					},
					Label:    "Create",
					Style:    discordgo.PrimaryButton,
					CustomID: "event_create",
				},
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "🔧",
					},
					Label:    "Edit",
					Style:    discordgo.SecondaryButton,
					CustomID: "event_edit",
				},
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "🗑️",
					},
					Label:    "Delete",
					Style:    discordgo.DangerButton,
					CustomID: "event_delete",
				},
			},
		},
	}

	resp.Embeds[0].Color = embedColor
	resp.Components = cudButtons

	return resp
}

/*func handleEventButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {
	C.Logger().Info("Executed interaction")
}
*/
