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
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/webapi"
)

var C = api.Component{
	// Metadata
	Code:         "bot_webapi",
	Name:         "Bot WebAPI",
	Categories:   api.Categories{api.CategoryInternal},
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
	compGroup := webapi.Router().Group("/components")
	compGroup.GET("/", ComponentsGet)

	commandsGroup := webapi.Router().Group("/commands")
	commandsGroup.GET("/", CommandsGet)
	commandsGroup.GET(fmt.Sprintf("/:%s", ParamCommandID), CommandGet)
	commandsGroup.GET(fmt.Sprintf("/:%s/options", ParamCommandID), CommandOptionsGet)

	return nil
}
