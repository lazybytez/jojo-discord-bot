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
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"gorm.io/gorm"
	"strings"
)

// CoreComponentPrefix is the prefix put in front of components that
// cannot be managed by server owners, as they are important core components
const CoreComponentPrefix = "bot_"

// RegisteredComponent represents a single component that is or was known
// to the system.
//
// Single purpose of this struct is to provide a entities
// table with which relations can be build to ensure integrity
// of the GuildComponentStatus and GlobalComponentStatus tables.
type RegisteredComponent struct {
	gorm.Model
	Code           string `gorm:"uniqueIndex"`
	Name           string
	Description    string
	DefaultEnabled bool
}

// RegisteredComponentEntityManager is the RegisteredComponent specific entity manager
// that allows easy access to global component status in the database.
type RegisteredComponentEntityManager struct {
	EntityManager

	availableComponents []string
}

// NewRegisteredComponentEntityManager creates a new RegisteredComponentEntityManager.
func NewRegisteredComponentEntityManager(entityManager EntityManager) *RegisteredComponentEntityManager {
	return &RegisteredComponentEntityManager{
		entityManager,
		make([]string, 0),
	}
}

// Get tries to get a RegisteredComponent from the
// cache. If no cache entry is present, a request to the database will be made.
// If no RegisteredComponent can be found, the function returns a new empty
// RegisteredComponent.
func (rgem *RegisteredComponentEntityManager) Get(registeredComponentCode string) (*RegisteredComponent, error) {
	cacheKey := rgem.getCacheKey(registeredComponentCode)
	cachedComp, ok := cache.Get(cacheKey, RegisteredComponent{})

	if ok {
		return &cachedComp, nil
	}

	regComp := &RegisteredComponent{}
	err := rgem.DB().GetFirstEntity(regComp, "code = ?", registeredComponentCode)
	if nil != err {
		return regComp, err
	}

	cache.Update(cacheKey, *regComp)

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

// Create saves the passed RegisteredComponent in the database.
// Use Update or Save to update an already existing RegisteredComponent.
func (rgem *RegisteredComponentEntityManager) Create(regComp *RegisteredComponent) error {
	err := rgem.DB().Create(regComp)
	if nil != err {
		return err
	}

	// Invalidate cache item (if present)
	cacheKey := rgem.getCacheKey(regComp.Code)
	cache.Invalidate(cacheKey, GlobalComponentStatus{})

	return nil
}

// Save updates the passed RegisteredComponent in the database.
// This does a generic update, use Update to do a precise and more performant update
// of the entity when only updating a single field!
func (rgem *RegisteredComponentEntityManager) Save(regComp *RegisteredComponent) error {
	err := rgem.DB().Save(regComp)
	if nil != err {
		return err
	}

	// Invalidate cache item (if present)
	cacheKey := rgem.getCacheKey(regComp.Code)
	cache.Invalidate(cacheKey, GlobalComponentStatus{})

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (rgem *RegisteredComponentEntityManager) Update(regComp *RegisteredComponent, column string, value interface{}) error {
	err := rgem.DB().UpdateEntity(regComp, column, value)
	if nil != err {
		return err
	}

	// Invalidate cache item (if present)
	cacheKey := rgem.getCacheKey(regComp.Code)
	cache.Invalidate(cacheKey, GlobalComponentStatus{})

	return nil
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

// IsCoreComponent checks whether the passed RegisteredComponent is a core
// component or not.
//
// Core components are components which are prefixed with the CoreComponentPrefix.
func (regComp *RegisteredComponent) IsCoreComponent() bool {
	return strings.HasPrefix(regComp.Code, CoreComponentPrefix)
}

// getCacheKey returns the computed cache key used to cache
// RegisteredComponent objects.
func (rgem *RegisteredComponentEntityManager) getCacheKey(registeredComponentCode string) string {
	return registeredComponentCode
}
