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
	"github.com/lazybytez/jojo-discord-bot/api/log"
)

// slashCommandLogPrefix is the prefix used by the log management during
// lifecycle events that cannot be assigned to a specific component.
const slashCommandLogPrefix = "slash_command_manager"

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
	// SyncApplicationComponentCommands ensures that the available discordgo.ApplicationCommand
	// are synced for the given component with the given guild.
	//
	// This means that disabled commands are enabled and enabled commands are disabled
	// depending on the component enable state.
	//
	// Also orphaned commands are cleaned up.
	// This is executed whenever a guild is joined or a component is toggled.
	SyncApplicationComponentCommands(session *discordgo.Session, guildId string)
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

	unregisterCommandHandler = session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	cmd.c = c.owner

	err := c.validateCommand(cmd)

	if nil != err {
		return err
	}

	if _, ok := componentCommandMap[cmd.Cmd.Name]; ok {
		err = errors.New("cannot register a command with the same name twice")

		c.owner.Logger().Err(
			err,
			"Failed to register the slash-cmd \"%v\" for component \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	componentCommandMap[cmd.Cmd.Name] = cmd

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
			"Failed to register the slash-Cmd \"%v\" for component \"%v\" on guild \"%v\": %v!",
			cmd.Cmd.Name,
			c.owner.Name,
			err.Error())

		return err
	}

	return nil
}

// SyncApplicationComponentCommands ensures that the available discordgo.ApplicationCommand
// are synced for the given component with the given guild.
//
// This means that disabled commands are enabled and enabled commands are disabled
// depending on the component enable state.
//
// Also orphaned commands are cleaned up.
// This is executed whenever a guild is joined or a component is toggled.
func (c *ComponentSlashCommandManager) SyncApplicationComponentCommands(
	session *discordgo.Session,
	guildId string,
) {
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, guildId)
	if nil != err {
		log.Err(
			slashCommandLogPrefix,
			err,
			"Failed to handle guild slash-command sync for guild \"%v\"!",
			guildId)

		return
	}

	log.Info(
		slashCommandLogPrefix,
		"Syncing slash-commands for guild \"%v\"...",
		guildId)
	registeredCommands = c.removeOrphanedCommands(session, guildId, registeredCommands)
	registeredCommands = c.removeCommandsByComponentState(session, guildId, registeredCommands)

	log.Info(
		slashCommandLogPrefix,
		"Finished syncing slash-commands for guild \"%v\"...",
		guildId)
}

// removeOrphanedCommands removes all slash-commands commands from a guild,
// that are no longer registered in the bots slash-command management.
//
// The function returns the passed list of commands,
// with the removed commands being removed.
func (c *ComponentSlashCommandManager) removeOrphanedCommands(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	for key, registeredCommand := range commands {
		if _, ok := componentCommandMap[registeredCommand.Name]; !ok {
			err := session.ApplicationCommandDelete(session.State.User.ID, guildId, registeredCommand.ID)
			if nil != err {
				log.Err(
					slashCommandLogPrefix,
					err,
					"Failed to remove orphaned slash-command \"%v\" from guild \"%v\"",
					registeredCommand.Name,
					guildId)

				continue
			}

			commands = append(commands[:key], commands[key+1:]...)
			log.Info(
				slashCommandLogPrefix,
				"Removed orphaned slash-command \"%v\" from guild \"%v\"!",
				registeredCommand.Name,
				guildId)
		}
	}

	return commands
}

// removeCommandsByComponentState removes commands from
// the specified guild depending on the owning components state.
func (c *ComponentSlashCommandManager) removeCommandsByComponentState(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	// First of all remove disabled existing commands
	for key, command := range commands {
		componentCommand, ok := componentCommandMap[command.Name]
		if !ok {
			log.Warn(
				slashCommandLogPrefix,
				"Missing component command for registered slash-command \"%v\"!",
				command.Name)

			continue
		}

		if IsComponentEnabled(componentCommand.c, guildId) {
			continue
		}

		err := session.ApplicationCommandDelete(session.State.User.ID, guildId, command.ID)
		if nil != err {
			componentCommand.c.Logger().Err(
				err,
				"Failed to remove disabled slash-command \"%v\" from guild \"%v\"",
				command.Name,
				guildId)

			continue
		}

		commands = append(commands[:key], commands[key+1:]...)

		componentCommand.c.Logger().Info(
			"Removed disabled slash-command \"%v\" from guild \"%v\"!",
			command.Name,
			guildId)
	}

	return commands
}

// addCommandsByComponentState removes commands from
// the specified guild depending on the owning components state.
func (c *ComponentSlashCommandManager) addCommandsByComponentState(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	for _, componentCommand := range componentCommandMap {
		if isCommandNameInApplicationCommandList(commands, componentCommand.Cmd.Name) {
			continue
		}

		if !IsComponentEnabled(componentCommand.c, guildId) {
			continue
		}

		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, componentCommand.Cmd)
		if nil != err {
			componentCommand.c.Logger().Err(
				err,
				"Failed to add enabled slash-command \"%v\" to guild \"%v\"",
				componentCommand.Cmd.Name,
				guildId)

			continue
		}

		commands = append(commands, createdCommand)
		componentCommand.c.Logger().Info(
			"Added enabled slash-command \"%v\" to guild \"%v\"!",
			componentCommand.Cmd.Name,
			guildId)
	}

	return commands
}

// isCommandNameInApplicationCommandList checks if a command with the provided name
// is present in the provided discordgo.ApplicationCommand slice.
func isCommandNameInApplicationCommandList(commands []*discordgo.ApplicationCommand, name string) bool {
	for _, command := range commands {
		if command.Name == name {
			return true
		}
	}

	return false
}

// ProcessSubCommands is an easy way to handle sub-commands and sub-command-groups.
// The function will return true if there was a valid sub-command to handle
func ProcessSubCommands(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
	handlers map[string]func(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
),
) bool {
	// First validate that there is at least one level of nesting
	command := i.ApplicationCommandData()
	if len(command.Options) < 1 {
		return false
	}

	if option == nil {
		option = command.Options[0]

		return runHandler(s, i, option, option.Name, handlers)
	}

	if len(option.Options) < 1 {
		return false
	}

	option = option.Options[0]

	return runHandler(s, i, option, option.Name, handlers)
}

// runHandler executes the actual handle
// of the found sub-command.
//
// Returns false if there is no suiting sub-command.
func runHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
	name string,
	handlers map[string]func(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
),
) bool {
	handler, ok := handlers[name]

	if !ok {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The executed (sub)command is invalid or does not exist!",
			},
		})

		return false
	}

	handler(s, i, option)

	return true
}
