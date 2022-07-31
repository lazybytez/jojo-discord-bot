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
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package components

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/log"
	"github.com/lazybytez/jojo-discord-bot/components/pingpong"
)

// LogComponentRegistry is the custom component name used
// to identify log messages from the component management system
const LogComponentRegistry = "Component Manager"

// Components contains all components that should be available.
//
// Enabled components should be registered here.
var Components = []*api.Component{
	pingpong.C,
}

// RegisterComponents handles the initialization of
// all components listed in the Components array.
//
// When it is not possible to register a component,
// an error will be printed into the log.
// The application will continue to run as nothing happened.
func RegisterComponents(discord *discordgo.Session) {
	log.Info(LogComponentRegistry, "Starting component load sequence...")
	for _, comp := range Components {
		if nil == comp.Lifecycle.LoadComponent {
			log.Debug(LogComponentRegistry, "Component \"%v\" does not have an load callback, not loading it!", comp.Name)
			continue
		}

		if !comp.State.Enabled {
			log.Info(LogComponentRegistry, "Component \"%v\" is not enabled, skipping!", comp.Name)
			continue
		}

		log.Info(LogComponentRegistry, "Loading component \"%v\"...", comp.Name)
		err := comp.RegisterComponent(discord)
		if nil != err {
			log.Warn(LogComponentRegistry, "Failed to load component with name \"%v\": %v", comp.Name, err.Error())
			continue
		}
		log.Info(LogComponentRegistry, "Successfully loaded component \"%v\"!", comp.Name)
	}
	log.Info(LogComponentRegistry, "Component load sequence completed!")
}

// UnloadComponents iterates through all registered api.Component
// registered in the Components array and calls their UnloadComponent
// function.
//
// If an api.Component does not have an UnloadComponent function defined,
// it will be ignored.
func UnloadComponents(discord *discordgo.Session) {
	log.Info(LogComponentRegistry, "Starting component unload sequence...")
	for _, comp := range Components {
		if nil == comp.Lifecycle.UnloadComponent {
			log.Debug(LogComponentRegistry, "Component \"%v\" does not have an unload callback, skipping!", comp.Name)
			continue
		}

		if !comp.State.Loaded {
			log.Warn(LogComponentRegistry, "Component \"%v\" has not been loaded, skipping!", comp.Name)
			continue
		}

		log.Info(LogComponentRegistry, "Unloading component \"%v\"...", comp.Name)
		err := comp.UnregisterComponent(discord)
		if nil != err {
			log.Warn(LogComponentRegistry, "Failed to unload component with name \"%v\": %v", comp.Name, err.Error())
			continue
		}
		log.Info(LogComponentRegistry, "Successfully unloaded component \"%v\"!", comp.Name)
	}
	log.Info(LogComponentRegistry, "Unload sequence completed!")
}
