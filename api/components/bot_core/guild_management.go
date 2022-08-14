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
	"github.com/lazybytez/jojo-discord-bot/api"
	"strconv"
)

// handleGuildRegisterOnJoin is triggered when the bot joins a guild.
//
// It ensures that every guild that isn't already known is registered
// in the database. It also keeps the name of the guild updated.
func handleGuildRegisterOnJoin(_ *discordgo.Session, g *discordgo.GuildCreate) {
	guildId, err := strconv.Atoi(g.ID)
	if nil != err {
		C.Logger().Warn("Joined guild with ID \"%v\" but could not convert ID to int!", g.ID)

		return
	}

	guild, ok := api.GetGuild(C, g.ID)
	if !ok {
		guild.GuildID = guildId
		guild.Name = g.Name

		api.Create(&guild)

		return
	}

	if guild.Name != g.Name {
		api.UpdateEntity(C, &guild, api.ColumnGuildName, g.Name)
	}
}

// handleGuildUpdateOnUpdate cares about updating the stored guild name
// in the database.
func handleGuildUpdateOnUpdate(_ *discordgo.Session, g *discordgo.GuildUpdate) {
	guild, ok := api.GetGuild(C, g.ID)
	if !ok {
		C.Logger().Warn("Could not update guild with ID \"%v\" named \"%v\" as it is missing in database!",
			g.ID,
			g.Name)
		return
	}

	if guild.Name != g.Name {
		api.UpdateEntity(C, &guild, api.ColumnGuildName, g.Name)
	}
}
