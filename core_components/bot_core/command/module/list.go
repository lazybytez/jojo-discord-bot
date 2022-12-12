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
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

// handleModuleList prints out a list of all commands and their status.
func handleModuleList(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) {
	compNamesAndStatus := generateComponentStatusTable(i)
	resp := createComponentStatusListResponse(compNamesAndStatus)

	slash_commands.Respond(C, s, i, resp)
}

// createComponentStatusListResponse creates an interaction response containing
// an embed that list all components and their status.
// Additionally, a legend is added, that describes the meaning of the different states.
func createComponentStatusListResponse(compNamesAndStatus string) *discordgo.InteractionResponseData {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(
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
			Value: entities.GlobalComponentStatusEnabledDisplay + " - Enabled\n" +
				entities.GuildComponentStatusDisabledDisplay + " - Disabled\n" +
				entities.GlobalComponentStatusDisabledDisplay + " - Globally disabled (Maintenance)",
			Inline: false,
		},
	}

	return resp
}

// generateComponentStatusTable generates a string with all component names
// and a string with matching component status divided by line-breaks.
func generateComponentStatusTable(i *discordgo.InteractionCreate) string {
	em := C.EntityManager()

	compNameAndStatus := &bytes.Buffer{}

	for _, regComp := range em.RegisteredComponent().GetAvailable() {
		if regComp.IsCoreComponent() {
			continue
		}

		if compNameAndStatus.Len() < 1 {
			// Ignore and continue on error
			_, err := fmt.Fprint(compNameAndStatus, "\n")
			if nil != err {
				C.Logger().Warn("Failed to write linebreak while building component list entry for \"%v\"",
					regComp.Name)
			}
		}

		globalStatus, err := em.GlobalComponentStatus().GetDisplayString(regComp.ID)
		if nil != err {
			getComponentStatusListRow(compNameAndStatus, regComp.Name, globalStatus)

			continue
		}

		guild, err := em.Guilds().Get(i.GuildID)
		if nil != err {
			continue
		}

		guildSpecificStatus, _ := em.GuildComponentStatus().GetDisplay(guild.ID, regComp.ID)
		getComponentStatusListRow(compNameAndStatus, regComp.Name, guildSpecificStatus)
	}

	return compNameAndStatus.String()
}

// getComponentStatusListRow writes a single row for the module list commands
// component status list.
// If the writing fails, nothing will happen.
func getComponentStatusListRow(buf *bytes.Buffer, name string, status string) {
	_, err := fmt.Fprintf(buf, "%v - %v\n", status, name)
	if nil != err {
		if nil != err {
			C.Logger().Warn("Failed to write row while building component list entry for \"%v\"",
				name)
		}
	}
}
