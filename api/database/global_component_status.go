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
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"time"
)

const GlobalComponentStatusEnabledDisplay = ":white_check_mark:"
const GlobalComponentStatusDisabledDisplay = ":no_entry:"

// GlobalComponentStatusEntityManager is the GlobalComponentStatus specific entity manager
// that allows easy access to global component status in the database.
type GlobalComponentStatusEntityManager struct {
	em    *EntityManager
	cache *cache.Cache[uint, GlobalComponentStatus]
}

// GlobalComponentStatusDBAccess is the interface that defines the capabilities
// of the GlobalComponentStatusEntityManager
type GlobalComponentStatusDBAccess interface {
	// Get tries to get a GlobalComponentStatus from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no GlobalComponentStatus can be found, the function returns a new empty
	// GlobalComponentStatus.
	Get(globalComponentStatusId uint) (*GlobalComponentStatus, error)
	// GetDisplayString returns the string that indicates whether a component is
	// enabled or disabled globally. The string can directly being used to print
	// out messages in Discord.
	GetDisplayString(globalComponentStatusId uint) (string, error)
	// Update adds or updates a cached item in the GlobalComponentStatusDBAccess cache.
	Update(globalComponentStatusId string, globalComponentStatus *GlobalComponentStatus) error
}

// GlobalComponentStatus returns the GlobalComponentStatusEntityManager that is currently active,
// which can be used to do GlobalComponentStatus specific database actions.
func (em *EntityManager) GlobalComponentStatus() *GlobalComponentStatusEntityManager {
	if nil != em.entityManagers.guild {
		gem := &GlobalComponentStatusEntityManager{
			em,
			cache.New[uint, GlobalComponentStatus](10 * time.Minute),
		}
		em.entityManagers.globalComponentStatusEntityManager = gem

		gem.cache.EnableAutoCleanup(10 * time.Minute)
	}

	return em.entityManagers.globalComponentStatusEntityManager
}

// GetDisplayString returns the string that indicates whether a component is
// enabled or disabled globally. The string can directly being used to print
// out messages in Discord.
func (gem *GlobalComponentStatusEntityManager) GetDisplayString(globalComponentStatusId uint) (string, error) {
	compState, err := gem.Get(globalComponentStatusId)
	if nil != err {
		return GlobalComponentStatusDisabledDisplay, err
	}

	if compState.Enabled {
		return GlobalComponentStatusEnabledDisplay, nil
	}

	return GlobalComponentStatusDisabledDisplay, nil
}

// Get tries to get a GlobalComponentStatus from the
// cache. If no cache entry is present, a request to the database will be made.
// If no GlobalComponentStatus can be found, the function returns a new empty
// GlobalComponentStatus.
func (gem *GlobalComponentStatusEntityManager) Get(globalComponentStatusId uint) (*GlobalComponentStatus, error) {
	comp, ok := cache.Get(gem.cache, globalComponentStatusId)

	if ok {
		return comp, nil
	}

	globalCompStatus := &GlobalComponentStatus{}
	err := gem.em.GetFirstEntity(globalCompStatus, ColumnComponent+" = ?", globalComponentStatusId)
	if nil != err {
		return globalCompStatus, err
	}

	gem.Update(globalComponentStatusId, globalCompStatus)

	return globalCompStatus, err
}

// Update adds or updates a cached item in the GlobalComponentStatus cache.
func (gem *GlobalComponentStatusEntityManager) Update(
	globalComponentStatusId uint,
	globalComponentStatus *GlobalComponentStatus,
) {
	cache.Update(gem.cache, globalComponentStatusId, globalComponentStatus)
}
