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
)

// LifecycleHooks allow to specify functions that should be called
// when components get loaded and unloaded.
//
// The defined functions allow some way of initializing a component.
type LifecycleHooks struct {
	LoadComponent   func(discord *discordgo.Session) error
	UnloadComponent func(discord *discordgo.Session) error
}

// State holds the state of a component.
// This includes states like:
//   - is the component enabled?
//   - is the component currently loaded?
type State struct {
	Loaded  bool
	Enabled bool
}

// Component is the base structure that must be initialized
// to create a new component.
//
// It holds basic metadata about the component
type Component struct {
	// Metadata
	Name         string
	Description  string
	DmPermission bool

	// State
	State State

	// Lifecycle hooks
	Lifecycle LifecycleHooks

	// Utilities
	// These are private and only managed by the API system.
	// Their initialization happens through call to the methods
	// used to get them (Example: logger -> Component.Logger()).
	logger         Logger
	handlerManager ComponentHandlerManager
	discord        *discordgo.Session
}

// RegistrableComponent is the interface that allows a component to be
// initialized and registered.
type RegistrableComponent interface {
	RegisterComponent(discord *discordgo.Session) error
	UnregisterComponent(discord *discordgo.Session) error
}

// RegisterComponent is used by the component registration system that
// automatically calls the RegisterComponent method for all Component instances in
// the components.Components array.
func (c *Component) RegisterComponent(discord *discordgo.Session) error {
	c.discord = discord

	err := c.Lifecycle.LoadComponent(discord)

	if err != nil {
		return err
	}
	c.State.Loaded = true

	return nil
}

// UnregisterComponent is used by the component registration system that
// automatically calls the UnregisterComponent method for all Component instances in
// the components.Components array.
//
// The function takes care of tasks like unregistering slash-commands and so on.
//
// It is used to give components the ability to gracefully shutdown.
func (c *Component) UnregisterComponent(discord *discordgo.Session) error {
	c.HandlerManager().unregisterAll()

	err := c.Lifecycle.UnloadComponent(discord)

	if err != nil {
		return err
	}
	c.State.Loaded = false

	return nil
}
