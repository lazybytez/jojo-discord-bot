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
	"strings"
	"time"
)

// CoreComponentPrefix is the prefix put in front of components that
// cannot be managed by server owners, as they are important core components
const CoreComponentPrefix = "bot_"

// RegisteredComponentEntityManager is the RegisteredComponent specific entity manager
// that allows easy access to global component status in the database.
type RegisteredComponentEntityManager struct {
	em                  *EntityManager
	cache               *cache.Cache[string, RegisteredComponent]
	availableComponents []string
}

// RegisteredComponentDBAccess is the interface that defines the capabilities
// of the RegisteredComponentEntityManager
type RegisteredComponentDBAccess interface {
	// Get tries to get a RegisteredComponent from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no RegisteredComponent can be found, the function returns a new empty
	// RegisteredComponent.
	Get(registeredComponentCode uint) (*RegisteredComponent, error)
	// GetAvailable returns all components that have been registered
	// during application bootstrap.
	GetAvailable() []*RegisteredComponent
	// Update adds or updates a cached item in the RegisteredComponentEntityManager cache.
	Update(registeredComponentId uint, registeredComponent *RegisteredComponent) error
	// MarkAsAvailable marks the passed component as available, by putting
	// the codes into an array.
	// Note that duplicates will be filtered.
	MarkAsAvailable(code string)
}

// RegisteredComponent returns the RegisteredComponentEntityManager that is currently active,
// which can be used to do RegisteredComponent specific database actions.
func (em *EntityManager) RegisteredComponent() *RegisteredComponentEntityManager {
	if nil == em.entityManagers.registeredComponentEntityManager {
		rgem := &RegisteredComponentEntityManager{
			em,
			cache.New[string, RegisteredComponent](10 * time.Minute),
			make([]string, 0),
		}
		em.entityManagers.registeredComponentEntityManager = rgem
	}

	return em.entityManagers.registeredComponentEntityManager
}

// Get tries to get a RegisteredComponent from the
// cache. If no cache entry is present, a request to the database will be made.
// If no RegisteredComponent can be found, the function returns a new empty
// RegisteredComponent.
func (rgem *RegisteredComponentEntityManager) Get(registeredComponentCode string) (*RegisteredComponent, error) {
	comp, ok := cache.Get(rgem.cache, registeredComponentCode)

	if ok {
		return comp, nil
	}

	regComp := &RegisteredComponent{}
	err := rgem.em.GetFirstEntity(regComp, "code = ?", registeredComponentCode)
	if nil != err {
		return regComp, err
	}

	rgem.Update(registeredComponentCode, regComp)

	return regComp, err
}

// GetAvailable returns all components that have been registered
// during application bootstrap.
func (rgem *RegisteredComponentEntityManager) GetAvailable() []*RegisteredComponent {
	availableComponents := make([]*RegisteredComponent, 0)

	for _, code := range rgem.availableComponents {
		regComp, err := rgem.Get(code)
		if nil != err {
			continue
		}

		availableComponents = append(availableComponents, regComp)
	}

	return availableComponents
}

// Update adds or updates a cached item in the RegisteredComponent cache.
func (rgem *RegisteredComponentEntityManager) Update(code string, component *RegisteredComponent) {
	cache.Update(rgem.cache, code, component)
}

// MarkAsAvailable marks the passed component as available, by putting
// the codes into an array.
// Note that duplicates will be filtered.
func (rgem *RegisteredComponentEntityManager) MarkAsAvailable(code string) {
	for _, comp := range rgem.availableComponents {
		if code == comp {
			return
		}
	}

	rgem.availableComponents = append(rgem.availableComponents, code)
}

// CoreComponentChecker is an interface providing functionality to
// check if a registered component is a core component.
type CoreComponentChecker interface {
	// IsCoreComponent checks whether the passed RegisteredComponent is a core
	// component or not.
	//
	// Core components are components which are prefixed with the CoreComponentPrefix.
	IsCoreComponent() bool
}

// IsCoreComponent checks whether the passed RegisteredComponent is a core
// component or not.
//
// Core components are components which are prefixed with the CoreComponentPrefix.
func (regComp *RegisteredComponent) IsCoreComponent() bool {
	return strings.HasPrefix(regComp.Code, CoreComponentPrefix)
}
