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

package bot_webapi

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"strings"
)

// ParamCommandID is the name of the parameter that carries a
// requested command name.
const ParamCommandID = "id"

// CommandIdSeparator separator used for command ids.
// A command id consists of the entire path of command names until the target is reached.
const CommandIdSeparator = "_"

// CommandDTOsWebApiCacheKey is the cache key used to store and retrieve all commands
// as CommandDTO instances from the cache.
const CommandDTOsWebApiCacheKey = "bot_web_api_commands_get_cache"

// CommandDTO is an intermediate data transfer object
// that can be output or received by the WebAPI.
// This type is used, because the bot both has the general api.Command.
// This type represents a  command with only API relevant data.
//
// @Description Command holds information about a slash-command like its name,
// @Description description and its options. Note that commands are always
// @Description built from the deepest level commands. This means a command is either a sub command or the
// @Description top-level command.
type CommandDTO struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Component   entities.ComponentCode `json:"component"`
	Category    api.Category           `json:"category"`
	Description string                 `json:"description"`
} //@Name Command

// CommandDTOsFromCommands creates an array of CommandDTO instances.
// The general data is computed from the passed api.Command.
// This function returns all commands that are currently registered.
func CommandDTOsFromCommands(cmds []*api.Command) []CommandDTO {
	commandDTOs := make([]CommandDTO, 0)

	for _, cmd := range cmds {
		subDTOs := commandDTOsFromCommand(cmd)

		commandDTOs = append(commandDTOs, subDTOs...)
	}

	return commandDTOs
}

// commandDTOsFromCommand converts a single api.Command into CommandDTO instances.
func commandDTOsFromCommand(cmd *api.Command) []CommandDTO {
	commandDTOs := make([]CommandDTO, 0)

	if nil != cmd.Cmd.Options {
		cmdDTO := commandDTOsFromCommandOptions(cmd.Cmd.Name,
			cmd.Category,
			cmd.GetComponentCode(),
			cmd.Cmd.Options)

		commandDTOs = append(commandDTOs, cmdDTO...)
	}

	if 0 == len(commandDTOs) {
		commandDTO := CommandDTO{
			ID:          getCommandIDFromCommandDTOName(cmd.Cmd.Name),
			Name:        cmd.Cmd.Name,
			Component:   cmd.GetComponentCode(),
			Category:    cmd.Category,
			Description: cmd.Cmd.Description,
		}

		commandDTOs = append(commandDTOs, commandDTO)
	}

	return commandDTOs
}

// commandDTOsFromCommandOptions converts options that are subcommands or subcommand groups to
// CommandDTO instances.
func commandDTOsFromCommandOptions(
	parentName string,
	category api.Category,
	component entities.ComponentCode,
	options []*discordgo.ApplicationCommandOption,
) []CommandDTO {
	if nil == options {
		return nil
	}

	commandDTOs := make([]CommandDTO, 0)

	for _, cmdOption := range options {
		switch cmdOption.Type {
		case discordgo.ApplicationCommandOptionSubCommandGroup:
			subCommandDTOs := commandDTOsFromCommandOptions(fmt.Sprintf("%s %s", parentName, cmdOption.Name),
				category,
				component,
				cmdOption.Options)

			commandDTOs = append(commandDTOs, subCommandDTOs...)
		case discordgo.ApplicationCommandOptionSubCommand:
			name := fmt.Sprintf("%s %s", parentName, cmdOption.Name)
			cmdDTO := CommandDTO{
				ID:          getCommandIDFromCommandDTOName(name),
				Name:        name,
				Component:   component,
				Category:    category,
				Description: cmdOption.Description,
			}

			commandDTOs = append(commandDTOs, cmdDTO)
		}
	}

	return commandDTOs
}

// getCommandDTOs returns the command DTOs from all commands that are currently registered in the bot.
// If an error occurs, an error will be returned.
func getCommandDTOs() []CommandDTO {
	commandsWebAPICache, ok := cache.Get(CommandDTOsWebApiCacheKey, []CommandDTO{})
	if ok {
		return commandsWebAPICache
	}

	commands := C.SlashCommandManager().GetCommands()
	commandDTOs := CommandDTOsFromCommands(commands)

	err := cache.Update(CommandDTOsWebApiCacheKey, commandDTOs)
	if nil != err {
		C.Logger().Warn(fmt.Sprintf("Failed to cache aggregated CommandDTOs: %s", err.Error()))
	}

	return commandDTOs
}

// getCommandIDFromCommandDTOName returns the ID of a command depending on its name
// of the CommandDTO.
func getCommandIDFromCommandDTOName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", CommandIdSeparator))
}
