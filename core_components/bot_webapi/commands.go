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
	"github.com/gin-gonic/gin"
	"github.com/lazybytez/jojo-discord-bot/api"
	"net/http"
	"strings"
)

// ParamCommandID is the name of the parameter that carries the
// requested command name.
const ParamCommandID = "id"

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
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Category    api.Category `json:"category"`
	Description string       `json:"description"`
} //@Name Command

// commandsWebAPICache acts as a cache to prevent multiple computation
// of the commands web api content.
var commandsWebAPICache []*CommandDTO = nil

// singleCommandWebAPICache acts as a cache to prevent multiple computation
// of the specific commands web api content.
var singleCommandWebAPICache map[string]*CommandDTO = map[string]*CommandDTO{}

// CommandDTOsFromCommands creates an array of CommandDTO instances.
// The general data is computed from the passed api.Command.
// This function returns all commands that are currently registered.
func CommandDTOsFromCommands(cmds []*api.Command) ([]*CommandDTO, error) {
	commandDTOs := make([]*CommandDTO, 0)

	for _, cmd := range cmds {
		subDTOs, err := commandDTOsFromCommand(cmd)

		if nil != err {
			return []*CommandDTO{}, err
		}

		commandDTOs = append(commandDTOs, subDTOs...)
	}

	return commandDTOs, nil
}

// commandDTOsFromCommand converts a single api.Command into CommandDTO instances.
func commandDTOsFromCommand(cmd *api.Command) ([]*CommandDTO, error) {
	commandDTOs := make([]*CommandDTO, 0)

	if nil != cmd.Cmd.Options {
		cmdDTO, err := commandDTOsFromCommandOptions(cmd.Cmd.Name, cmd.Category, cmd.Cmd.Options)
		if nil != err {
			return []*CommandDTO{}, err
		}

		commandDTOs = append(commandDTOs, cmdDTO...)
	}

	if 0 == len(commandDTOs) {
		commandDTO := &CommandDTO{
			ID:          getCommandIDFromCommandDTOName(cmd.Cmd.Name),
			Name:        cmd.Cmd.Name,
			Category:    cmd.Category,
			Description: cmd.Cmd.Description,
		}

		commandDTOs = append(commandDTOs, commandDTO)
	}

	return commandDTOs, nil
}

// commandDTOsFromCommandOptions converts options that are subcommands or subcommand groups to
// CommandDTO instances.
func commandDTOsFromCommandOptions(
	parentName string,
	category api.Category,
	options []*discordgo.ApplicationCommandOption,
) ([]*CommandDTO, error) {
	if nil == options {
		return nil, nil
	}

	commandDTOs := make([]*CommandDTO, 0)

	for _, cmdOption := range options {
		switch cmdOption.Type {
		case discordgo.ApplicationCommandOptionSubCommandGroup:
			subCommandDTOs, err := commandDTOsFromCommandOptions(fmt.Sprintf("%s %s", parentName, cmdOption.Name),
				category,
				cmdOption.Options)

			if nil != err {
				return commandDTOs, err
			}

			commandDTOs = append(commandDTOs, subCommandDTOs...)
		case discordgo.ApplicationCommandOptionSubCommand:
			name := fmt.Sprintf("%s %s", parentName, cmdOption.Name)
			cmdDTO := &CommandDTO{
				ID:          getCommandIDFromCommandDTOName(name),
				Name:        name,
				Category:    category,
				Description: cmdOption.Description,
			}

			commandDTOs = append(commandDTOs, cmdDTO)
		}
	}

	return commandDTOs, nil
}

// CommandsGet endpoint
//
// @Summary     Get all available commands of the bot
// @Description This endpoint collects all available commands and returns them.
// @Description The result on a success contains relevant information like name, description and category of commands.
// @Description Note that this endpoint does not return detailed information like the options of a command.
// @Description To obtain the available command options, the command must be queried on its own using the
// @Description single command options get endpoint.
// @Tags        Component System
// @Produce     json
// @Success     200 {array} CommandDTO "Returns an array of commands with all their available data."
// @Router      /commands/ [get]
func CommandsGet(g *gin.Context) {
	if nil != commandsWebAPICache {
		g.JSON(http.StatusOK, commandsWebAPICache)

		return
	}

	commandDTOs, err := getCommandDTOs()

	if nil != err {
		C.Logger().Err(err, "Failed to convert some commands to CommandDTO!")
	}

	commandsWebAPICache = commandDTOs

	g.JSON(http.StatusOK, commandDTOs)
}

// CommandGet endpoint
//
// @Summary     Get a specific command of the bot.
// @Description This endpoint returns information about the specific requested command.
// @Description The result on a success contains relevant information like name, description and category of commands.
// @Description Note that this endpoint does not return detailed information like the options of a command.
// @Description To obtain the available command options, the command must be queried on its own using the
// @Description single command options get endpoint.
// @Tags        Component System
// @Param		id path string true "ID of the command to search for"
// @Produce     json
// @Success     200 {array} CommandDTO "Returns an object of a specific command with all its available data."
// @Router      /commands/{id} [get]
func CommandGet(g *gin.Context) {
	cmdID := g.Param(ParamCommandID)

	if "" == cmdID {
		g.AbortWithStatus(400)
		return
	}

	commandDTOs, err := getCommandDTOs()

	if nil != err {
		C.Logger().Err(err, "Failed to convert some commands to CommandDTO!")
	}

	cmdDTO, ok := singleCommandWebAPICache[cmdID]
	if !ok {
		for _, cmd := range commandDTOs {
			if cmd.ID == cmdID {
				cmdDTO = cmd

				break
			}
		}

		if nil == cmdDTO {
			C.Logger().Err(err, "Could not find a command with the given name!")
			return
		}

		singleCommandWebAPICache[cmdDTO.ID] = cmdDTO
	}

	g.JSON(http.StatusOK, cmdDTO)
}

func CommandOptionsGet(g *gin.Context) {

}

// getCommandDTOs returns the command DTOs from all commands that are currently registered in the bot.
// If an error occurs, an error will be returned.
func getCommandDTOs() ([]*CommandDTO, error) {
	if nil != commandsWebAPICache {
		return commandsWebAPICache, nil
	}

	commands := C.SlashCommandManager().GetCommands()
	commandDTOs, err := CommandDTOsFromCommands(commands)

	if nil != err {
		return []*CommandDTO{}, err
	}

	commandsWebAPICache = commandDTOs

	return commandsWebAPICache, nil
}

// getCommandIDFromCommandDTOName returns the ID of a command depending on its name
// of the CommandDTO.
func getCommandIDFromCommandDTOName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}
