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
	enableCommandResponseHeader                            = "Enable Bot Audit Log"
	enableCannotConfigureWithoutChannelResponseName        = ":x: Whoops, no bot audit log without a channel!"
	enableCannotConfigureWithoutChannelResponseValue       = "To enable the bot audit log for the first time, you must enter a valid channel to write the logs to!"
	enableAlreadyConfiguredForChannelResponseName          = ":x: Nothing to do here!"
	enableAlreadyConfiguredForChannelResponseValueTemplate = "The bot audit log is already enabled and configured to use the channel %s!"
	enableSuccessResponseName                              = ":white_check_mark: Done!"
	enableSuccessResponseValueTemplate                     = "The bot audit log is now enabled and configured to use the channel %s!"
)

func handleAuditLogEnable(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	options *discordgo.ApplicationCommandInteractionDataOption,
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

	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate(enableCommandResponseHeader, "")

	guild, err := C.EntityManager().Guilds().Get(i.GuildID)
	if nil != err {
		slash_commands.RespondWithGenericErrorMessage(C, s, i, resp)

		return
	}

	var channel *discordgo.Channel
	for _, option := range options.Options {
		if option.Name == "channel" {
			channel = option.ChannelValue(s)

			break
		}
	}

	guildAuditLogConfig, err := C.EntityManager().AuditLogConfig().GetByGuildId(guild.ID)
	if nil != err && nil == channel {
		slash_commands.RespondWithSimpleEmbedMessage(C,
			s,
			i,
			resp,
			enableCannotConfigureWithoutChannelResponseName,
			enableCannotConfigureWithoutChannelResponseValue)

		return
	}
	if nil != err {
		// Prepare new audit log config entity
		guildAuditLogConfig.GuildID = guild.ID
		guildAuditLogConfig.Guild = *guild
	}

	newChannelIdInt := uint64(0)
	if nil != channel {
		newChannelIdInt, err = strconv.ParseUint(channel.ID, 10, 64)
		if nil != err {
			slash_commands.RespondWithGenericErrorMessage(C, s, i, resp)

			return
		}
	}

	if guildAuditLogConfig.ChannelId != nil && *guildAuditLogConfig.ChannelId == newChannelIdInt && guildAuditLogConfig.Enabled {
		slash_commands.RespondWithSimpleEmbedMessage(C,
			s,
			i,
			resp,
			enableAlreadyConfiguredForChannelResponseName,
			fmt.Sprintf(enableAlreadyConfiguredForChannelResponseValueTemplate, channel.Mention()))

		return
	}

	if nil != channel {
		guildAuditLogConfig.ChannelId = &newChannelIdInt
	}

	if nil == guildAuditLogConfig.ChannelId && nil == channel {
		slash_commands.RespondWithSimpleEmbedMessage(C,
			s,
			i,
			resp,
			enableCannotConfigureWithoutChannelResponseName,
			enableCannotConfigureWithoutChannelResponseValue)

		return
	}

	guildAuditLogConfig.Enabled = true

	err = C.EntityManager().AuditLogConfig().Save(guildAuditLogConfig)
	if nil != err {
		slash_commands.RespondWithGenericErrorMessage(C, s, i, resp)

		return
	}

	notifyAuditLogChannelConfigured(s, channel, i.Member)
	slash_commands.RespondWithSimpleEmbedMessage(C,
		s,
		i,
		resp,
		enableSuccessResponseName,
		fmt.Sprintf(enableSuccessResponseValueTemplate, channel.Mention()))

	C.BotAuditLogger().Log(
		dgoGuild,
		user,
		fmt.Sprintf("The bot audit log announcements have been enabled for channel %s!", channel.Mention()),
		false)
}

// notifyAuditLogChannelConfigured sends an information message to the channel that
// has been configured as the new audit log channel.
func notifyAuditLogChannelConfigured(session *discordgo.Session, channel *discordgo.Channel, member *discordgo.Member) {
	if nil == member {
		C.Logger().Warn("Failed to notify users on guild \"%s\" about bot audit log configuration "+
			"with channel \"%s\" as member of interaction is nil!",
			channel.GuildID,
			channel.ID)
	}

	_, err := session.ChannelMessageSend(channel.ID, fmt.Sprintf(":white_check_mark: %s configured this "+
		"channel to receive future bot audit log messages.\n\n"+
		"If this channel was choosen by accident and you want to reconfigure the bot audit log, use the "+
		"command `/jojo auditlog enable <channel>` to do so!\n"+
		"If this action was unintended, use the command `/jojo auditlog disable` to disable the bot audit log.",
		member.Mention()))

	if nil != err {
		C.Logger().Warn("Failed to notify users on guild \"%s\" about bot audit log configuration "+
			"with channel \"%s\" triggered by user \"%s\"!",
			channel.GuildID,
			channel.ID,
			member.User.ID)
	}
}
