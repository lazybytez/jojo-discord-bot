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
	"github.com/lazybytez/jojo-discord-bot/api/database"
)

// handleInitialComponentStatusOnGuildJoin ensures that components that are enabled by default
// are written to the database and - if newly written enabled by default.
func handleInitialComponentStatusOnGuildJoin(_ *discordgo.Session, create *discordgo.GuildCreate) {
	guild, ok := database.GetGuild(C, create.ID)

	if !ok {
		C.Logger().Warn("Tried to initialize guild component status but could not get guild \"%v\" from DB!",
			create.ID)

		return
	}

	for _, comp := range api.Components {
		if !comp.State.DefaultEnabled || api.IsCoreComponent(comp) {
			continue
		}

		regComp, ok := database.GetRegisteredComponent(C, comp.Code)

		if !ok {
			C.Logger().Warn("Tried to get registered component \"%v\" from DB but failed!",
				comp.Code)

			return
		}

		componentStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
		if ok {
			continue
		}

		componentStatus.Component = *regComp
		componentStatus.Guild = *guild
		componentStatus.Enabled = true

		database.Create(componentStatus)
	}
}
