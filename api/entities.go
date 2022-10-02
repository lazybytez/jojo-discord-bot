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
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/services"
)

// GuildEntityManager is an entity manager
// that provides functionality for entities.Guild CRUD operations.
type GuildEntityManager interface {
	// Get tries to get a Guild from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no Guild can be found, the function returns a new empty
	// Guild.
	Get(guildId string) (*entities.Guild, error)
	// Count returns the number of all guilds stored in the database
	Count() (int64, error)

	// Create saves the passed Guild in the database.
	// Use Update or Save to update an already existing Guild.
	Create(guild *entities.Guild) error
	// Save updates the passed Guild in the database.
	// This does a generic update, use Update to do a precise and more performant update
	// of the entity when only updating a single field!
	Save(guild *entities.Guild) error
	// Update updates the defined field on the entity and saves it in the database.
	Update(guild *entities.Guild, column string, value interface{}) error
}

// RegisteredComponentEntityManager is an entity manager
// that provides functionality for entities.RegisteredComponent CRUD operations.
type RegisteredComponentEntityManager interface {
	// Get tries to get a RegisteredComponent from the
	// cache. If no cache entry is present, a request to the entities will be made.
	// If no RegisteredComponent can be found, the function returns a new empty
	// RegisteredComponent.
	Get(registeredComponentCode uint) (*entities.RegisteredComponent, error)
	// GetAvailable returns all components that have been registered
	// during application bootstrap.
	GetAvailable() []*entities.RegisteredComponent

	// Create saves the passed RegisteredComponent in the database.
	// Use Update or Save to update an already existing RegisteredComponent.
	Create(regComp *entities.RegisteredComponent) error
	// Save updates the passed RegisteredComponent in the database.
	// This does a generic update, use Update to do a precise and more performant update
	// of the entity when only updating a single field!
	Save(regComp *entities.RegisteredComponent) error
	// Update updates the defined field on the entity and saves it in the database.
	Update(regComp *entities.RegisteredComponent, column string, value interface{}) error
	// MarkAsAvailable marks the passed component as available, by putting
	// the codes into an array.
	// Note that duplicates will be filtered.
	MarkAsAvailable(code string)
}

// GlobalComponentStatusEntityManager is an entity manager
// that provides functionality for entities.GlobalComponentStatus CRUD operations.
type GlobalComponentStatusEntityManager interface {
	// Get tries to get a GlobalComponentStatus from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no GlobalComponentStatus can be found, the function returns a new empty
	// GlobalComponentStatus.
	Get(globalComponentStatusId uint) (*entities.GlobalComponentStatus, error)
	// GetDisplayString returns the string that indicates whether a component is
	// enabled or disabled globally. The string can directly being used to print
	// out messages in Discord.
	GetDisplayString(globalComponentStatusId uint) (string, error)

	// Create saves the passed GlobalComponentStatus in the database.
	// Use Update or Save to update an already existing GlobalComponentStatus.
	Create(globalComponentStatus *entities.GlobalComponentStatus) error
	// Save updates the passed GlobalComponentStatus in the database.
	// This does a generic update, use Update to do a precise and more performant update
	// of the entity when only updating a single field!
	Save(globalComponentStatus *entities.GlobalComponentStatus) error
	// Update updates the defined field on the entity and saves it in the database.
	Update(globalComponentStatus *entities.GlobalComponentStatus, column string, value interface{}) error
}

// GuildComponentStatusEntityManager is an entity manager
// that provides functionality for entities.GuildComponentStatus CRUD operations.
type GuildComponentStatusEntityManager interface {
	// Get tries to get a GuildComponentStatus from the
	// cache. If no cache entry is present, a request to the entities will be made.
	// If no GuildComponentStatus can be found, the function returns a new empty
	// GuildComponentStatus.
	Get(guildId uint, componentId uint) (*entities.Guild, error)
	// GetDisplay returns the status of a component in a form
	// that can be directly displayed in Discord.
	GetDisplay(guildId uint, componentId uint) (string, error)

	// Create saves the passed Guild in the database.
	// Use Update or Save to update an already existing Guild.
	Create(guildComponentStatus *entities.GuildComponentStatus) error
	// Save updates the passed Guild in the database.
	// This does a generic update, use Update to do a precise and more performant update
	// of the entity when only updating a single field!
	Save(guildComponentStatus *entities.GuildComponentStatus) error
	// Update updates the defined field on the entity and saves it in the database.
	Update(component *entities.GuildComponentStatus, column string, value interface{}) error

	// getComponentStatusCacheKey concatenates the passed guild and component ids to create
	// a new unique cache key for the component status
	getComponentStatusCacheKey(guildId uint, componentId uint) string
}

// EntityManager contains methods that allow to retrieve
// entity specific entity managers that provide caching and dedicated functions
// to work with specific entities.
// The entity specific managers allow to perform CRUD operations.
//
// Additionally, the main entity manager allows to register (auto migrate)
type EntityManager interface {
	// RegisterEntity registers a new entity (struct) and runs its automated
	// migration to ensure the entities schema is up-to-date.
	RegisterEntity(entityType interface{}) error
	// DB returns the current services.DatabaseAccess instance that
	// wraps gorm.DB and allows lower level database access.
	// Note that it is highly recommended to depend on methods of the
	// entity specific managers instead of services.DatabaseAccess.
	DB() services.DatabaseAccess

	// RegisteredComponent returns the RegisteredComponentEntityManager that is currently active,
	// which can be used to do RegisteredComponent specific entities actions.
	RegisteredComponent() *RegisteredComponentEntityManager
	// Guilds returns the GuildEntityManager that is currently active,
	// which can be used to do Guild specific entities actions.
	Guilds() *GuildEntityManager
	// GlobalComponentStatus returns the GlobalComponentStatusEntityManager that is currently active,
	// which can be used to do GlobalComponentStatus specific entities actions.
	GlobalComponentStatus() *GlobalComponentStatusEntityManager
	// GuildComponentStatus returns the GuildComponentStatusEntityManager that is currently active,
	// which can be used to do GuildComponentStatus specific entities actions.
	GuildComponentStatus() *GuildComponentStatusEntityManager
}

// entityManager is the internal database.GormDatabaseAccess instance
// used by the entities api and provided to components
// that need low-level access using gorm.DB.
var entityManager EntityManager

// EntityManager returns the currently active EntityManager.
// The currently active EntityManager is shared across all components.
//
// The entities.DatabaseAccess allows to interact with the entities of the application
// or access the raw gorm.DB instance, which is used for database access.
func (c *Component) EntityManager() EntityManager {
	return entityManager
}
