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
)

// handleModuleShow prints out a list of all commands and their status.
func handleModuleShow(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := generateInteractionResponseDataTemplate("Module Information", "")

	comp := findComponent(option)
	if nil == comp || api.IsCoreComponent(comp) {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	regComp, ok := api.GetRegisteredComponent(C, comp.Code)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guild, ok := api.GetGuild(C, i.GuildID)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	globalStatusOutput, _ := api.GetGlobalStatusDisplayString(C, regComp.ID)
	guildSpecificStatusOutput, _ := api.GetGuildComponentStatusDisplay(C, guild.ID, regComp.ID)

	populateComponentStatusEmbedFields(resp, comp, guildSpecificStatusOutput, globalStatusOutput)

	respond(s, i, resp)
}

// populateComponentStatusEmbedFields cares about filling up the interaction
// response templates embed with the status of the requested component.
func populateComponentStatusEmbedFields(
	resp *discordgo.InteractionResponseData,
	comp *api.Component,
	guildSpecificStatusOutput string,
	globalStatusOutput string,
) {
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
}
