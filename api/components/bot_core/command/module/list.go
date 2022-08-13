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

package module

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/database"
)

// handleModuleList prints out a list of all commands and their status.
func handleModuleList(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	compNames := ""
	compStatus := ""

	for _, comp := range api.Components {
		if api.IsCoreComponent(comp) {
			continue
		}

		if "" != compNames {
			compNames += "\n"
		}
		compNames += comp.Name

		regComp, ok := database.GetRegisteredComponent(C, comp.Code)
		if !ok {
			continue
		}

		if "" != compStatus {
			compStatus += "\n"
		}

		globalStatus, ok := database.GetGlobalComponentStatus(C, regComp.ID)
		if !ok {
			continue
		}

		if !globalStatus.Enabled {
			compStatus += ":no_entry:"

			continue
		}

		guild, ok := database.GetGuild(C, i.GuildID)
		if !ok {
			continue
		}

		guildSpecificStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
		if !ok || !guildSpecificStatus.Enabled {
			compStatus += ":x:"

			continue
		}

		compStatus += ":white_check_mark:"
	}

	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Module Status",
				Description: "Overview of all modules and whether they are enabled or not",
				Color:       api.DefaultEmbedColor,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Module",
						Value:  compNames,
						Inline: true,
					},
					{
						Name:   "Status",
						Value:  compStatus,
						Inline: true,
					},
					{
						Name: "Legend",
						Value: ":white_check_mark: - Enabled\n" +
							":x: - Disabled\n" +
							":no_entry: - Globally disabled (Maintenance)",
						Inline: false,
					},
				},
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
