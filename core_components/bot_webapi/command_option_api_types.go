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
	Name             string            `json:"name"`
	NameLocalization map[string]string `json:"nameLocalization"`
	Value            string            `json:"value"`
} //@Name CommandOptionChoice

// computeCommandOptionDTOsForCommand computes the CommandOptionDTO instances for the target command.
// The function first searches for the proper command and then collects its options.
// The options are then converted to CommandOptionDTO and CommandOptionChoiceDTO instances.
func computeCommandOptionDTOsForCommand(commands []*api.Command, id string) ([]CommandOptionDTO, error) {
	commandPathParts := strings.Split(id, "_")

	var foundCommand *api.Command
	for cmdPath := range commandPathParts {

	}

	for _, cmd := range commands {

	}

	return []CommandOptionDTO{}, nil
}
