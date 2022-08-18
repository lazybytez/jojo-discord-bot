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

package internal

import "C"
import (
	"github.com/bwmarrin/discordgo"
	apiDatabase "github.com/lazybytez/jojo-discord-bot/api/database"
	"github.com/lazybytez/jojo-discord-bot/api/log"
	"github.com/lazybytez/jojo-discord-bot/components"
)

// logComponentRegistry is the custom component name used
// to identify log messages from the component management system
const logComponentRegistry = "component_manager"

// RegisterComponents registers all available components in the database
// and fills the available components in the database API, to provide
// a unified API to get component information.
func RegisterComponents() {
	log.Info(logComponentRegistry, "Registering components in database...")
	em := apiDatabase.GetEntityManager()
	for _, component := range components.Components {
		registeredComponent, err := em.RegisteredComponent().Get(component.Code)

		if nil != err {
			registeredComponent.Code = component.Code
			registeredComponent.Name = component.Name
			registeredComponent.Description = component.Description
			registeredComponent.DefaultEnabled = component.State.DefaultEnabled

			err := em.Create(registeredComponent)
			if nil != err {
				log.Warn(
					logComponentRegistry,
					"Failed to register component with code \"%v\" in database!",
					registeredComponent.Code)
			}

			em.RegisteredComponent().MarkAsAvailable(component.Code)

			continue
		}

		changed := false
		if registeredComponent.Name != component.Name {
			registeredComponent.Name = component.Name
			changed = true
		}

		if registeredComponent.Description != component.Description {
			registeredComponent.Description = component.Description
			changed = true
		}

		if registeredComponent.DefaultEnabled != component.State.DefaultEnabled {
			registeredComponent.DefaultEnabled = component.State.DefaultEnabled
			changed = true
		}

		if changed {
			err := em.Save(registeredComponent)
			if nil != err {
				log.Warn(
					logComponentRegistry,
					"Failed to update registered component for component with code \"%v\" in database!",
					registeredComponent.Code)
			}
		}

		apiDatabase.GetEntityManager().RegisteredComponent().MarkAsAvailable(component.Code)
	}
	log.Info(logComponentRegistry, "Components have been successfully registered...")
}

// LoadComponents handles the initialization of
// all components listed in the Components array.
//
// When it is not possible to register a component,
// an error will be printed into the log.
// The application will continue to run as nothing happened.
func LoadComponents(discord *discordgo.Session) {
	log.Info(logComponentRegistry, "Starting component load sequence...")
	for _, comp := range components.Components {
		if nil == comp.Lifecycle.LoadComponent {
			log.Debug(logComponentRegistry, "Component \"%v\" does not have an load callback, not loading it!", comp.Name)
			continue
		}

		log.Info(logComponentRegistry, "Loading component \"%v\"...", comp.Name)
		err := comp.RegisterComponent(discord)
		if nil != err {
			log.Warn(logComponentRegistry, "Failed to load component with name \"%v\": %v", comp.Name, err.Error())
			continue
		}
		log.Info(logComponentRegistry, "Successfully loaded component \"%v\"!", comp.Name)
	}
	log.Info(logComponentRegistry, "Component load sequence completed!")
}

// UnloadComponents iterates through all registered api.Component
// registered in the Components array and calls their UnloadComponent
// function.
//
// If an api.Component does not have an UnloadComponent function defined,
// it will be ignored.
func UnloadComponents(discord *discordgo.Session) {
	log.Info(logComponentRegistry, "Starting component unload sequence...")
	for _, comp := range components.Components {
		if nil == comp.Lifecycle.UnloadComponent {
			log.Debug(logComponentRegistry, "Component \"%v\" does not have an unload callback, skipping!", comp.Name)
			continue
		}

		if !comp.State.Loaded {
			log.Warn(logComponentRegistry, "Component \"%v\" has not been loaded, skipping!", comp.Name)
			continue
		}

		log.Info(logComponentRegistry, "Unloading component \"%v\"...", comp.Name)
		err := comp.UnregisterComponent(discord)
		if nil != err {
			log.Warn(logComponentRegistry, "Failed to unload component with name \"%v\": %v", comp.Name, err.Error())
			continue
		}
		log.Info(logComponentRegistry, "Successfully unloaded component \"%v\"!", comp.Name)
	}
	log.Info(logComponentRegistry, "Unload sequence completed!")
}
