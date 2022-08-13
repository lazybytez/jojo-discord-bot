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

var C *api.Component

func init() {
	C = &api.Component{
		// Metadata
		Code:        "bot_core",
		Name:        "Bot Core",
		Description: "This component handles core routines and database management.",

		State: &api.State{
			DefaultEnabled: true,
		},

		Lifecycle: api.LifecycleHooks{
			LoadComponent: LoadComponent,
		},
	}
}

// LoadComponent loads the bot core component
// and handles migration of core entities
// and registration of important core event handlers.
func LoadComponent(discord *discordgo.Session) error {
	prepareDatabase()
	initializeComponentManagement()

	_, _ = C.HandlerManager().Register("register_guilds", handleGuildRegisterOnJoin)
	_, _ = C.HandlerManager().Register("update_registered_guilds", handleGuildUpdateOnUpdate)
	_, _ = C.HandlerManager().Register(
		"populate_default_guild_component_status",
		handleInitialComponentStatusOnGuildJoin)

	_ = C.SlashCommandManager().Register(jojoCommand)

	return nil
}

// prepareDatabase updates the schema with the core entities which lay
// in the database package.
func prepareDatabase() {
	// Guild related entities
	_ = database.RegisterEntity(C, &database.Guild{})

	// Component related entities
	_ = database.RegisterEntity(C, &database.RegisteredComponent{})
	_ = database.RegisterEntity(C, &database.ComponentStatus{})
	_ = database.RegisterEntity(C, &database.GlobalComponentStatus{})
}

// initializeComponentManagement initializes the component management
// by populating the database with necessary data and pre-warming the cache.
func initializeComponentManagement() {
	registerAvailableComponents()
	ensureGlobalComponentStatusExists()
}
