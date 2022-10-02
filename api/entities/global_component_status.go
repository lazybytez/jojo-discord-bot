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

package entities

import (
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
	ComponentID uint                `gorm:"index:idx_global_component_status_component_id;"`
	Component   RegisteredComponent `gorm:"constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// GlobalComponentStatusEntityManager is the GlobalComponentStatus specific entity manager
// that allows easy access to global component status in the entities.
type GlobalComponentStatusEntityManager struct {
	*EntityManager
	cache *cache.Cache[uint, GlobalComponentStatus]
}

// GlobalComponentStatus returns the GlobalComponentStatusEntityManager that is currently active,
// which can be used to do GlobalComponentStatus specific entities actions.
func (em *EntityManager) GlobalComponentStatus() *GlobalComponentStatusEntityManager {
	if nil == em.globalComponentStatusEntityManager {
		gem := &GlobalComponentStatusEntityManager{
			em,
			cache.New[uint, GlobalComponentStatus](10 * time.Minute),
		}
		em.globalComponentStatusEntityManager = gem

		err := gem.cache.EnableAutoCleanup(10 * time.Minute)
		if nil != err {
			em.logger.Err(err, "Failed to initialize periodic cache cleanup task "+
				"for GlobalComponmentStatus entity manager!")
		}
	}

	return em.globalComponentStatusEntityManager
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
// cache. If no cache entry is present, a request to the entities will be made.
// If no GlobalComponentStatus can be found, the function returns a new empty
// GlobalComponentStatus.
func (gem *GlobalComponentStatusEntityManager) Get(globalComponentStatusId uint) (*GlobalComponentStatus, error) {
	comp, ok := cache.Get(gem.cache, globalComponentStatusId)

	if ok {
		return comp, nil
	}

	globalCompStatus := &GlobalComponentStatus{}
	err := gem.database.GetFirstEntity(globalCompStatus, ColumnComponent+" = ?", globalComponentStatusId)
	if nil != err {
		return globalCompStatus, err
	}

	cache.Update(gem.cache, globalCompStatus.ID, globalCompStatus)

	return globalCompStatus, err
}

// Create saves the passed GlobalComponentStatus in the database.
// Use Update or Save to update an already existing GlobalComponentStatus.
func (gem *GlobalComponentStatusEntityManager) Create(globalComponentStatus *GlobalComponentStatus) error {
	err := gem.database.Create(globalComponentStatus)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(gem.cache, globalComponentStatus.ID, globalComponentStatus)

	return nil
}

// Save updates the passed GlobalComponentStatus in the database.
// This does a generic update, use Update to do a precise and more performant update
// of the entity when only updating a single field!
func (gem *GlobalComponentStatusEntityManager) Save(globalComponentStatus *GlobalComponentStatus) error {
	err := gem.database.Save(globalComponentStatus)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(gem.cache, globalComponentStatus.ID, globalComponentStatus)

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (gem *GlobalComponentStatusEntityManager) Update(
	globalComponentStatus *GlobalComponentStatus,
	column string,
	value interface{},
) error {
	err := gem.database.UpdateEntity(globalComponentStatus, column, value)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(gem.cache, globalComponentStatus.ID, globalComponentStatus)

	return nil
}
