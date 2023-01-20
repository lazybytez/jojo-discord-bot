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
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommandDTOsWebApiCacheKey is the cache key used to store and retrieve all commands
// as CommandDTO instances from the cache.
const CommandDTOsWebApiCacheKey = "bot_web_api_commands_get_cache"

// CommandsGet endpoint
//
// @Summary     Get all available commands of the bot
// @Description This endpoint collects all available commands and returns them.
// @Description The result on a success contains relevant information like name, description and category of commands.
// @Description Note that this endpoint does not return detailed information like the options of a command.
// @Description To obtain the available command options, the command must be queried on its own using the
// @Description single command options get endpoint.
// @Tags        Command System
// @Produce     json
// @Success     200 {array} CommandDTO "An array consisting of objects containing information about commands"
// @Failure		500 {object} webapi.ErrorResponse "An error indicating that an internal error happened"
// @Router      /commands [get]
func CommandsGet(g *gin.Context) {
	commandDTOs := getCommandDTOs()

	g.JSON(http.StatusOK, commandDTOs)
}
