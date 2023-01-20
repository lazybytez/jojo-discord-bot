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
	"github.com/gin-gonic/gin"
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"github.com/lazybytez/jojo-discord-bot/webapi"
	"net/http"
	"time"
)

// SpecificCommandDTOWebApiCacheKey is the cache key format used to store and retrieve a specific command
// as CommandDTO instance from the cache. There is exactly one placeholder in this constant,
// that should be replaced with the ID of the command to get.
const SpecificCommandDTOWebApiCacheKey = "bot_web_api_specific_command_get_%s_cache"

// CommandGet endpoint
//
// @Summary     Get a specific command of the bot.
// @Description This endpoint returns information about the specific requested command.
// @Description The result on a success contains relevant information like name, description and category of commands.
// @Description Note that this endpoint does not return detailed information like the options of a command.
// @Description To obtain the available command options, the command must be queried on its own using the
// @Description single command options get endpoint.
// @Tags        Command System
// @Param		id path string true "ID of the command to search for"
// @Produce     json
// @Success     200 {array} CommandDTO "An object containing information about a specific command"
// @Failure		404 {object} webapi.ErrorResponse "An error indicating that the requested resource could not be found"
// @Failure		500 {object} webapi.ErrorResponse "An error indicating that an internal error happened"
// @Router      /commands/{id} [get]
func CommandGet(g *gin.Context) {
	cmdID := g.Param(ParamCommandID)

	if "" == cmdID {
		webapi.RespondWithError(g, webapi.ErrorResponse{
			Status:    http.StatusBadRequest,
			Error:     "Missing command id parameter",
			Message:   "No value has been supplied for the command id path parameter",
			Timestamp: time.Now(),
		})

		return
	}

	commandDTOs := getCommandDTOs()

	cacheKey := getSpecificCommandDTOCacheKey(cmdID)
	cmdDTO, ok := cache.Get(cacheKey, CommandDTO{})
	if !ok {
		for _, cmd := range commandDTOs {
			if cmd.ID == cmdID {
				cmdDTO = cmd
				ok = true

				break
			}
		}

		if !ok {
			webapi.RespondWithError(g, webapi.ErrorResponse{
				Status:    http.StatusNotFound,
				Error:     "Failed to find the desired command",
				Message:   "Could not find a command matching the given criteria",
				Timestamp: time.Now(),
			})

			return
		}

		cache.Update(cacheKey, cmdDTO)
	}

	g.JSON(http.StatusOK, cmdDTO)
}

// getSpecificCommandDTOCacheKey returns the cache key to get a specific CommandDTO from cache.
func getSpecificCommandDTOCacheKey(id string) string {
	return fmt.Sprintf(SpecificCommandDTOWebApiCacheKey, id)
}
