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
	"github.com/davecgh/go-spew/spew"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

// UserAction is a string that can be output
// to tell the user the action that is currently applied
// to a module
type UserAction string

const (
	UserActionEnable  = "enabled"
	UserActionDisable = "disabled"
)

// respondWithTogglingComponent responds with a message
// telling the user the command is still processing.
// This is necessary as the module enable/disable process
// can take a few seconds.
func respondWithTogglingComponent(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName string,
	action UserAction,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name: ":alarm_clock: Processing...",
			Value: spew.Sprintf("The module \"%s\" is being %s, please wait...",
				componentName,
				action),
		},
	}

	resp.Embeds[0].Fields = embeds

	slash_commands.Respond(C, s, i, resp)
}

// respondWithMissingComponent fills the passed discordgo.InteractionResponseData
// with an embed field that indicates that the specified component could not be found.
func respondWithMissingComponent(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName interface{},
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Error",
			Value: fmt.Sprintf("No module with name \"%v\" could be found!", componentName),
		},
	}

	resp.Embeds[0].Fields = embeds

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
