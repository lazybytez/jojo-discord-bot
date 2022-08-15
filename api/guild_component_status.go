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

package api

import (
	"fmt"
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"github.com/lazybytez/jojo-discord-bot/api/database"
	"gorm.io/gorm"
	"time"
)

var GuildComponentStatusEnabledDisplay = ":white_check_mark:"
var GuildComponentStatusDisabledDisplay = ":x:"

// ComponentStatus holds the status of a component on a specific server
type ComponentStatus struct {
	gorm.Model
	GuildID     uint
	Guild       database.Guild `gorm:"index:idx_guild_component;constraint:OnDelete:CASCADE;"`
	ComponentID uint
	Component   database.RegisteredComponent `gorm:"index:idx_guild_component;index:idx_component;constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// componentStatusCache is the cache used to reduce
// amount of database calls for the global component status.
var componentStatusCache = cache.New[string, ComponentStatus](10 * time.Minute)

// init ensure cache cleanup task is running
func init() {
	componentStatusCache.EnableAutoCleanup(10 * time.Minute)
}

// GetGuildComponentStatus tries to get a ComponentStatus from the
// cache. If no cache entry is present, a request to the database will be made.
// If no ComponentStatus can be found, the function returns a new empty
// ComponentStatus.
func GetGuildComponentStatus(c *Component, guildId uint, componentId uint) (*ComponentStatus, bool) {
	cacheKey := getComponentStatusCacheKey(guildId, componentId)
	comp, ok := cache.Get(componentStatusCache, cacheKey)

	if ok {
		return comp, true
	}

	regComp := &ComponentStatus{}
	queryStr := database.ColumnGuild + " = ? AND " + database.ColumnComponent + " = ?"
	ok = database.GetFirstEntity(regComp, queryStr, guildId, componentId)

	UpdateComponentStatus(c, guildId, componentId, regComp)

	return regComp, ok
}

// GetGuildComponentStatusDisplay returns the status of a component in a form
// that can be directly displayed in Discord.
func GetGuildComponentStatusDisplay(c *Component, guildId uint, componentId uint) (string, bool) {
	compState, ok := GetGuildComponentStatus(c, guildId, componentId)
	if !ok {
		return GuildComponentStatusDisabledDisplay, false
	}

	if compState.Enabled {
		return GuildComponentStatusEnabledDisplay, true
	}

	return GuildComponentStatusDisabledDisplay, false
}

// UpdateComponentStatus adds or updates a cached item in the ComponentStatus cache.
func UpdateComponentStatus(_ *Component, guildId uint, componentId uint, component *ComponentStatus) {
	cache.Update(componentStatusCache, getComponentStatusCacheKey(guildId, componentId), component)
}

// getComponentStatusCacheKey concatenates the passed guild and component ids to create
// a new unique cache key for the component status
func getComponentStatusCacheKey(guildId uint, componentId uint) string {
	return fmt.Sprintf("%v_%v", guildId, componentId)
}
