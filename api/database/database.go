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
	"github.com/lazybytez/jojo-discord-bot/api/log"
	"github.com/lazybytez/jojo-discord-bot/api/util"
	"gorm.io/gorm"
)

// Due to generics, the decision has been made to
// run the database management as simple functions
// instead of creating an interface like it has been done for
// event handlers and slash commands.
// This might change in the future.
//
// We still wrap all GORM functions to:
//  a) provide a unified API
//  b) keep the ability open to intercept actions and add custom logic
//
// Most functions expect parameters in a similar way they would be passed
// to GORM.

// loggerPrefix is the prefix used for logging
// when no component is associated with an action.
const loggerPrefix = "database_api"

// entityManager is the DBAccess used during application
// lifetime. It is initialized using the Init function.
var entityManager EntityManager

// EntityManager is a container struct that holds the
// gorm.DB instance to use for database interaction.
//
// When calling DBAccess() on a api.Component,
// in reality this type is returned internally.
type EntityManager struct {
	database       *gorm.DB
	logger         log.Logging
	entityManagers entityManagers
}

// GetEntityManager returns the currently active DBAccess.
//
// When needing database access in a component, use the components
// GetEntityManager() method instead! This function is meant for use by the
// API and Core
func GetEntityManager() *EntityManager {
	return &entityManager
}

// DBAccess is a wrapper around gorm.DB and forces the application
// to redirect all calls through the API.
//
// Reason for this is to keep control on the API up to a specific degree.
// Also, it allows to provide a unified way of accessing the database
// that perfectly suites the applications structure.
type DBAccess interface {
	// RegisterEntity registers a new entity (struct) and runs its automated
	// migration to ensure the database schema is up-to-date.
	RegisterEntity(entityType interface{}) error

	// Create creates the passed entity in the database
	Create(entity interface{}) error
	// Save upserts the passed entity in the database
	Save(entity interface{}) error
	// UpdateEntity can be used to update the passed entity in the database
	UpdateEntity(entityContainer interface{}, column string, value interface{}) error

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
	GetEntities(entities []interface{}, conditions ...interface{}) error

	// WorkOn returns a gorm.DB pointer that allows to do a custom search or actions on entities.
	//
	// The returned gorm.DB instance is created by using gorm.DB.Model() and is therefore
	// already prepared to get started with applying filters.
	// This function is the only interface point to get direct access to gorm.DB
	WorkOn(entityContainer interface{}) *gorm.DB

	// entitySpecificManagerAccess provides methods to work with discrete entities
	entitySpecificManagerAccess
}

// Init database instance for the database API
// and register default entity types that are managed by
// the bot core.
func Init(db *gorm.DB) error {
	entityManager = EntityManager{
		db,
		log.New(loggerPrefix, nil),
		entityManagers{},
	}

	return registerDefaultEntities(&entityManager)
}

// RegisterEntity registers a new entity (struct) and runs its automated
// migration to ensure the database schema is up-to-date.
func (em *EntityManager) RegisterEntity(entityType interface{}) error {
	err := em.database.AutoMigrate(entityType)
	if nil != err {
		em.logger.Err(err, "Failed to auto-migrated entity \"%v\"", util.ExtractTypeName(entityType))

		return err
	}

	em.logger.Info("Auto-migrated entity \"%v\"", util.ExtractTypeName(entityType))

	return nil
}

// GetFirstEntity fills the passed entity container with the first
// found entity matching the passed conditions.
//
// Returns an error if no record could be found.
func (em *EntityManager) GetFirstEntity(entityContainer interface{}, conditions ...interface{}) error {
	return em.database.First(entityContainer, conditions...).Error
}

// GetLastEntity fills the passed entity container with the last
// found entity matching the passed conditions.
//
// Returns false if no entries could be found.
func (em *EntityManager) GetLastEntity(entityContainer interface{}, conditions ...interface{}) error {
	return em.database.Last(entityContainer, conditions...).Error
}

// GetEntities fills the passed entities slice with the entities
// that have been found for the specified condition.
func (em *EntityManager) GetEntities(entities []interface{}, conditions ...interface{}) error {
	return em.database.Find(entities, conditions...).Error
}

// Create creates the passed entity in the database
func (em *EntityManager) Create(entity interface{}) error {
	return em.database.Create(entity).Error
}

// Save upserts the passed entity in the database
func (em *EntityManager) Save(entity interface{}) error {
	return em.database.Save(entity).Error
}

// UpdateEntity can be used to update the passed entity in the database
func (em *EntityManager) UpdateEntity(entityContainer interface{}, column string, value interface{}) error {
	return em.database.Model(entityContainer).Update(column, value).Error
}

// DeleteEntity deletes the passed entity from the database.
func (em *EntityManager) DeleteEntity(entityContainer interface{}) error {
	return em.database.Delete(entityContainer).Error
}

// WorkOn returns a gorm.DB pointer that allows to do a custom search or actions on entities.
//
// The returned gorm.DB instance is created by using gorm.DB.Model() and is therefore
// already prepared to get started with applying filters.
// This function is the only interface point to get direct access to gorm.DB
func (em *EntityManager) WorkOn(entityContainer interface{}) *gorm.DB {
	return em.database.Model(entityContainer)
}
