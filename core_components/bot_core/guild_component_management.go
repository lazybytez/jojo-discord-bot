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
)

// handleInitialComponentStatusOnGuildJoin ensures that components that are enabled by default
// are written to the database and - if newly written enabled by default.
func handleInitialComponentStatusOnGuildJoin(_ *discordgo.Session, create *discordgo.GuildCreate) {
	em := C.EntityManager()
	guild, err := em.Guilds().Get(create.ID)

	if nil != err {
		C.Logger().Warn("Tried to initialize guild component status but could not get guild \"%v\" from DB!",
			create.ID)

		return
	}

	for _, regComp := range em.RegisteredComponent().GetAvailable() {
		if !regComp.DefaultEnabled || regComp.IsCoreComponent() {
			continue
		}

		componentStatus, err := em.GuildComponentStatus().Get(guild.ID, regComp.ID)
		if nil == err {
			continue
		}

		componentStatus.Component = *regComp
		componentStatus.Guild = *guild
		componentStatus.Enabled = true

		err = em.GuildComponentStatus().Create(componentStatus)
		if nil != err {
			C.Logger().Warn("Could not enable default component \"%v\" for guild \"%v\"",
				regComp.Code,
				guild.GuildID)
		}
	}
}
