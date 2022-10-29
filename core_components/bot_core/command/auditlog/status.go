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
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
	"strconv"
)

// handleModuleDisable enables the targeted module.
func handleAuditLogStatus(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := slash_commands.GenerateInteractionResponseTemplate("Bot Audit Log", "")

	guildIdInt, err := strconv.ParseUint(i.GuildID, 10, 64)
	if nil != err {
		respondWithError(s, i, resp)

		return
	}

	guildAuditLogConfig, err := C.EntityManager().AuditLogConfig().GetByGuildId(guildIdInt)
	if nil != err {
		respondWithNotConfiguredStatus(s, i, resp)

		return
	}

	populateAuditLogConfigStatusEmbedFields(resp,
		getStatusDisplay(guildAuditLogConfig.Enabled),
		getConfiguredChannel(s, guildAuditLogConfig.ChannelId),
	)

	slash_commands.Respond(C, s, i, resp)
}

// getStatusDisplay returns the string representation
// of the passed audit log enablement status.
func getStatusDisplay(auditLogEnabled bool) string {
	if auditLogEnabled {
		return ":white_check_mark:"
	}

	return ":x:"
}

// getConfiguredChannel returns the currently configured channel
// for the bot audit log.
func getConfiguredChannel(session *discordgo.Session, channel uint64) string {
	channelIdStr := strconv.FormatUint(channel, 10)

	dgChannel, err := session.Channel(channelIdStr)
	if nil == err {
		return fmt.Sprintf("<#%v>", dgChannel)
	}

	return "Not configured!"
}

// populateComponentStatusEmbedFields cares about filling up the interaction
// response templates embed with the status of the requested component.
func populateAuditLogConfigStatusEmbedFields(
	resp *discordgo.InteractionResponseData,
	auditLogStatus string,
	auditLogChannel string,
) {
	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Enabled",
			Value:  auditLogStatus,
			Inline: false,
		},
		{
			Name:   "Configured Channel",
			Value:  auditLogChannel,
			Inline: false,
		},
	}
}

// respondWithNotConfiguredStatus responds with a message
// telling the user that the audit log is not configured.
//
// This function should be used when no database row exists.
func respondWithNotConfiguredStatus(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	populateAuditLogConfigStatusEmbedFields(resp,
		":x:",
		"Not configured!")

	slash_commands.Respond(C, s, i, resp)
}

// respondWithError responds with a message
// telling the user the command failed.
func respondWithError(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Damn, something went wrong!",
			Value: "Something unexpected happened while processing the command!",
		},
	}

	resp.Embeds[0].Fields = embeds

	slash_commands.Respond(C, s, i, resp)
}
