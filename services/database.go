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

package services

import "gorm.io/gorm"

// DatabaseAccess is an interface thar provides functions for low level
// entities operations that bypass the entity manager API.
//
// This is part of the low-level API. Prefer using the entity API
// instead of this one!
//
// Unlike other parts of the API, this is a direct wrapper around gorm.DB
// and forces the application to redirect all calls through the API.
// Also, all entities of the application are based on GORMs definition.
// This means it is not possible to just swap the used ORM implementation,
// like it would be for other parts of the API.
//
// This interface exists to provide a unified and easy way to access the
// GORM functions and gorm.DB itself. A nice side effect is, that the wrapped functions
// have improved error handling by automatically passing down errors.
type DatabaseAccess interface {
	// RegisterEntity registers a new entity (struct) and ensures the
	// entities schema knows about the entity.
	RegisterEntity(entityType interface{}) error

	// Create creates the passed entity in the entities
	Create(entity interface{}) error
	// Save upserts the passed entity in the entities
	Save(entity interface{}) error
	// UpdateEntity can be used to update the passed entity in the entities
	UpdateEntity(entityContainer interface{}, column string, value interface{}) error
	// DeleteEntity deletes the passed entity from the entities.
	DeleteEntity(entityContainer interface{}) error

	// GetFirstEntity fills the passed entity container with the first
	// found entity matching the passed conditions.
	//
	// Returns an error if no record could be found.
	GetFirstEntity(entityContainer interface{}, conditions ...interface{}) error
	// GetLastEntity fills the passed entity container with the last
	// found entity matching the passed conditions.
	//
	// Returns false if no entries could be found.
	GetLastEntity(entityContainer interface{}, conditions ...interface{}) error
	// GetEntities fills the passed entities slice with the entities
	// that have been found for the specified condition.
	GetEntities(entities interface{}, conditions ...interface{}) error

	// WorkOn returns a gorm.DB pointer that allows to do a custom search or actions on entities.
	//
	// The returned gorm.DB instance is created by using gorm.DB.Model() and is therefore
	// already prepared to get started with applying filters.
	// This function is the only interface point to get direct access to gorm.DB
	WorkOn(entityContainer interface{}) *gorm.DB

	// DB returns the reference to the used gorm.DB instance.
	DB() *gorm.DB
}
