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
	"gorm.io/gorm"
	"sync"
)

const ColumnName = "code"

// registeredComponentCache is a simple struct that
// is used to create cache instances that allow to cache
// RegisteredComponent
type registeredComponentCacheContainer struct {
	// registeredComponentCache is a map that holds component
	// codes and their reference.
	//
	// It acts as a cache to prevent exhaustive database calls
	// for something that is required during the entire application
	// lifetime and does not change.
	cache map[string]*RegisteredComponent
	lock  sync.RWMutex
}

// registeredComponentCache is the used instance of the registeredComponentCacheContainer
// that allows caching of RegisteredComponent
var registeredComponentCache = &registeredComponentCacheContainer{
	map[string]*RegisteredComponent{},
	sync.RWMutex{},
}

// RegisteredComponent represents a single component that is or was known
// to the system.
//
// Single purpose of this struct is to provide a database
// table with which relations can be build to ensure integrity
// of the ComponentStatus and GlobalComponentStatus tables.
type RegisteredComponent struct {
	gorm.Model
	Code string `gorm:"uniqueIndex"`
}

// GetRegisteredComponent tries to get a RegisteredComponent from the
// cache. If no cache entry is present, a request to the database will be made.
// If no RegisteredComponent can be found, the function returns a new empty
// RegisteredComponent.
func GetRegisteredComponent(c *api.Component, code string) (*RegisteredComponent, bool) {
	registeredComponentCache.lock.RLock()
	comp, ok := registeredComponentCache.cache[code]
	registeredComponentCache.lock.RUnlock()

	if ok {
		return comp, true
	}

	regComp := &RegisteredComponent{}
	ok = GetFirstEntity(c, regComp, "code = ?", code)
	if !ok {
		UpdateCache(code, regComp)

		return regComp, false
	}

	UpdateCache(code, regComp)

	return regComp, true
}

// UpdateCache adds or updates a cached item in the RegisteredComponent cache.
func UpdateCache(code string, component *RegisteredComponent) {
	registeredComponentCache.lock.Lock()
	defer registeredComponentCache.lock.Unlock()

	registeredComponentCache.cache[code] = component
}
