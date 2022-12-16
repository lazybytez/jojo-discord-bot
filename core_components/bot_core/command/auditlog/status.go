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

const (
	statusCommandResponseHeader                = "Bot Audit Log Status"
	statusBotAuditLogEnabledName               = "Enabled"
	statusBotAuditLogEnabledValue              = ":white_check_mark:"
	statusBotAuditLogDisabledValue             = ":x:"
	statusBotAuditLogChannelName               = "Channel"
	statusBotAuditLogChannelNotConfiguredValue = "Not configured!"
)

// handleModuleDisable enables the targeted module.
func handleAuditLogStatus(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(statusCommandResponseHeader, "")

	guild, err := C.EntityManager().Guilds().Get(i.GuildID)
	if nil != err {
		slash_commands.RespondWithGenericErrorMessage(C, s, i, resp)

		return
	}

	guildAuditLogConfig, err := C.EntityManager().AuditLogConfig().GetByGuildId(guild.ID)
	if nil != err {
		respondWithAuditLogConfigStatus(s,
			i,
			resp,
			statusBotAuditLogDisabledValue,
			statusBotAuditLogChannelNotConfiguredValue)

		return
	}

	currentChannelId := guildAuditLogConfig.ChannelId
	if nil == currentChannelId {
		respondWithAuditLogConfigStatus(s,
			i,
			resp,
			statusBotAuditLogDisabledValue,
			statusBotAuditLogChannelNotConfiguredValue)

		return
	}

	channelStatus, channelFound := getConfiguredChannel(s, *currentChannelId)
	statusBadge := getStatusDisplay(guildAuditLogConfig.Enabled, channelFound)

	respondWithAuditLogConfigStatus(s,
		i,
		resp,
		statusBadge,
		channelStatus)
}

// getStatusDisplay returns the string representation
// of the passed audit log enablement status.
func getStatusDisplay(auditLogEnabled bool, channelFound bool) string {
	if auditLogEnabled && channelFound {
		return statusBotAuditLogEnabledValue
	}

	return statusBotAuditLogDisabledValue
}

// getConfiguredChannel returns the currently configured channel
// for the bot audit log. It also returns a boolean indicating whether the configured channel could
// be found on the guild.
func getConfiguredChannel(session *discordgo.Session, channel uint64) (string, bool) {
	channelIdStr := strconv.FormatUint(channel, 10)

	dgChannel, err := session.Channel(channelIdStr)
	if nil == err {
		return fmt.Sprintf("<#%v>", dgChannel.ID), true
	}

	return fmt.Sprintf("Seems like the configured channel with id `%d` has been deleted. "+
		"Please re-enable the bot auditlog to receive further bot audit log messages!",
		channel), false
}

// respondWithAuditLogConfigStatus cares about filling up the interaction
// response templates embed with the status of the bot audit log.
// After preparing the response, the response is sent.
func respondWithAuditLogConfigStatus(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	auditLogStatus string,
	auditLogChannel string,
) {
	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   statusBotAuditLogEnabledName,
			Value:  auditLogStatus,
			Inline: false,
		},
		{
			Name:   statusBotAuditLogChannelName,
			Value:  auditLogChannel,
			Inline: false,
		},
	}

	slash_commands.Respond(C, s, i, resp)
}
