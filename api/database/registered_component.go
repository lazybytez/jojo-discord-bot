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

// RegisteredComponentEntityManager is the RegisteredComponent specific entity manager
// that allows easy access to global component status in the database.
type RegisteredComponentEntityManager struct {
	em    *EntityManager
	cache *cache.Cache[uint, RegisteredComponent]
}

// RegisteredComponentDBAccess is the interface that defines the capabilities
// of the RegisteredComponentEntityManager
type RegisteredComponentDBAccess interface {
	// Get tries to get a RegisteredComponent from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no RegisteredComponent can be found, the function returns a new empty
	// RegisteredComponent.
	Get(registeredComponentCode uint) (*RegisteredComponent, error)
	// Update adds or updates a cached item in the RegisteredComponentEntityManager cache.
	Update(registeredComponentId uint, registeredComponent *RegisteredComponent) error
}

// RegisteredComponent returns the RegisteredComponentEntityManager that is currently active,
// which can be used to do RegisteredComponent specific database actions.
func (em *EntityManager) RegisteredComponent() *RegisteredComponentEntityManager {
	if nil != em.entityManagers.guild {
		rgem := &RegisteredComponentEntityManager{
			em,
			cache.New[uint, RegisteredComponent](10 * time.Minute),
		}
		em.entityManagers.registeredComponentEntityManager = rgem
	}

	return em.entityManagers.registeredComponentEntityManager
}

// registeredComponentCache is the used instance of the registeredComponentCacheContainer
// that allows caching of RegisteredComponent
var registeredComponentCache = cache.New[string, RegisteredComponent](0)

// Get tries to get a RegisteredComponent from the
// cache. If no cache entry is present, a request to the database will be made.
// If no RegisteredComponent can be found, the function returns a new empty
// RegisteredComponent.
func (rgem *RegisteredComponentEntityManager) Get(registeredComponentCode string) (*RegisteredComponent, error) {
	comp, ok := cache.Get(registeredComponentCache, registeredComponentCode)

	if ok {
		return comp, nil
	}

	regComp := &RegisteredComponent{}
	err := rgem.em.GetFirstEntity(regComp, "code = ?", registeredComponentCode)
	if nil != err {
		return regComp, err
	}

	UpdateRegisteredComponent(registeredComponentCode, regComp)

	return regComp, err
}

// UpdateRegisteredComponent adds or updates a cached item in the RegisteredComponent cache.
func UpdateRegisteredComponent(code string, component *RegisteredComponent) {
	cache.Update(registeredComponentCache, code, component)
}
