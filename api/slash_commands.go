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
	"errors"
	"github.com/bwmarrin/discordgo"
)

// TODO: Improve this system in a future revision.

// componentCommandMap is a map that holds the discordgo.ApplicationCommand
// as value and the name of the discordgo.ApplicationCommand as key.
var componentCommandMap map[string]*Command

// isInitialized holds the current state if the
// command handling is fully initialized or not
var isInitialized = false

// init slash command sub-system
func init() {
	componentCommandMap = make(map[string]*Command)
}

// Command is a struct that acts as a container for
// discordgo.ApplicationCommand and the assigned command Handler.
//
// Create an instance of the struct and pass to Register a command
type Command struct {
	Cmd     *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
	c       *Component
}

// ComponentSlashCommandManager is a type that is used to hold
// the owner that keeps the information about the
// component used by the slash Command manager methods.
type ComponentSlashCommandManager struct {
	owner *Component
}

// SlashCommandManager provides a standardized interface
// how slash commands should be created and registered
// in the application
type SlashCommandManager interface {
	// Register allows to register a command
	//
	// It requires a Command to be passed.
	// The Command holds the common discordgo.ApplicationCommand
	// and the function that should handle the command.
	Register(cmd *Command) error
}

// unregisterCommandHandler holds the function that can be used to unregister
// the command Handler registered by InitCommandHandling
var unregisterCommandHandler func()

// InitCommandHandling initializes the command handling
// by registering the necessary event Handler.
func InitCommandHandling(session *discordgo.Session) error {
	if nil != unregisterCommandHandler {
		return errors.New("cannot initialize command handling system twice")
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if command, ok := componentCommandMap[i.ApplicationCommandData().Name]; ok {
			command.Handler(s, i)
		}
	})

	for _, command := range componentCommandMap {
		createCommand(session, command)
	}

	isInitialized = true

	return nil
}

func createCommand(session *discordgo.Session, command *Command) {
	_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.Cmd)

	if nil != err {
		command.c.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": Could not create application Cmd!",
			command.Cmd.Name,
			command.c.Name)
	}
}

// DeinitCommandHandling unregisters the event Handler
// that is registered by InitCommandHandling.
func DeinitCommandHandling(session *discordgo.Session) {
	if nil == unregisterCommandHandler {
		return
	}

	unregisterCommandHandler()

	for _, command := range componentCommandMap {
		_ = session.ApplicationCommandDelete(session.State.User.ID, "", command.Cmd.ID)
	}

	isInitialized = false
}

// SlashCommandManager is used to obtain the components slash Command management
//
// On first call, this function initializes the private Component.slashCommandManager
// field. On consecutive calls, the already present SlashCommandManager will be used.
func (c *Component) SlashCommandManager() SlashCommandManager {
	if nil == c.slashCommandManager {
		c.slashCommandManager = &ComponentSlashCommandManager{owner: c}
	}

	return c.slashCommandManager
}

// Register allows to register a command
//
// It requires a Command to be passed.
// The Command holds the common discordgo.ApplicationCommand
// and the function that should handle the command.
func (c *ComponentSlashCommandManager) Register(cmd *Command) error {
	cmd.c = c.owner

	err := c.validateCommand(cmd)

	if nil != err {
		return err
	}

	if _, ok := componentCommandMap[cmd.Cmd.Name]; ok {
		err = errors.New("cannot register a command with the same name twice")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	componentCommandMap[cmd.Cmd.Name] = cmd

	if isInitialized {
		createCommand(c.owner.discord, cmd)
	}

	return nil
}

// validateCommand validates the passed command to ensure it is valid
// and can be registered properly.
func (c *ComponentSlashCommandManager) validateCommand(cmd *Command) error {
	if nil == cmd.Cmd {
		err := errors.New("the discordgo.ApplicationCommand of the passed command is nil")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	if nil == cmd.Handler {
		err := errors.New("the Handler of the passed command is nil")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	if nil == cmd.Handler {
		err := errors.New("the Handler of the passed command is nil")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	return nil
}
