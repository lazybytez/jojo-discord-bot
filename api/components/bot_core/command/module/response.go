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
)

// generateShowInteractionResponseTemplate creates a prefilled discordgo.InteractionResponseData
// that is prepared to be filled with specific data or errors.
//
// The template will get an empty embed with the specified name and description.
func generateInteractionResponseDataTemplate(name string, description string) *discordgo.InteractionResponseData {
	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       name,
				Description: description,
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}
	return resp
}

// respond to the passed interaction with the passed
// discordgo.InteractionResponseData as a message in the channel
// where the interaction has been triggered.
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, resp *discordgo.InteractionResponseData) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})

	if nil != err {
		C.Logger().Err(err, "Failed to deliver interaction response on slash-command!")
	}
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
