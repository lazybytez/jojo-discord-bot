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
	"github.com/lazybytez/jojo-discord-bot/services"
)

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
	// Logger returns the services.Logger of the EntityManager.
	Logger() services.Logger
}
