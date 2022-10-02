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
	"github.com/lazybytez/jojo-discord-bot/api/util"
	"github.com/lazybytez/jojo-discord-bot/services"
)

// defaultEntities holds a list of all entities that should be
// registered by default. The list is registered in the order
// the entities are added to the list.
var defaultEntities = []interface{}{
	&Guild{},
	&RegisteredComponent{},
	&GlobalComponentStatus{},
	&GuildComponentStatus{},
}

// EntityManager is a struct embedded by GormDatabaseAccessor
// that holds the instances of the entity specific entity managers
type EntityManager struct {
	database services.DatabaseAccess
	logger   services.Logger

	guild                              *GuildEntityManager
	globalComponentStatusEntityManager *GlobalComponentStatusEntityManager
	registeredComponentEntityManager   *RegisteredComponentEntityManager
	guildComponentStatusEntityManager  *GuildComponentStatusEntityManager
}

// NewEntityManager creates a new EntityManager using the given services.Logger and
// services.DatabaseAccess instances.
func NewEntityManager(database services.DatabaseAccess, logger services.Logger) EntityManager {
	return EntityManager{
		database: database,
		logger:   logger,
	}
}

// RegisterEntity registers a new entity (struct) and runs its automated
// migration to ensure the entities schema is up-to-date.
func (em *EntityManager) RegisterEntity(entityType interface{}) error {
	err := em.database.RegisterEntity(entityType)
	if nil != err {
		em.logger.Err(err, "Failed to auto-migrated entity \"%v\"", util.ExtractTypeName(entityType))

		return err
	}

	em.logger.Info("Auto-migrated entity \"%v\"", util.ExtractTypeName(entityType))

	return nil
}

// registerDefaultEntities takes care of letting gorm
// know about all entities in this file.
func registerDefaultEntities(em EntityManager) error {
	for _, entity := range defaultEntities {
		err := em.RegisterEntity(entity)
		if nil != err {
			return err
		}
	}

	return nil
}
