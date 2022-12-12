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

package sync_commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"time"
)

var C *api.Component

// getLastGuildSyncCacheKey returns the cache key used to store
// the last time a command sync has been run.
func getLastGuildSyncCacheKey(guildID string) string {
	return fmt.Sprintf("last_command_sync_%s", guildID)
}

// HandleSyncCommandSubCommand handles the execution of the
// "sync_commands" subcommand.
//
// The command allows to trigger re-sync of the commands registered for
// in case a bot administrator thinks there is an inconsistency.
func HandleSyncCommandSubCommand(
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

	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate("Slash Command Synchronisation", "")

	cacheKey := getLastGuildSyncCacheKey(i.GuildID)
	lastSync, ok := cache.Get(cacheKey, time.Time{})

	if ok && time.Since(lastSync) < 10*time.Minute {
		respondWithOnCoolDown(s, i, resp)

		return
	}

	respondWithProcessing(s, i, resp)
	C.BotAuditLogger().Log(
		dgoGuild,
		user,
		"A slash-command re-sync has been triggered",
		true)

	C.Logger().Info(
		"Manual slash-command sync has been triggered for guild \"%v\"",
		i.GuildID)
	C.SlashCommandManager().SyncApplicationComponentCommands(s, i.GuildID)

	currentTime := time.Now()
	cache.Update(cacheKey, currentTime)

	finishWitSuccess(s, i, resp)
	C.BotAuditLogger().Log(
		dgoGuild,
		user,
		"A slash-command re-sync has been finished",
		true)
}

// respondWithProcessing responds with a message
// telling the user the command is still processing.
func respondWithProcessing(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":alarm_clock: Processing...",
			Value: "Synchronisation is in progress and can take up to a minute, please wait...",
		},
	}

	resp.Embeds[0].Fields = embeds

	slash_commands.Respond(C, s, i, resp)
}

// respondWithOnCoolDown responds with a message
// telling the user the command is still on cool-down.
func respondWithOnCoolDown(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Too fast!",
			Value: "This command can only be used once every 10 minutes!",
		},
	}

	resp.Embeds[0].Fields = embeds

	slash_commands.Respond(C, s, i, resp)
}

// respondWithOnCoolDown responds with a message
// telling the user the command is still on cool-down.
func finishWitSuccess(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":white_check_mark: Done!",
			Value: "Commands have been sucessfully synced with your guild!",
		},
	}

	// We have an already appropriate embed formed, therefore just
	// edit and use it.
	resp.Embeds[0].Fields = embeds

	slash_commands.EditResponse(C, s, i, &discordgo.WebhookEdit{
		Embeds: &resp.Embeds,
	})
}
