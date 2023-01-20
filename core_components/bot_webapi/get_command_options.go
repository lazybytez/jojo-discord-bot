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

const CommandOptionsDTOWebApiCacheKey = "bot_web_api_specific_command_options_get_%s_cache"

// CommandOptionsGet endpoint
//
// @Summary     Get the options of a specific command of the bot.
// @Description This endpoint returns information about the options of the specified command.
// @Description The result on a success contains relevant information about the options and their choices.
// @Tags        Command System
// @Param		id path string true "ID of the command to search for options"
// @Produce     json
// @Success     200 {array} CommandOptionDTO "An object containing information about the options of a specific command"
// @Failure		404 {object} webapi.ErrorResponse "An error indicating that the requested resource could not be found"
// @Failure		500 {object} webapi.ErrorResponse "An error indicating that an internal error happened"
// @Router      /commands/{id}/options [get]
func CommandOptionsGet(g *gin.Context) {
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

	cacheKey := getCommandOptionsDTOCacheKey(cmdID)
	cachedCommandOptions, ok := cache.Get(cacheKey, []CommandOptionDTO{})
	if ok {
		g.JSON(http.StatusOK, cachedCommandOptions)

		return
	}

	computedCommandOptionDTOs, err := computeCommandOptionDTOsForCommand(C.SlashCommandManager().GetCommands(), cmdID)
	if nil != err {
		webapi.RespondWithError(g, webapi.ErrorResponse{
			Status:    400,
			Error:     "Failed to find the desired command",
			Message:   "Could not find a command matching the given criteria",
			Timestamp: time.Now(),
		})
	}

	cache.Get(cacheKey, computedCommandOptionDTOs)

	g.JSON(http.StatusOK, computedCommandOptionDTOs)
}

// getCommandOptionsDTOCacheKey returns the cache key to get the CommandOptionDTO from cache
// for a specific command.
func getCommandOptionsDTOCacheKey(id string) string {
	return fmt.Sprintf(CommandOptionsDTOWebApiCacheKey, id)
}
