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

package pingpong

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// This component is an example on how to create a component
// The variable C holds the api.Component itself.
// The init function initializes the component with all required
// metadata.
// To make a component autoload properly, it is necessary to:
//  1. Enable it in the api.State of the component
//  2. Add an api.Lifecycle with a valid callback
//     that initializes the component

// C is the instance of the ping pong component.
// Can be used to register the component or get information about it.
var C = api.Component{
	// Metadata
	Code:        "ping_pong",
	Name:        "Ping Pong",
	Description: "This module plays ping pong with you and returns Latency (maybe)",

	State: &api.State{
		DefaultEnabled: true,
	},
}

// init initializes the component with its metadata
func init() {
	api.RegisterComponent(&C, LoadComponent)
}

// LoadComponent loads the Ping-Pong Component
func LoadComponent(_ *discordgo.Session) error {
	_ = C.SlashCommandManager().Register(pingCommand)
	_ = C.SlashCommandManager().Register(pongCommand)

	return nil
}
