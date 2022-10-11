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

package bot_webapi

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/webapi"
)

var C = api.Component{
	// Metadata
	Code:         "bot_webapi",
	Name:         "Bot WebAPI",
	Description:  "This component handles setup of the web api for the bots core api endpoints.",
	LoadPriority: 999,

	State: &api.State{
		DefaultEnabled: true,
	},
}

func init() {
	api.RegisterComponent(&C, LoadComponent)
}

// LoadComponent loads the bot core component
// and handles migration of core entities
// and registration of important core event handlers.
func LoadComponent(_ *discordgo.Session) error {
	eg := webapi.Router().Group("/components")
	eg.GET("/", ComponentsGet)

	return nil
}
