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

package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"strconv"
)

// BotAuditLogger represents the service type to do bot audit logging.
type BotAuditLogger struct {
	c *Component
}

// BotAuditLogger returns the bot audit logger for the current component,
// which allows to create audit log entries.
func (c *Component) BotAuditLogger() *BotAuditLogger {
	if nil == c.botAuditLogger {
		c.botAuditLogger = &BotAuditLogger{c}
	}

	return c.botAuditLogger
}

// Log creates an audit log entry.
// The function always creates a database entry using the passed parameters.
// When announce is true and a channel has been configured for bot audit logs,
// the audit log entry will also be announced.
func (bal *BotAuditLogger) Log(guild *discordgo.Guild, user *discordgo.User, msg string, announce bool) {
	dbGuild, err := bal.c.EntityManager().Guilds().Get(guild.ID)
	if nil != err {
		bal.c.Logger().Err(err, "Tried to create bot audit log entry with message \"%s\", "+
			"but could not retrieve guild with ID \"%s\" from DB",
			msg,
			guild.ID)

		return
	}

	regComp, err := bal.c.EntityManager().RegisteredComponent().Get(bal.c.Code)
	if nil != err {
		bal.c.Logger().Err(err, "Tried to create bot audit log entry with message \"%s\", "+
			"but could not retrieve registered component with code \"%s\" from DB",
			msg,
			bal.c.Code)

		return
	}

	userIdInt, err := strconv.ParseUint(user.ID, 10, 64)
	if nil != err {
		bal.c.Logger().Err(err, "Tried to create bot audit log entry with message \"%s\", "+
			"but could not convert user id \"%s\" to uint64",
			msg,
			user.ID)

		return
	}

	auditLog := &entities.AuditLog{
		GuildID:               dbGuild.ID,
		Guild:                 *dbGuild,
		RegisteredComponentID: regComp.ID,
		RegisteredComponent:   *regComp,
		UserID:                userIdInt,
		Message:               msg,
	}

	err = bal.c.EntityManager().AuditLog().Create(auditLog)

	if announce {
		bal.announceLog(user, auditLog)
	}
}

// announceLog posts the supplied log entry on the configured bot audit log
// channel. If no bot audit log has been configured on the guild, this function
// won't do anything.
func (bal *BotAuditLogger) announceLog(user *discordgo.User, log *entities.AuditLog) {
	auditLogConfig, err := bal.c.EntityManager().AuditLogConfig().GetByGuildId(log.GuildID)
	if nil != err {
		// At this point audit log should not be configured for guild, therefore skip
		return
	}

	if !auditLogConfig.Enabled || nil == auditLogConfig.ChannelId {
		return
	}

	messagesSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Bot Audit Log",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Component",
						Value:  bal.c.Name,
						Inline: false,
					},
					{
						Name:   "User",
						Value:  user.Mention(),
						Inline: true,
					},
					{
						Name:   "Message",
						Value:  log.Message,
						Inline: false,
					},
				},
			},
		},
	}

	_, err = bal.c.discord.ChannelMessageSendComplex(
		strconv.FormatUint(*auditLogConfig.ChannelId, 10),
		messagesSend)
	if nil != err {
		bal.c.Logger().Err(err, "Tried to announce bot audit log message for guild \"%d\", "+
			"but either the configured channel \"%d\" does not exist or permissions to send messages is missing!",
			auditLogConfig.Guild.GuildID,
			*auditLogConfig.ChannelId)

		return
	}
}
