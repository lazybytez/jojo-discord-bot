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

package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/lazybytez/jojo-discord-bot/build"
	"github.com/lazybytez/jojo-discord-bot/webapi"
	"net/http"
)

// StatsDTO is a DTO used to store statistics that can be
// output in the web api.
//
// @Description Statistics holds statistics about the bot like the current version.
type StatsDTO struct {
	GuildCount        int    `json:"guild_count"`
	SlashCommandCount int    `json:"slash_command_count"`
	Version           string `json:"version"`
} //@Name Statistics

// StatsGet endpoint
//
// @Summary     Get general statistics about the bot
// @Description This endpoint collects publicly available statistics and returns them.
// @Description The result contains general information like the current count of guilds or version.
// @Description Note that some statistics might be cached and only updated every few minutes.
// @Tags        General
// @Produce     json
// @Success     200 {object} StatsDTO "The statistics collected by the bot"
// @Failure		500 {object} webapi.ErrorResponse "An error indicating that an internal error happened"
// @Router      /stats [get]
func StatsGet(g *gin.Context) {
	statsDto := StatsDTO{
		GuildCount:        collectGuildCount(),
		SlashCommandCount: collectSlashCommandCount(),
		Version:           build.ComputeVersionString(),
	}

	g.JSON(http.StatusOK, statsDto)
}

// registerRoutes registers the web api routes
// provided by the statistics component.
func registerRoutes() {
	eg := webapi.Router().Group("/stats")
	eg.GET("/", StatsGet)
}
