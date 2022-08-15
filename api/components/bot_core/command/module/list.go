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
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/database"
)

// handleModuleList prints out a list of all commands and their status.
func handleModuleList(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) {
	compNamesAndStatus := generateComponentStatusTable(i)
	resp := createComponentStatusListResponse(compNamesAndStatus)

	respond(s, i, resp)
}

// createComponentStatusListResponse creates an interaction response containing
// an embed that list all components and their status.
// Additionally, a legend is added, that describes the meaning of the different states.
func createComponentStatusListResponse(compNamesAndStatus string) *discordgo.InteractionResponseData {
	resp := generateInteractionResponseDataTemplate(
		"Module Status",
		"Overview of all modules and whether they are enabled or not")

	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Status - Module",
			Value:  compNamesAndStatus,
			Inline: true,
		},
		{
			Name: "Legend",
			Value: database.GlobalComponentStatusEnabledDisplay + " - Enabled\n" +
				api.GuildComponentStatusDisabledDisplay + " - Disabled\n" +
				database.GlobalComponentStatusDisabledDisplay + " - Globally disabled (Maintenance)",
			Inline: false,
		},
	}

	return resp
}

// generateComponentStatusTable generates a string with all component names
// and a string with matching component status divided by line-breaks.
func generateComponentStatusTable(i *discordgo.InteractionCreate) string {
	compNameAndStatus := &bytes.Buffer{}

	for _, comp := range api.Components {
		if api.IsCoreComponent(comp) {
			continue
		}

		if compNameAndStatus.Len() < 1 {
			// Ignore and continue on error
			_, err := fmt.Fprint(compNameAndStatus, "\n")
			if nil != err {
				C.Logger().Warn("Failed to write linebreak while building component list entry for \"%v\"",
					comp.Name)
			}
		}

		regComp, ok := database.GetRegisteredComponent(C, comp.Code)
		if !ok {
			continue
		}

		globalStatus, ok := database.GetGlobalStatusDisplayString(C, regComp.ID)
		if !ok {
			getComponentStatusListRow(compNameAndStatus, comp.Name, globalStatus)

			continue
		}

		guild, ok := database.GetGuild(C, i.GuildID)
		if !ok {
			continue
		}

		guildSpecificStatus, _ := api.GetGuildComponentStatusDisplay(C, guild.ID, regComp.ID)
		getComponentStatusListRow(compNameAndStatus, comp.Name, guildSpecificStatus)
	}

	return compNameAndStatus.String()
}

// getComponentStatusListRow writes a single row for the module list commands
// component status list.
// If the writing fails, nothing will happen.
func getComponentStatusListRow(buf *bytes.Buffer, name string, status string) {
	_, err := fmt.Fprintf(buf, "%v - %v", status, name)
	if nil != err {
		if nil != err {
			C.Logger().Warn("Failed to write row while building component list entry for \"%v\"",
				name)
		}
	}
}
