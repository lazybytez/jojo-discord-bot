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
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
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

	// slashCommandDmPermission is used when commands are declared as global
	slashCommandDmPermission = true
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
	Cmd      *discordgo.ApplicationCommand
	Global   bool
	Category Category
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate)
	c        *Component
}

// SlashCommandManager is a type that is used to hold
// the owner that keeps the information about the
// component used by the slash Command manager methods.
type SlashCommandManager struct {
	owner *Component
}

// CommonSlashCommandManager provides a standardized interface
// how slash commands should be created and registered
// in the application.
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
	// GetCommandsForComponent returns all commands for the
	// specified component. The component needs to be specified by its unique code.
	GetCommandsForComponent(code entities.ComponentCode) []*Command
	// GetCommands returns all currently registered commands.
	GetCommands() []*Command
	// GetCommandCount returns the number of registered slash commands
	GetCommandCount() int
}

// InitCommandHandling initializes the command handling
// by registering the necessary event Handler.
func InitCommandHandling(session *discordgo.Session) error {
	if nil != unregisterCommandHandler {
		return errors.New("cannot initialize command handling system twice")
	}

	unregisterCommandHandler = session.AddHandler(handleCommandDispatch)

	return nil
}

// handleCommandDispatch handles the processing of a command
// that has been executed by a user. It already takes
// component status in account to ensure that inconsistent command
// states do not end in prohibited execution of a command.
func handleCommandDispatch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if command, ok := componentCommandMap[i.ApplicationCommandData().Name]; ok {
		user := i.User
		if nil == user {
			user = i.Member.User
		}
		if nil == user {
			command.c.Logger().Warn("Cannot handle the command \"%s\" without a user!", command.Cmd.Name)
		}

		if !IsComponentEnabled(command.c, i.GuildID) {
			resp := &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "JOJO Discord Bot",
						Description: "",
						Color:       DefaultEmbedColor,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  ":no_entry_sign: STOP :no_entry_sign:",
								Value: "An unexpected error happened while computing your error message!",
							},
						},
					},
				},
			}

			component := command.c
			registeredComponent, err := component.EntityManager().RegisteredComponent().Get(command.c.Code)
			if nil == err {
				globalStatus, err := component.EntityManager().GlobalComponentStatus().Get(registeredComponent.ID)
				if nil == err {
					switch globalStatus.Enabled {
					case true:
						resp.Embeds[0].Fields[0].Value = fmt.Sprintf("The command `/%s` is disabled on this "+
							"guild! Ask your guilds administrator to enable the `%s` component to use this command!",
							command.Cmd.Name,
							command.c.Name)
					case false:
						resp.Embeds[0].Fields[0].Value = fmt.Sprintf("The command `/%s` is globally disabled. "+
							"This might be due to some maintenance on the `%s` module.",
							command.Cmd.Name,
							command.c.Name)
					}
				}
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: resp,
			})

			if nil != err {
				command.c.Logger().Err(err, "Failed to deliver interaction response on slash-command!")

				return
			}

			command.c.Logger().Info("The user \"%s#%s\" with id \"%s\" tried to execute the "+
				"disabled command \"%s\" with options \"%s\" %s",
				user.Username,
				user.Discriminator,
				user.ID,
				command.c.SlashCommandManager().computeFullCommandStringFromInteractionData(i.ApplicationCommandData()),
				command.c.SlashCommandManager().computeConfiguredOptionsString(i.ApplicationCommandData().Options),
				getGuildOrGlobalLogPart(i.GuildID, "on"))

			return
		}

		command.c.Logger().Info("The user \"%s#%s\" with id \"%s\" executed the "+
			"command \"%s\" with options \"%s\" %s",
			user.Username,
			user.Discriminator,
			user.ID,
			command.c.SlashCommandManager().computeFullCommandStringFromInteractionData(i.ApplicationCommandData()),
			command.c.SlashCommandManager().computeConfiguredOptionsString(i.ApplicationCommandData().Options),
			getGuildOrGlobalLogPart(i.GuildID, "on"))

		command.Handler(s, i)
	}
}

// computeFullCommandStringFromInteractionData returns the full command name (e.g. jojo module enable) from the
// passed discordgo.ApplicationCommandInteractionData.
func (c *SlashCommandManager) computeFullCommandStringFromInteractionData(
	cmd discordgo.ApplicationCommandInteractionData,
) string {
	if len(cmd.Options) < 1 {
		return cmd.Name
	}

	option := cmd.Options[0]
	if option == nil {
		return cmd.Name
	}

	subCommand := c.computeSubCommandStringFromInteractionData(option)
	if subCommand == "" {
		return cmd.Name
	}

	return fmt.Sprintf("%s %s", cmd.Name, subCommand)
}

// computeSubCommandStringFromInteractionData returns, if available, the concatenated subcommand group name
// and subcommand name of the passed option.
func (c *SlashCommandManager) computeSubCommandStringFromInteractionData(
	option *discordgo.ApplicationCommandInteractionDataOption,
) string {
	if option.Type == discordgo.ApplicationCommandOptionSubCommand {
		return option.Name
	}

	if option.Type != discordgo.ApplicationCommandOptionSubCommandGroup {
		// This is not a subcommand for sure.
		return ""
	}

	if len(option.Options) == 0 {
		return fmt.Sprintf("%s %s", option.Name, "<no subcommand declared>")
	}

	subCommand := option.Options[0]
	return fmt.Sprintf("%s %s", option.Name, c.computeSubCommandStringFromInteractionData(subCommand))
}

// computeConfiguredOptions creates a string out of all configured options of a application command interaction.
func (c *SlashCommandManager) computeConfiguredOptionsString(
	options []*discordgo.ApplicationCommandInteractionDataOption,
) string {
	configuredOptions := ""

	if len(options) == 0 {
		return ""
	}

	if options[0].Type == discordgo.ApplicationCommandOptionSubCommand {
		return c.computeConfiguredOptionsString(options[0].Options)
	}

	if options[0].Type == discordgo.ApplicationCommandOptionSubCommandGroup {
		subCommandGroup := options[0]

		if len(subCommandGroup.Options) == 0 {
			return ""
		}

		return c.computeConfiguredOptionsString(subCommandGroup.Options)
	}

	for _, option := range options {
		if "" == configuredOptions {
			configuredOptions = fmt.Sprintf("%s=%v", option.Name, option.Value)

			continue
		}

		configuredOptions = fmt.Sprintf("%s; %s=%v", configuredOptions, option.Name, option.Value)
	}

	return configuredOptions
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

	if cmd.Global {
		cmd.Cmd.DMPermission = &slashCommandDmPermission
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

	return nil
}

// GetCommands returns all currently registered commands.
func (c *SlashCommandManager) GetCommands() []*Command {
	commands := make([]*Command, 0)

	for _, command := range componentCommandMap {
		commands = append(commands, command)
	}

	return commands
}

// GetCommandsForComponent returns all commands for the
// specified component. The component needs to be specified by its unique code.
func (c *SlashCommandManager) GetCommandsForComponent(code entities.ComponentCode) []*Command {
	commands := make([]*Command, 0)

	for _, command := range componentCommandMap {
		if command.c.Code == code {
			commands = append(commands, command)
		}
	}

	return commands
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
			"Failed to handle guild slash-command sync for guild \"%s\"!",
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
		"Finished syncing slash-commands for guild \"%s\"...",
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
	stillAvailableCommands := commands

	for key, registeredCommand := range commands {
		presentCompCmd, ok := componentCommandMap[registeredCommand.Name]
		if !ok || (ok && "" == guildId && !presentCompCmd.Global) || (ok && "" != guildId && presentCompCmd.Global) {
			err := session.ApplicationCommandDelete(session.State.User.ID, guildId, registeredCommand.ID)
			if nil != err {
				slashCommandManagerLogger.Err(
					err,
					"Failed to remove orphaned slash-command \"%s\" %s!",
					registeredCommand.Name,
					getGuildOrGlobalLogPart(guildId, "from"))

				continue
			}

			slicedCommands := make([]*discordgo.ApplicationCommand, 0)
			if len(commands)-1 >= key+1 {
				slicedCommands = commands[key+1:]
			}
			stillAvailableCommands = append(stillAvailableCommands[:], slicedCommands...)

			slashCommandManagerLogger.Info(
				"Removed orphaned slash-command \"%s\" %s!",
				registeredCommand.Name,
				getGuildOrGlobalLogPart(guildId, "from"))
		}
	}

	return stillAvailableCommands
}

// removeCommandsByComponentState removes commands from
// the specified guild depending on the owning components state.
func (c *SlashCommandManager) removeCommandsByComponentState(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	stillAvailableCommands := commands

	// First of all remove disabled existing commands
	for key, command := range commands {
		componentCommand, ok := componentCommandMap[command.Name]
		if !ok {
			slashCommandManagerLogger.Warn(
				"Missing component command for registered slash-command \"%s\"!",
				command.Name)

			continue
		}

		if componentCommand.Global && "" != guildId || !componentCommand.Global && "" == guildId {
			continue
		}

		if IsComponentEnabled(componentCommand.c, guildId) {
			continue
		}

		err := session.ApplicationCommandDelete(session.State.User.ID, guildId, command.ID)
		if nil != err {
			componentCommand.c.Logger().Err(
				err,
				"Failed to remove disabled slash-command \"%s\" %s!",
				command.Name,
				getGuildOrGlobalLogPart(guildId, "from"))

			continue
		}

		slicedCommands := make([]*discordgo.ApplicationCommand, 0)
		if len(commands)-1 >= key+1 {
			slicedCommands = commands[key+1:]
		}
		stillAvailableCommands = append(stillAvailableCommands[:], slicedCommands...)

		componentCommand.c.Logger().Info(
			"Removed disabled slash-command \"%s\" %s!",
			command.Name,
			getGuildOrGlobalLogPart(guildId, "from"))
	}

	return stillAvailableCommands
}

// addCommandsByComponentState removes commands from
// the specified guild depending on the owning components state.
func (c *SlashCommandManager) addCommandsByComponentState(
	session *discordgo.Session,
	guildId string,
	commands []*discordgo.ApplicationCommand,
) []*discordgo.ApplicationCommand {
	for _, componentCommand := range componentCommandMap {
		if componentCommand.Global && "" != guildId || !componentCommand.Global && "" == guildId {
			continue
		}

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
				"Failed to add enabled slash-command \"%v\" %s!",
				componentCommand.Cmd.Name,
				getGuildOrGlobalLogPart(guildId, "to"))

			continue
		}

		commands = append(commands, createdCommand)
		componentCommand.c.Logger().Info(
			"Added enabled slash-command \"%s\" %s!",
			componentCommand.Cmd.Name,
			getGuildOrGlobalLogPart(guildId, "to"))
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
			slashCommandManagerLogger.Warn("Cannot check for command updates for \"%s\" "+
				"as a corresponding component command is missing!",
				command.Name)

			continue
		}

		if componentCommand.Global && "" != guildId || !componentCommand.Global && "" == guildId {
			continue
		}

		if c.compareCommands(command, componentCommand.Cmd) {
			continue
		}

		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, componentCommand.Cmd)
		if nil != err {
			componentCommand.c.Logger().Err(
				err,
				"Failed to add slash-command \"%s\" %s during command update!",
				componentCommand.Cmd.Name,
				getGuildOrGlobalLogPart(guildId, "to"))

			continue
		}

		commands[key] = createdCommand
		componentCommand.c.Logger().Info(
			"Updated slash-command \"%s\" %s!",
			componentCommand.Cmd.Name,
			getGuildOrGlobalLogPart(guildId, "for"))
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

// getGuildOrGlobalLogPart returns "globally" for an empty guild id
// or "of guild <id>" when a non-empty guild id has been passed.
func getGuildOrGlobalLogPart(guildId string, prefix string) string {
	if "" == guildId {
		return "globally"
	}

	return fmt.Sprintf("%s guild \"%s\"", prefix, guildId)
}

// compareCommands compares to discordgo.Application commands
// using some key factors and returns if they are equal or not
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

	if !util.PointerValuesEqual(a.DefaultMemberPermissions, b.DefaultMemberPermissions) {
		return false
	}

	if !util.PointerValuesEqual(a.DMPermission, b.DMPermission) {
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

	if !util.PointerValuesEqual(a.MinLength, b.MinLength) {
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

// GetComponentCode returns the code of the component that owns the command.
func (c Command) GetComponentCode() entities.ComponentCode {
	return c.c.Code
}
