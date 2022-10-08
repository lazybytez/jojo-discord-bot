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
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

// handleModuleShow prints out a list of all commands and their status.
func handleModuleShow(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := slash_commands.GenerateInteractionResponseTemplate("Module Information", "")

	regComp := findComponent(option)
	if nil == regComp || regComp.IsCoreComponent() {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	em := C.EntityManager()

	guild, err := em.Guilds().Get(i.GuildID)
	if nil != err {
		respondWithMissingComponent(s, i, resp, regComp.Name)

		return
	}

	globalStatusOutput, _ := em.GlobalComponentStatus().GetDisplayString(regComp.ID)
	guildSpecificStatusOutput, _ := em.GuildComponentStatus().GetDisplay(guild.ID, regComp.ID)

	populateComponentStatusEmbedFields(resp, regComp, guildSpecificStatusOutput, globalStatusOutput)

	slash_commands.Respond(C, s, i, resp)
}

// populateComponentStatusEmbedFields cares about filling up the interaction
// response templates embed with the status of the requested component.
func populateComponentStatusEmbedFields(
	resp *discordgo.InteractionResponseData,
	comp *entities.RegisteredComponent,
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
