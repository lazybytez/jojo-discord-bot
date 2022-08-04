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

// componentCommandMap is a map that holds the discordgo.ApplicationCommand
// as value and the name of the discordgo.ApplicationCommand as key.
var componentCommandMap map[string]*Command

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

	return nil
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
	var err error
	if nil != err {
		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": Could not create application Cmd!",
			cmd.Cmd.Name,
			c.owner.Name)

		return err
	}

	if nil == cmd.Cmd {
		err = errors.New("the discordgo.ApplicationCommand of the passed command is nil")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	if nil == cmd.Handler {
		err = errors.New("the Handler of the passed command is nil")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	if nil == cmd.Handler {
		err = errors.New("the Handler of the passed command is nil")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-Cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

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

	return nil
}

// FinishCommandRegistration is called by the internal system
// and finishes the final registration steps necessary to get commands working.
func FinishCommandRegistration(session *discordgo.Session) {
	for _, command := range componentCommandMap {
		_, _ = session.ApplicationCommandCreate(session.State.User.ID, "", command.Cmd)
	}
}
