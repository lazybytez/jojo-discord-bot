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
	"github.com/lazybytez/jojo-discord-bot/api"
	"net/http"
)

// ComponentDTO is an intermediate data transfer object
// that can be output or received by the WebAPI.
// This type is used, because the bot both has the general api.Component
// and the entities.RegisteredComponent types. This type represents
// a component with only API relevant data.
type ComponentDTO struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	GlobalEnabled bool   `json:"global_enabled"`
	GuildEnabled  bool   `json:"guild_enabled"`
}

// ComponentDTOFromComponent creates a new ComponentDTO.
// The general data is taken from the passed api.Component.
// The enabled states are pulled from the database.
//
// Note that the guild status will only being pulled when a valid GuildID is passed.
// The passed GuildID must be present in the database.
// If no GuildID is passed or the ID is invalid, the Guild enabled status will result in false.
func ComponentDTOFromComponent(c *api.Component, guildId string) (ComponentDTO, error) {
	regComp, err := C.EntityManager().RegisteredComponent().Get(c.Code)
	if nil != err {
		return ComponentDTO{}, err
	}

	gcs, err := C.EntityManager().GlobalComponentStatus().Get(regComp.ID)
	if nil != err {
		return ComponentDTO{}, err
	}

	guildComponentStatus := false
	if guildId <= "" {
		if guild, err := C.EntityManager().Guilds().Get(guildId); nil == err {
			if guildCSE, err := C.EntityManager().GuildComponentStatus().Get(guild.ID, regComp.ID); nil == err {
				guildComponentStatus = guildCSE.Enabled
			}
		}
	}

	return ComponentDTO{
		Code:          c.Code,
		Name:          c.Name,
		Description:   c.Description,
		GlobalEnabled: gcs.Enabled,
		GuildEnabled:  guildComponentStatus,
	}, nil
}

// ComponentsGet endpoint
// @Summary Endpoint that returns all available components
// @Schemes
// @Description This endpoint collects all available components and returns them.
// @Tags core
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /v1/components [get]
func ComponentsGet(g *gin.Context) {
	componentDTOs := make([]ComponentDTO, len(api.Components))
	var err error

	for i, comp := range api.Components {
		componentDTOs[i], err = ComponentDTOFromComponent(comp, "")
		if nil != err {
			C.Logger().Err(err, "Failed to convert component with code \"%s\" to ComponentDTO!", comp.Code)
		}
	}

	g.JSON(http.StatusOK, componentDTOs)
}
