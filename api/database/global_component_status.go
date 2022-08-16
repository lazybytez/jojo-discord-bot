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

package database

import (
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"gorm.io/gorm"
	"time"
)

const GlobalComponentStatusEnabledDisplay = ":white_check_mark:"
const GlobalComponentStatusDisabledDisplay = ":no_entry:"

// GlobalComponentStatus holds the status of a component in the global context.
// This allows disabling a bugging component globally if necessary.
type GlobalComponentStatus struct {
	gorm.Model
	ComponentID uint
	Component   RegisteredComponent `gorm:"index:idx_guild_component;index:idx_component;constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// globalComponentStatusCache is the cache used to reduce
// amount of database calls for the global component status.
var globalComponentStatusCache = cache.New[uint, GlobalComponentStatus](10 * time.Minute)

// GetGlobalStatusDisplayString returns the string that indicates whether a component is
// enabled or disabled globally. The string can directly being used to print
// out messages in Discord.
func GetGlobalStatusDisplayString(c *api.Component, registeredComponentId uint) (string, bool) {
	compState, ok := GetGlobalComponentStatus(c, registeredComponentId)
	if !ok {
		return GlobalComponentStatusDisabledDisplay, false
	}

	if compState.Enabled {
		return GlobalComponentStatusEnabledDisplay, true
	}

	return GlobalComponentStatusDisabledDisplay, false
}

// GetGlobalComponentStatus tries to get a GlobalComponentStatus from the
// cache. If no cache entry is present, a request to the database will be made.
// If no GlobalComponentStatus can be found, the function returns a new empty
// GlobalComponentStatus.
func GetGlobalComponentStatus(c *api.Component, registeredComponentId uint) (*GlobalComponentStatus, bool) {
	comp, ok := cache.Get(globalComponentStatusCache, registeredComponentId)

	if ok {
		return comp, true
	}

	regComp := &GlobalComponentStatus{}
	ok = GetFirstEntity(c, regComp, ColumnComponent+" = ?", registeredComponentId)

	UpdateGlobalComponentStatus(c, registeredComponentId, regComp)

	return regComp, ok
}

// UpdateGlobalComponentStatus adds or updates a cached item in the GlobalComponentStatus cache.
func UpdateGlobalComponentStatus(_ *api.Component, registeredComponentId uint, component *GlobalComponentStatus) {
	cache.Update(globalComponentStatusCache, registeredComponentId, component)
}
