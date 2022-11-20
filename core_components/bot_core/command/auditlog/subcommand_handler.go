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

package auditlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

var C *api.Component

// HandleAuditLogCommandSubCommand handles the execution of the
// "auditlog" subcommand.
//
// The command allows to manage the bot audit log configuration.
func HandleAuditLogCommandSubCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	if nil == i.Member {
		slash_commands.RespondWithCommandIsGuildOnly(C, s, i, "auditlog")

		return
	}

	subCommands := map[string]func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
		option *discordgo.ApplicationCommandInteractionDataOption,
	){
		"status":  handleAuditLogStatus,
		"enable":  handleAuditLogEnable,
		"disable": handleAuditLogDisable,
	}

	success := api.ProcessSubCommands(
		s,
		i,
		option,
		subCommands)

	if !success {
		if !success {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The executed (sub)command is invalid or does not exist!",
				},
			})
		}
	}
}
