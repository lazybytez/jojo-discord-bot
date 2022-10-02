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

// GormDatabaseAccess is a container struct that holds the
// gorm.DB instance to use for entities interaction.
type GormDatabaseAccess struct {
	database *gorm.DB
}

// New creates a new GormDatabaseAccess instance for the entities API
// and register default entity types that are managed by
// the bot core.
func New(db *gorm.DB) *GormDatabaseAccess {
	return &GormDatabaseAccess{
		db,
	}
}

// RegisterEntity registers a new entity (struct) and runs its automated
// migration to ensure the entities schema is up-to-date.
func (gda *GormDatabaseAccess) RegisterEntity(entityType interface{}) error {
	err := gda.database.AutoMigrate(entityType)
	if nil != err {
		return err
	}

	return nil
}

// GetFirstEntity fills the passed entity container with the first
// found entity matching the passed conditions.
//
// Returns an error if no record could be found.
func (gda *GormDatabaseAccess) GetFirstEntity(entityContainer interface{}, conditions ...interface{}) error {
	return gda.database.First(entityContainer, conditions...).Error
}

// GetLastEntity fills the passed entity container with the last
// found entity matching the passed conditions.
//
// Returns false if no entries could be found.
func (gda *GormDatabaseAccess) GetLastEntity(entityContainer interface{}, conditions ...interface{}) error {
	return gda.database.Last(entityContainer, conditions...).Error
}

// GetEntities fills the passed entities slice with the entities
// that have been found for the specified condition.
func (gda *GormDatabaseAccess) GetEntities(entities interface{}, conditions ...interface{}) error {
	return gda.database.Find(entities, conditions...).Error
}

// Create creates the passed entity in the entities
func (gda *GormDatabaseAccess) Create(entity interface{}) error {
	return gda.database.Create(entity).Error
}

// Save upserts the passed entity in the entities
func (gda *GormDatabaseAccess) Save(entity interface{}) error {
	return gda.database.Save(entity).Error
}

// UpdateEntity can be used to update the passed entity in the entities
func (gda *GormDatabaseAccess) UpdateEntity(entityContainer interface{}, column string, value interface{}) error {
	return gda.database.Model(entityContainer).Update(column, value).Error
}

// DeleteEntity deletes the passed entity from the entities.
func (gda *GormDatabaseAccess) DeleteEntity(entityContainer interface{}) error {
	return gda.database.Delete(entityContainer).Error
}

// WorkOn returns a gorm.DB pointer that allows to do a custom search or actions on entities.
//
// The returned gorm.DB instance is created by using gorm.DB.Model() and is therefore
// already prepared to get started with applying filters.
// This function is the only interface point to get direct access to gorm.DB
func (gda *GormDatabaseAccess) WorkOn(entityContainer interface{}) *gorm.DB {
	return gda.database.Model(entityContainer)
}

// DB returns the gorm.DB instance used by the entities api.
func (gda *GormDatabaseAccess) DB() *gorm.DB {
	return gda.database
}
