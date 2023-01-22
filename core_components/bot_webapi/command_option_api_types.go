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
	"strings"
)

// CommandOptionDTO is an intermediate data transfer object
// that can be output or received by the WebAPI.
// This type only contains necessary data to output one or more
// command options on the web api.
//
// @Description CommandOption holds information about a single option of a command.
type CommandOptionDTO struct {
	Owner   string                   `json:"id"`
	Name    string                   `json:"name"`
	Type    int                      `json:"type"`
	Choices []CommandOptionChoiceDTO `json:"choices"`
} //@Name CommandOption

// CommandOptionChoiceDTO is an intermediate data transfer object
// that can be output or received by the WebAPI.
// This type only contains necessary data to output a command choice on the web api.
//
// @Description CommandOptionChoice holds information about a single choice of a command option.
// @Description This is only used as embedded data in a CommandOption normally.
type CommandOptionChoiceDTO struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
} //@Name CommandOptionChoice

// computeCommandOptionDTOsForCommand computes the CommandOptionDTO instances for the target command.
// The function first searches for the proper command and then collects its options.
// The options are then converted to CommandOptionDTO and CommandOptionChoiceDTO instances.
func computeCommandOptionDTOsForCommand(commands []*api.Command, id string) ([]CommandOptionDTO, error) {
	commandOptions := getCommandOptionsByCommandId(commands, id)

	if nil == commandOptions {
		return []CommandOptionDTO{}, fmt.Errorf(
			"there is no command with the id \"%s\"", id)
	}

	commandOptionDTOs := make([]CommandOptionDTO, 0)

	for i := 0; i < len(commandOptions); i++ {
		option := commandOptions[i]

		if option.Type == discordgo.ApplicationCommandOptionSubCommandGroup ||
			option.Type == discordgo.ApplicationCommandOptionSubCommand {
			continue
		}

		var choices []CommandOptionChoiceDTO
		if nil != option.Choices {
			choices = make([]CommandOptionChoiceDTO, 0)

			for _, choice := range option.Choices {
				// The application should only process supported types ofr values
				var value interface{}
				switch choice.Value.(type) {
				case string, int, int8, int16, int32, int64, float32, float64, bool, byte:
					value = choice.Value
				default:
					value = nil
				}

				choices = append(choices, CommandOptionChoiceDTO{
					choice.Name,
					value,
				})
			}
		}

		commandOptionDTOs = append(commandOptionDTOs, CommandOptionDTO{
			id,
			option.Name,
			int(option.Type),
			choices,
		})
	}

	return commandOptionDTOs, nil
}

// getCommandOptionsByCommandId returns the options of a specific (sub)command.
func getCommandOptionsByCommandId(commands []*api.Command, id string) []*discordgo.ApplicationCommandOption {
	commandPathArray := strings.Split(id, CommandIdSeparator)
	nestingLevel := len(commandPathArray)

	baseCommand := findApplicationCommandByCommandPath(commandPathArray, commands)
	if nil == baseCommand {
		return nil
	}

	options, lastOption := findOptionsForCommand(baseCommand, commandPathArray, nestingLevel)

	if hasOnlySubCommands(options) {
		return nil
	}

	if nestingLevel == 1 && nil == lastOption {
		return options
	}

	if nil == lastOption || (lastOption.Type == discordgo.ApplicationCommandOptionSubCommandGroup) {
		return nil
	}

	return options
}

// findOptionsForCommand tries to find the options for the given commandPath array.
// The command path contains a split command id.
func findOptionsForCommand(
	baseCommand *discordgo.ApplicationCommand,
	commandPath []string,
	nestingLevel int,
) ([]*discordgo.ApplicationCommandOption, *discordgo.ApplicationCommandOption) {
	options := baseCommand.Options
	if nil == options {
		options = []*discordgo.ApplicationCommandOption{}
	}

	var subCmdOption *discordgo.ApplicationCommandOption
	for i := 1; i < nestingLevel; i++ {
		currentSubCommand := commandPath[i]
		for _, option := range options {
			if option.Type != discordgo.ApplicationCommandOptionSubCommandGroup &&
				option.Type != discordgo.ApplicationCommandOptionSubCommand {
				continue
			}

			if option.Name == currentSubCommand {
				options = option.Options
				subCmdOption = option

				break
			}
		}
	}

	if nil == options {
		options = []*discordgo.ApplicationCommandOption{}
	}

	return options, subCmdOption
}

// hasOnlySubCommands checks if an array of discordgo.ApplicationCommandOption has only
// sub command groups and/or sub commands.
func hasOnlySubCommands(options []*discordgo.ApplicationCommandOption) bool {
	if 0 == len(options) {
		return false
	}

	for _, option := range options {
		if option.Type != discordgo.ApplicationCommandOptionSubCommandGroup &&
			option.Type != discordgo.ApplicationCommandOptionSubCommand {
			return false
		}
	}

	return true
}

// findApplicationCommandByCommandPath returns the discordgo.ApplicationCommand that matches the first element
// of the supplied commandPath array. If there is no matching discordgo.ApplicationCommand or commandPath
// is empty, the function will return null.
func findApplicationCommandByCommandPath(commandPath []string, commands []*api.Command) *discordgo.ApplicationCommand {
	if 0 == len(commandPath) {
		return nil
	}

	for _, cmd := range commands {
		if cmd.Cmd.Name == commandPath[0] {
			return cmd.Cmd
		}
	}

	return nil
}
