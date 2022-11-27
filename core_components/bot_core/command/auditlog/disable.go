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
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

const (
	disableCommandResponseHeader                      = "Disable Bot Audit Log"
	disableAuditLogNeverConfiguredBeforeResponseName  = ":x: Oh no, no configuration could be found!"
	disableAuditLogNeverConfiguredBeforeResponseValue = "The bot audit log is already disabled, as it has not been configured before!"
	disableAuditLogAlreadyDisabledResponseName        = ":x: Nothing to do here!"
	disableAuditLogAlreadyDisabledResponseValue       = "The bot audit log is already disabled!"
	disableAuditSuccessResponseName                   = ":white_check_mark: Done!"
	disableAuditLogSuccessResponseValue               = "The bot audit log is now disabled! You can enable it again at any time!"
)

func handleAuditLogDisable(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	_ *discordgo.ApplicationCommandInteractionDataOption,
) {
	dgoGuild, err := s.Guild(i.GuildID)
	if nil != err {
		C.Logger().Err(err, "Failed to get guild with id \"%S\" to create "+
			"bot audit log when enabling a module on guild!",
			i.GuildID)

		return
	}

	user := i.User
	if nil == user {
		user = i.Member.User
	}

	resp := slash_commands.GenerateInteractionResponseTemplate(disableCommandResponseHeader, "")

	guild, err := C.EntityManager().Guilds().Get(i.GuildID)
	if nil != err {
		slash_commands.RespondWithGenericErrorMessage(C, s, i, resp)

		return
	}

	guildAuditLogConfig, err := C.EntityManager().AuditLogConfig().GetByGuildId(guild.ID)
	if nil != err {
		slash_commands.RespondWithSimpleEmbedMessage(C,
			s,
			i,
			resp,
			disableAuditLogNeverConfiguredBeforeResponseName,
			disableAuditLogNeverConfiguredBeforeResponseValue)

		return
	}

	if !guildAuditLogConfig.Enabled {
		slash_commands.RespondWithSimpleEmbedMessage(C,
			s,
			i,
			resp,
			disableAuditLogAlreadyDisabledResponseName,
			disableAuditLogAlreadyDisabledResponseValue)

		return
	}

	guildAuditLogConfig.Enabled = false
	guildAuditLogConfig.ChannelId = nil

	err = C.EntityManager().AuditLogConfig().Save(guildAuditLogConfig)
	if nil != err {
		slash_commands.RespondWithGenericErrorMessage(C, s, i, resp)

		return
	}

	slash_commands.RespondWithSimpleEmbedMessage(C,
		s,
		i,
		resp,
		disableAuditSuccessResponseName,
		disableAuditLogSuccessResponseValue)

	C.BotAuditLogger().Log(
		dgoGuild,
		user,
		"The bot audit log announcements have been disabled!",
		false)
}
