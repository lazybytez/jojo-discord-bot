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

package statistics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// C is the instance of the statistics component.
var C *api.Component

// init initializes the component with its metadata
func init() {
	C = &api.Component{
		// Metadata
		Name:         "Statistics Component",
		Description:  "This Component returns statistics about the bot and the runtime.",
		DmPermission: true,

		State: api.State{
			Enabled: true,
		},

		Lifecycle: api.LifecycleHooks{
			LoadComponent: LoadComponent,
		},
	}
}

// LoadComponent loads the Ping-Pong Component
func LoadComponent(discord *discordgo.Session) error {
	_ = C.SlashCommandManager().Register(statsCommand)

	return nil
}
