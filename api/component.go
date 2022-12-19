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

package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/services"
	"strings"
)

// CoreComponentPrefix is the prefix put in front of components that
// cannot be managed by server owners, as they are important core components
const CoreComponentPrefix = "bot_"

// State holds the state of a component.
// This includes states like:
//   - is the component enabled?
//   - is the component currently loaded?
type State struct {
	// Loaded is used to determine if the component
	// has been loaded properly or not.
	Loaded bool

	// DefaultEnabled is the default status to set for the component
	// in the database when the bot joins a new guild.
	DefaultEnabled bool
}

// Component is the base structure that must be initialized
// to create a new component.
//
// It holds basic metadata about the component
type Component struct {
	// Metadata
	Code         entities.ComponentCode
	Name         string
	Description  string
	LoadPriority int

	// State
	State *State

	// Lifecycle hooks
	loadComponentFunction func(_ *discordgo.Session) error

	// Utilities
	// These are private and only managed by the API system.
	// Their initialization happens through call to the methods
	// used to get them (Example: logger -> Component.Logger()).
	// See ServiceManager interface on how to obtain them.
	discordApi          DiscordApiWrapper
	logger              services.Logger
	handlerManager      ComponentHandlerManager
	slashCommandManager *SlashCommandManager
	discord             *discordgo.Session
	botAuditLogger      *BotAuditLogger
}

// RegistrableComponent is the interface that allows a component to be
// initialized and registered.
type RegistrableComponent interface {
	LoadComponent(discord *discordgo.Session) error
	UnloadComponent(discord *discordgo.Session) error
}

// ServiceManager is a simple interface that defines the methods
// that provide the APIs features, like Logger.
//
// Although it is most of the time not best practice, the API returns
// interfaces. We decided to use this design as interfaces are great to offer
// an API where internal things or complete subsystems can be swapped
// without breaking components. They act as contracts in this application.
// This allows us to maintain a separate logger implementation in the `services/logger`
// that is compatible with this API. So, when using the API package, think of using
// contracts.
//
// tl;dr: everything from the services package is low level and should not be used
// directly when possible (although it is not prohibited). The api package is the consumer of
// the services package and implements interfaces. These interfaces are provided to components.
// Because of that, we decided to return interfaces instead of structs in the API.
type ServiceManager interface {
	// Logger is used to obtain the Logger of a component
	//
	// On first call, this function initializes the private Component.logger
	// field. On consecutive calls, the already present Logger will be used.
	Logger() services.Logger
	// HandlerManager returns the management interface for event handlers.
	//
	// It allows the registration, decoration and general
	// management of event handlers.
	//
	// It should be always used when event handlers to listen  for
	// Discord events are necessary. It natively handles stuff like logging
	// event Handler status.
	HandlerManager() ComponentHandlerManager
	// SlashCommandManager is used to obtain the components slash Command management
	//
	// On first call, this function initializes the private Component.slashCommandManager
	// field. On consecutive calls, the already present CommonSlashCommandManager will be used.
	SlashCommandManager() CommonSlashCommandManager
	// DatabaseAccess returns the currently active DatabaseAccess.
	// The currently active DatabaseAccess is shared across all components.
	//
	// The DatabaseAccess allows to interact with the entities of the application.
	// Prefer using the EntityManager instead, as DatabaseAccess is considered
	// a low-level api.
	DatabaseAccess() *EntityManager
	// BotAuditLogger returns the bot audit logger for the current component,
	// which allows to create audit log entries.
	BotAuditLogger() *BotAuditLogger
	// DiscordApi is used to obtain the components slash DiscordApiWrapper management
	//
	// On first call, this function initializes the private Component.discordAPi
	// field. On consecutive calls, the already present DiscordGoApiWrapper will be used.
	DiscordApi() DiscordApiWrapper
	// BotStatusManager returns the current StatusManager which
	// allows to add additional status to the bot.
	BotStatusManager() StatusManager
}

// LoadComponent is used by the component registration system that
// automatically calls the LoadComponent method for all Component instances in
// the components.Components array.
func (c *Component) LoadComponent(discord *discordgo.Session) error {
	c.discord = discord

	err := c.loadComponentFunction(discord)

	if err != nil {
		return err
	}
	c.State.Loaded = true

	return nil
}

// UnloadComponent is used by the component registration system that
// automatically calls the UnregisterComponent method for all Component instances in
// the components.Components array.
//
// The function takes care of tasks like unregistering slash-commands and so on.
//
// It is used to give components the ability to gracefully shutdown.
func (c *Component) UnloadComponent(*discordgo.Session) error {
	c.HandlerManager().UnregisterAll()

	c.State.Loaded = false

	return nil
}

// IsCoreComponent checks whether the passed component is a core
// component or not.
//
// Core components are components which are prefixed with the CoreComponentPrefix.
func IsCoreComponent(c *Component) bool {
	return strings.HasPrefix(string(c.Code), CoreComponentPrefix)
}
