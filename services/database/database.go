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
	"gorm.io/gorm"
)

// GormDatabaseAccessor is a container struct that holds the
// gorm.DB instance to use for database interaction.
type GormDatabaseAccessor struct {
	database *gorm.DB
}

// New creates a new GormDatabaseAccessor instance for the database API
// and register default entity types that are managed by
// the bot core.
func New(db *gorm.DB) *GormDatabaseAccessor {
	return &GormDatabaseAccessor{
		db,
	}
}

// RegisterEntity registers a new entity (struct) and runs its automated
// migration to ensure the database schema is up-to-date.
func (em *GormDatabaseAccessor) RegisterEntity(entityType interface{}) error {
	err := em.database.AutoMigrate(entityType)
	if nil != err {
		return err
	}

	return nil
}

// GetFirstEntity fills the passed entity container with the first
// found entity matching the passed conditions.
//
// Returns an error if no record could be found.
func (em *GormDatabaseAccessor) GetFirstEntity(entityContainer interface{}, conditions ...interface{}) error {
	return em.database.First(entityContainer, conditions...).Error
}

// GetLastEntity fills the passed entity container with the last
// found entity matching the passed conditions.
//
// Returns false if no entries could be found.
func (em *GormDatabaseAccessor) GetLastEntity(entityContainer interface{}, conditions ...interface{}) error {
	return em.database.Last(entityContainer, conditions...).Error
}

// GetEntities fills the passed entities slice with the entities
// that have been found for the specified condition.
func (em *GormDatabaseAccessor) GetEntities(entities interface{}, conditions ...interface{}) error {
	return em.database.Find(entities, conditions...).Error
}

// Create creates the passed entity in the database
func (em *GormDatabaseAccessor) Create(entity interface{}) error {
	return em.database.Create(entity).Error
}

// Save upserts the passed entity in the database
func (em *GormDatabaseAccessor) Save(entity interface{}) error {
	return em.database.Save(entity).Error
}

// UpdateEntity can be used to update the passed entity in the database
func (em *GormDatabaseAccessor) UpdateEntity(entityContainer interface{}, column string, value interface{}) error {
	return em.database.Model(entityContainer).Update(column, value).Error
}

// DeleteEntity deletes the passed entity from the database.
func (em *GormDatabaseAccessor) DeleteEntity(entityContainer interface{}) error {
	return em.database.Delete(entityContainer).Error
}

// WorkOn returns a gorm.DB pointer that allows to do a custom search or actions on entities.
//
// The returned gorm.DB instance is created by using gorm.DB.Model() and is therefore
// already prepared to get started with applying filters.
// This function is the only interface point to get direct access to gorm.DB
func (em *GormDatabaseAccessor) WorkOn(entityContainer interface{}) *gorm.DB {
	return em.database.Model(entityContainer)
}
