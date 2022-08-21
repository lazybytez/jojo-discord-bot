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

package bot_core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/database"
	"strconv"
)

// handleGuildRegisterOnJoin is triggered when the bot joins a guild.
//
// It ensures that every guild that isn't already known is registered
// in the database. It also keeps the name of the guild updated.
func handleGuildRegisterOnJoin(_ *discordgo.Session, g *discordgo.GuildCreate) {
	guildId, err := strconv.ParseUint(g.ID, 10, 64)
	if nil != err {
		C.Logger().Warn("Joined guild with ID \"%v\" but could not convert ID to int!", g.ID)

		return
	}
	em := C.EntityManager()

	var guild *database.Guild
	guild, err = em.Guilds().Get(g.ID)
	if nil != err {
		guild.GuildID = guildId
		guild.Name = g.Name

		err = em.Create(&guild)
		if nil != err {
			C.Logger().Warn("Failed to create guild with ID \"%v\" in database!", g.ID)
		}

		return
	}

	if guild.Name != g.Name {
		err = em.UpdateEntity(&guild, database.ColumnName, g.Name)
		if nil != err {
			C.Logger().Warn("Failed to update guild with ID \"%v\" in database!", g.ID)
		}
	}
}

// handleCommandSyncOnGuildJoin ensures that the appropriate set of slash commands
// is registered for the guild.
func handleCommandSyncOnGuildJoin(session *discordgo.Session, g *discordgo.GuildCreate) {
	C.SlashCommandManager().SyncApplicationComponentCommands(session, g.ID)
}

// handleGuildUpdateOnUpdate cares about updating the stored guild name
// in the database.
func handleGuildUpdateOnUpdate(_ *discordgo.Session, g *discordgo.GuildUpdate) {
	em := C.EntityManager()
	guild, err := em.Guilds().Get(g.ID)
	if err != nil {
		C.Logger().Warn("Could not update guild with ID \"%v\" named \"%v\" as it is missing in database!",
			g.ID,
			g.Name)
		return
	}

	updateGuildOnNameChange(em, guild, g)
}

// updateGuildOnNameChange updates the name of the passed Guild
// in the database, if the name changed.
func updateGuildOnNameChange(em *database.EntityManager, guild *database.Guild, g *discordgo.GuildUpdate) {
	if guild.Name != g.Name {
		err := em.UpdateEntity(&guild, database.ColumnName, g.Name)
		if nil != err {
			C.Logger().Warn("Failed to update guild with ID \"%v\" in database!", g.ID)
		}
	}
}
