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
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"github.com/lazybytez/jojo-discord-bot/webapi"
	"net/http"
	"time"
)

// ComponentDTOsResponseWebApiCacheKey is the cache key used to store and retrieve all components
// as ComponentDTO instances from the cache.
const ComponentDTOsResponseWebApiCacheKey = "bot_web_api_components_get_cache"

// ComponentDTO is an intermediate data transfer object
// that can be output or received by the WebAPI.
// This type is used, because the bot both has the general api.Component
// and the entities.RegisteredComponent types. This type represents
// a component with only API relevant data.
//
// @Description Component holds the metadata of a component like its name and
// @Description and description. Additionally, it holds the current global
// @Description and guild status.
type ComponentDTO struct {
	Code          entities.ComponentCode `json:"code"`
	Name          string                 `json:"name"`
	Categories    api.Categories         `json:"categories"`
	Description   string                 `json:"description"`
	GlobalEnabled bool                   `json:"global_enabled"`
	GuildEnabled  bool                   `json:"guild_enabled"`
} //@Name Component

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
		Categories:    c.Categories,
		Description:   c.Description,
		GlobalEnabled: gcs.Enabled,
		GuildEnabled:  guildComponentStatus,
	}, nil
}

// ComponentsGet endpoint
//
// @Summary     Get all available components of the bot
// @Description This endpoint collects all available components and returns them.
// @Description The result on a success contains all relevant information, which includes name and description of components.
// @Description Additionally, the endpoint also returns the status of the components.
// @Description
// @Description The guild status is currently not populated and always false!
// @Tags        Component System
// @Produce     json
// @Success     200 {array} ComponentDTO "An array consisting of objects containing information about components"
// @Failure		500 {object} webapi.ErrorResponse "An error indicating that an internal error happened"
// @Router      /components [get]
func ComponentsGet(g *gin.Context) {
	cachedResponse, ok := cache.Get(ComponentDTOsResponseWebApiCacheKey, []ComponentDTO{})
	if ok {
		g.JSON(http.StatusOK, cachedResponse)

		return
	}

	componentDTOs := make([]ComponentDTO, len(api.Components))
	var err error

	for i, comp := range api.Components {
		componentDTOs[i], err = ComponentDTOFromComponent(comp, "")
		if nil != err {
			C.Logger().Err(err, "Failed to convert component with code \"%s\" to ComponentDTO!", comp.Code)

			webapi.RespondWithError(g, webapi.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "Failed to prepare components",
				Message: fmt.Sprintf(
					"The server failed to prepare the component \"%s\"",
					comp.Name),
				Timestamp: time.Now(),
			})
		}
	}

	err = cache.Update(ComponentDTOsResponseWebApiCacheKey, componentDTOs)
	if nil != err {
		C.Logger().Err(err, "Failed to cache ComponentDTOs for GET components web api endpoint!")
	}

	g.JSON(http.StatusOK, componentDTOs)
}
