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
	"fmt"
	"github.com/lazybytez/jojo-discord-bot/api"
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

// LoggerPrefix is the prefix used for logging
// when no component is associated with an action.
const LoggerPrefix = "DB API"

var database *gorm.DB

// Init database instance for the database API
// and register default entity types that are managed by
// the bot core.
func Init(db *gorm.DB) error {
	if nil != database {
		return fmt.Errorf("database API already initialized")
	}

	database = db

	return autoMigrateDefaultEntities()
}

// autoMigrateDefaultEntities runs the automated migrations
// for all core entities which are not managed through the
// components of the application but rather through the core
func autoMigrateDefaultEntities() error {
	return RegisterEntity(nil, &Guild{})
}

// RegisterEntity registers a new entity (struct) and runs its automated
// migration to ensure the database schema is up-to-date.
func RegisterEntity[C any](c *api.Component, entityType *C) error {
	componentName := LoggerPrefix
	if nil != c {
		componentName = c.Name
	}

	err := database.AutoMigrate(entityType)
	if nil != err {
		return err
	}

	log.Info(componentName, "Auto-migrated entity \"%v\"", util.ExtractTypeName(entityType))

	return nil
}

// Create creates the passed entity in the database
func Create[C any](entity *C) {
	database.Create(entity)
}

// GetFirstEntity fills the passed entity container with the first
// found entity matching the passed conditions.
//
// Returns false if no entries could be found.
func GetFirstEntity[C any](c *api.Component, entityContainer *C, conditions ...interface{}) bool {
	db := database.First(entityContainer, conditions...)

	if nil != db.Error {
		log.Err(c.Name, db.Error, "Something went wrong when retrieving first entity!")
	}

	return db.RowsAffected > 0
}

// GetLastEntity fills the passed entity container with the last
// found entity matching the passed conditions.
//
// Returns false if no entries could be found.
func GetLastEntity[C any](c *api.Component, entityContainer *C, conditions ...interface{}) bool {
	db := database.Last(entityContainer, conditions...)

	if nil != db.Error {
		log.Err(c.Name, db.Error, "Something went wrong when retrieving last entity!")
	}

	return db.RowsAffected > 0
}

// GetEntities fills the passed entity container slice with the entities
// that have been found for the specified condition.
func GetEntities[C any](c *api.Component, entityContainer []*C, conditions ...interface{}) bool {
	db := database.Find(entityContainer, conditions...)

	if nil != db.Error {
		log.Err(c.Name, db.Error, "Something went wrong when retrieving entities!")
	}

	return db.RowsAffected > 0
}

// GetEntitiesComplex returns a gorm.DB pointer that allows to do a custom search.
//
// The returned gorm.DB instance is created by using gorm.DB.Model and is therefore
// already prepared to get started with applying filters.
func GetEntitiesComplex[C any](entityContainer []*C) *gorm.DB {
	return database.Model(entityContainer)
}

// UpdateEntity can be used to update the passed entity in the database
func UpdateEntity[C any](c *api.Component, entityContainer *C, column string, value interface{}) bool {
	db := database.Model(entityContainer).Update(column, value)
	if nil != db.Error {
		log.Err(c.Name, db.Error, "Something went wrong when retrieving an entity!")
	}

	return nil != db.Error && db.RowsAffected > 0
}

// DeleteEntity deletes the passed entity from the database.
func DeleteEntity[C any](c *api.Component, entityContainer *C) bool {
	db := database.Delete(entityContainer)

	if nil != db.Error {
		log.Err(c.Name, db.Error, "Something went wrong when deleting an entity!")
	}

	return nil != db.Error && db.RowsAffected > 0
}
