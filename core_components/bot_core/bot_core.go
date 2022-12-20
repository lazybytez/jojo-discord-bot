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
)

var C = api.Component{
	// Metadata
	Code:         "bot_core",
	Name:         "Bot Core",
	Categories:   api.Categories{api.CategoryInternal},
	Description:  "This component handles core routines and entity management.",
	LoadPriority: 1000,

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
	initializeComponentManagement()

	_, _ = C.HandlerManager().Register("guild_join", onGuildJoin)
	_, _ = C.HandlerManager().Register("update_registered_guilds", handleGuildUpdateOnUpdate)

	// We need to handle the JOJO command special as it needs access to the component list.
	// This is only possible after the API has been properly initialized and the components.Components
	// list has been accessed once.
	//
	// Therefore, we configure and register the command when this core component is
	// loaded, as at this point the API should know the components too.
	initAndRegisterJojoCommand()

	return nil
}

// initializeComponentManagement initializes the component management
// by populating the database with necessary data and pre-warming the cache.
func initializeComponentManagement() {
	ensureGlobalComponentStatusExists()
}

// onGuildJoin is an event handler called when to bot joins a guild.
// It ensures the guild registration and global status registrations
// are happening in the right order.
func onGuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	handleGuildRegisterOnJoin(s, g)
	handleInitialComponentStatusOnGuildJoin(s, g)
	handleCommandSyncOnGuildJoin(s, g)
}
