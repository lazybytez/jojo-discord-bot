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

// handleModuleShow prints out a list of all commands and their status.
func handleModuleShow(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Module Information",
				Color:  api.DefaultEmbedColor,
				Fields: []*discordgo.MessageEmbedField{},
			},
		},
	}

	var comp *api.Component
	for _, c := range api.Components {
		if c.Code == option.Options[0].Value {
			comp = c
			break
		}
	}

	if nil == comp || api.IsCoreComponent(comp) {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	regComp, ok := database.GetRegisteredComponent(C, comp.Code)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	globalStatus, ok := database.GetGlobalComponentStatus(C, regComp.ID)
	globalStatusOutput := ":white_check_mark:"
	if !ok {
		globalStatusOutput = ":no_entry:"
	}

	if !globalStatus.Enabled {
		globalStatusOutput += ":no_entry:"
	}

	guild, ok := database.GetGuild(C, i.GuildID)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guildSpecificStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
	guildSpecificStatusOutput := ":white_check_mark:"
	if !ok || !guildSpecificStatus.Enabled {
		guildSpecificStatusOutput += ":x:"
	}

	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Name",
			Value:  comp.Name,
			Inline: false,
		},
		{
			Name:   "Description",
			Value:  comp.Description,
			Inline: false,
		},
		{
			Name:   "Guild Status",
			Value:  guildSpecificStatusOutput,
			Inline: true,
		},
		{
			Name:   "Global Status",
			Value:  globalStatusOutput,
			Inline: true,
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
