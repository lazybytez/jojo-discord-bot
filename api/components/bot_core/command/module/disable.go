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
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/database"
)

// handleModuleDisable enables the targeted module.
func handleModuleDisable(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := generateInteractionResponseDataTemplate("Disable Module", "")

	comp := findComponent(option)
	if nil == comp || api.IsCoreComponent(comp) {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	regComp, ok := database.GetRegisteredComponent(C, comp.Code)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guild, ok := database.GetGuild(C, i.GuildID)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	if !disableComponentForGuild(guild, regComp) {
		respondWithAlreadyDisabled(s, i, resp, comp.Name)

		return
	}

	generateModuleDisableSuccessfulEmbedField(resp, comp)
	respond(s, i, resp)
}

func disableComponentForGuild(
	guild *database.Guild,
	regComp *database.RegisteredComponent,
) bool {
	guildSpecificStatus, ok := api.GetGuildComponentStatus(C, guild.ID, regComp.ID)
	if !ok {
		// No database entry = disabled
		return false
	}

	if !guildSpecificStatus.Enabled {
		return false
	}

	guildSpecificStatus.Enabled = false
	database.Save(guildSpecificStatus)

	return true
}

// generateModuleDisableSuccessfulEmbedField creates the necessary embed fields
// used to response to a successful module disable command.
func generateModuleDisableSuccessfulEmbedField(resp *discordgo.InteractionResponseData, comp *api.Component) {
	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Module",
			Value:  comp.Name,
			Inline: false,
		},
		{
			Name:   "Status",
			Value:  ":x: - The module has been disabled!",
			Inline: false,
		},
	}
}

// respondWithAlreadyDisabled fills the passed discordgo.InteractionResponseData
// with an embed field that indicates that the specified component is already disabled.
func respondWithAlreadyDisabled(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName string,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Error",
			Value: fmt.Sprintf("Module with name \"%v\" is already disabled!", componentName),
		},
	}

	resp.Embeds[0].Fields = embeds

	respond(s, i, resp)
}
