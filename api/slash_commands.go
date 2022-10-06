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
	"github.com/lazybytez/jojo-discord-bot/api/util"
	"github.com/lazybytez/jojo-discord-bot/services/logger"
)

// slashCommandLogPrefix is the prefix used by the log management during
// lifecycle events that cannot be assigned to a specific component.
const slashCommandLogPrefix = "slash_command_manager"

var (
	// componentCommandMap is a map that holds the discordgo.ApplicationCommand
	// as value and the name of the discordgo.ApplicationCommand as key.
	componentCommandMap map[string]*Command

	// unregisterCommandHandler holds the function that can be used to unregister
	// the command Handler registered by InitCommandHandling
	unregisterCommandHandler func()

	// slashCommandManagerLogger is the logger used by the slash command management
	// when there is no component a log message could be assigned to
	slashCommandManagerLogger *logger.Logger
)

// init slash command sub-system
func init() {
	componentCommandMap = make(map[string]*Command)
	slashCommandManagerLogger = logger.New(slashCommandLogPrefix, nil)
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

// SlashCommandManager is a type that is used to hold
// the owner that keeps the information about the
// component used by the slash Command manager methods.
type SlashCommandManager struct {
	owner *Component
}

// CommonSlashCommandManager provides a standardized interface
// how slash commands should be created and registered
// in the application
type CommonSlashCommandManager interface {
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
	// GetCommandCount returns the number of registered slash commands
	GetCommandCount() int
}

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

	// Drop global registered commands, if any.
	// We only allow guild specific commands
	//
	// This is for backward compatibility, as in the past
	// we registered commands globally.
	//
	// TODO: Remove in some safe future version!
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if nil != err {
		slashCommandManagerLogger.Err(
			err,
			"Failed to retrieve globally registered commands!")
	}
	for _, cmd := range commands {
		err = session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID)
		if nil != err {
			slashCommandManagerLogger.Err(
				err,
				"Failed to remove global slash-command with name \"%v\"!",
				cmd.Name)
		}
	}

	return nil
}

// DeinitCommandHandling unregisters the event Handler
// that is registered by InitCommandHandling.
func DeinitCommandHandling() {
	if nil == unregisterCommandHandler {
		return
	}

	unregisterCommandHandler()
}

// SlashCommandManager is used to obtain the components slash Command management
//
// On first call, this function initializes the private Component.slashCommandManager
// field. On consecutive calls, the already present CommonSlashCommandManager will be used.
func (c *Component) SlashCommandManager() *SlashCommandManager {
	if nil == c.slashCommandManager {
		c.slashCommandManager = &SlashCommandManager{owner: c}
	}

	return c.slashCommandManager
}

// Register allows to register a command
//
// It requires a Command to be passed.
// The Command holds the common discordgo.ApplicationCommand
// and the function that should handle the command.
func (c *SlashCommandManager) Register(cmd *Command) error {
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

	if 0 == cmd.Cmd.Type {
		cmd.Cmd.Type = discordgo.ChatApplicationCommand
	}

	componentCommandMap[cmd.Cmd.Name] = cmd

	return nil
}

// validateCommand validates the passed command to ensure it is valid
// and can be registered properly.
func (c *SlashCommandManager) validateCommand(cmd *Command) error {
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
			"Failed to register the slash-Cmd \"%v\" for component \"%v\"!",
			cmd.Cmd.Name,
			c.owner.Name)

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
//
// Sync is a four-step process:
//   - remove orphaned commands
//   - remove disabled commands
//   - add new commands
//   - update existing commands
func (c *SlashCommandManager) SyncApplicationComponentCommands(
	session *discordgo.Session,
	guildId string,
) {
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, guildId)
	if nil != err {
		slashCommandManagerLogger.Err(
			err,
			"Failed to handle guild slash-command sync for guild \"%v\"!",
			guildId)

		return
	}

	slashCommandManagerLogger.Info(
		"Syncing slash-commands for guild \"%v\"...",
		guildId)
	registeredCommands = c.removeOrphanedCommands(session, guildId, registeredCommands)
	registeredCommands = c.removeCommandsByComponentState(session, guildId, registeredCommands)
	registeredCommands = c.addCommandsByComponentState(session, guildId, registeredCommands)
	_ = c.updateRegisteredCommands(session, guildId, registeredCommands)

	slashCommandManagerLogger.Info(
		"Finished syncing slash-commands for guild \"%v\"...",
		guildId)
}

// removeOrphanedCommands removes all slash-commands commands from a guild,
// that are no longer registered in the bots slash-command management.
//
// The function returns the passed list of commands,
// with the removed commands being removed.
func (c *SlashCommandManager) removeOrphanedCommands(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	for key, registeredCommand := range commands {
		if _, ok := componentCommandMap[registeredCommand.Name]; !ok {
			err := session.ApplicationCommandDelete(session.State.User.ID, guildId, registeredCommand.ID)
			if nil != err {
				slashCommandManagerLogger.Err(
					err,
					"Failed to remove orphaned slash-command \"%v\" from guild \"%v\"!",
					registeredCommand.Name,
					guildId)

				continue
			}

			slicedCommands := make([]*discordgo.ApplicationCommand, 0)
			if len(commands)-1 >= key+1 {
				slicedCommands = commands[key+1:]
			}
			commands = append(commands[:key], slicedCommands...)
			slashCommandManagerLogger.Info(
				"Removed orphaned slash-command \"%v\" from guild \"%v\"!",
				registeredCommand.Name,
				guildId)
		}
	}

	return commands
}

// removeCommandsByComponentState removes commands from
// the specified guild depending on the owning components state.
func (c *SlashCommandManager) removeCommandsByComponentState(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	// First of all remove disabled existing commands
	for key, command := range commands {
		componentCommand, ok := componentCommandMap[command.Name]
		if !ok {
			slashCommandManagerLogger.Warn(
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
				"Failed to remove disabled slash-command \"%v\" from guild \"%v\"!",
				command.Name,
				guildId)

			continue
		}

		slicedCommands := make([]*discordgo.ApplicationCommand, 0)
		if len(commands)-1 >= key+1 {
			slicedCommands = commands[key+1:]
		}
		commands = append(commands[:key], slicedCommands...)

		componentCommand.c.Logger().Info(
			"Removed disabled slash-command \"%v\" from guild \"%v\"!",
			command.Name,
			guildId)
	}

	return commands
}

// addCommandsByComponentState removes commands from
// the specified guild depending on the owning components state.
func (c *SlashCommandManager) addCommandsByComponentState(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	for _, componentCommand := range componentCommandMap {
		if c.isCommandNameInApplicationCommandList(commands, componentCommand.Cmd.Name) {
			continue
		}

		if !IsComponentEnabled(componentCommand.c, guildId) {
			continue
		}

		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, componentCommand.Cmd)
		if nil != err {
			componentCommand.c.Logger().Err(
				err,
				"Failed to add enabled slash-command \"%v\" to guild \"%v\"!",
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

// updateRegisteredCommands checks all registered commands
// and updates (deletes and creates) differing commands.
func (c *SlashCommandManager) updateRegisteredCommands(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	for key, command := range commands {
		componentCommand, ok := componentCommandMap[command.Name]
		if !ok {
			slashCommandManagerLogger.Warn("Cannot check for command updates for \"%v\" "+
				"as a corresponding component command is missing!",
				command.Name)

			continue
		}

		if c.compareCommands(command, componentCommand.Cmd) {
			continue
		}

		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, componentCommand.Cmd)
		if nil != err {
			componentCommand.c.Logger().Err(
				err,
				"Failed to add slash-command \"%v\" to guild \"%v\" during command update!",
				componentCommand.Cmd.Name,
				guildId)

			continue
		}

		commands[key] = createdCommand
		componentCommand.c.Logger().Info(
			"Updated slash-command \"%v\" of guild \"%v\"!",
			componentCommand.Cmd.Name,
			guildId)
	}

	return commands
}

// isCommandNameInApplicationCommandList checks if a command with the provided name
// is present in the provided discordgo.ApplicationCommand slice.
func (c *SlashCommandManager) isCommandNameInApplicationCommandList(commands []*discordgo.ApplicationCommand, name string) bool {
	for _, command := range commands {
		if command.Name == name {
			return true
		}
	}

	return false
}

// compareCommands compares to discordgo.Application commands
// using some key factors and returns of they are equal or not
func (c *SlashCommandManager) compareCommands(
	a *discordgo.ApplicationCommand,
	b *discordgo.ApplicationCommand,
) bool {
	// Common data
	if a.Name != b.Name {
		return false
	}
	if !util.MapsEqual(a.NameLocalizations, b.NameLocalizations) {
		return false
	}

	if a.Description != b.Description {
		return false
	}
	if !util.MapsEqual(a.DescriptionLocalizations, b.DescriptionLocalizations) {
		return false
	}

	if a.DefaultMemberPermissions != b.DefaultMemberPermissions {
		return false
	}

	if a.DMPermission != b.DMPermission {
		return false
	}

	if uint8(a.Type) != uint8(b.Type) {
		return false
	}

	// Options
	if len(a.Options) != len(b.Options) {
		return false
	}

	for key, optionA := range a.Options {
		optionB := b.Options[key]
		if nil == optionB {
			return false
		}

		if !c.compareCommandOptions(optionA, optionB) {
			return false
		}
	}

	return true
}

// compareCommandOptions compares the options of a command with another
func (c *SlashCommandManager) compareCommandOptions(a *discordgo.ApplicationCommandOption, b *discordgo.ApplicationCommandOption) bool {
	// Common data
	if a.Name != b.Name {
		return false
	}
	if !util.MapsEqual(&a.NameLocalizations, &b.NameLocalizations) {
		return false
	}

	if a.Description != b.Description {
		return false
	}
	if !util.MapsEqual(&a.DescriptionLocalizations, &b.DescriptionLocalizations) {
		return false
	}

	if uint8(a.Type) != uint8(b.Type) {
		return false
	}

	if a.Required != b.Required {
		return false
	}

	if a.Autocomplete != b.Autocomplete {
		return false
	}

	if a.MaxLength != b.MaxLength {
		return false
	}

	if a.MinLength != b.MinLength {
		return false
	}

	// Channel types
	if !util.ArraysEqual(&a.ChannelTypes, &b.ChannelTypes) {
		return false
	}

	// Options
	if len(a.Options) != len(b.Options) {
		return false
	}

	for k, optionA := range a.Options {
		optionB := b.Options[k]
		if nil == optionB {
			return false
		}

		if !c.compareCommandOptions(optionA, optionB) {
			return false
		}
	}

	return compareCommandOptionChoices(a.Choices, b.Choices)
}

// compareCommandOptionChoices compares the choices of a command option
// and returns the result.
func compareCommandOptionChoices(
	a []*discordgo.ApplicationCommandOptionChoice,
	b []*discordgo.ApplicationCommandOptionChoice,
) bool {
	if len(a) != len(b) {
		return false
	}

	for k, choiceA := range a {
		choiceB := b[k]
		if nil == choiceB {
			return false
		}

		if choiceA.Name != choiceB.Name {
			return false
		}
		if !util.MapsEqual(&choiceA.NameLocalizations, &choiceB.NameLocalizations) {
			return false
		}

		if choiceA.Value != choiceB.Value {
			return false
		}
	}

	return true
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

// GetCommandCount returns the number of registered slash commands
func (c *SlashCommandManager) GetCommandCount() int {
	return len(componentCommandMap)
}
