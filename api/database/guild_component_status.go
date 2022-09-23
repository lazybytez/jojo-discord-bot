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
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"time"
)

const GuildComponentStatusEnabledDisplay = ":white_check_mark:"
const GuildComponentStatusDisabledDisplay = ":x:"

// GuildComponentStatusEntityManager is the GuildCom specific entity manager
// that allows easy access to guilds in the database.
type GuildComponentStatusEntityManager struct {
	em    *GormEntityManager
	cache *cache.Cache[string, GuildComponentStatus]
}

// GuildComponentStatusDBAccess is the interface that defines the capabilities
// of the GuildComponentStatusEntityManager
type GuildComponentStatusDBAccess interface {
	// Get tries to get a GuildComponentStatus from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no GuildComponentStatus can be found, the function returns a new empty
	// GuildComponentStatus.
	Get(guildId uint, componentId uint) (*Guild, error)
	// GetDisplay returns the status of a component in a form
	// that can be directly displayed in Discord.
	GetDisplay(guildId uint, componentId uint) (string, error)
	// Update adds or updates a cached item in the GuildComponentStatus cache.
	Update(guildId string, guild *Guild) error

	// getComponentStatusCacheKey concatenates the passed guild and component ids to create
	// a new unique cache key for the component status
	getComponentStatusCacheKey(guildId uint, componentId uint) string
}

// GuildComponentStatus returns the GuildComponentStatusEntityManager that is currently active,
// which can be used to do GuildComponentStatus specific database actions.
func (em *GormEntityManager) GuildComponentStatus() *GuildComponentStatusEntityManager {
	if nil == em.entityManagers.guildComponentStatusEntityManager {
		gem := &GuildComponentStatusEntityManager{
			em,
			cache.New[string, GuildComponentStatus](10 * time.Minute),
		}
		em.entityManagers.guildComponentStatusEntityManager = gem

		err := gem.cache.EnableAutoCleanup(10 * time.Minute)
		if nil != err {
			em.logger.Err(err, "Failed to initialize periodic cache cleanup task "+
				"for GuildComponentStatus entity manager!")
		}
	}

	return em.entityManagers.guildComponentStatusEntityManager
}

// Get tries to get a GuildComponentStatus from the
// cache. If no cache entry is present, a request to the database will be made.
// If no GuildComponentStatus can be found, the function returns a new empty
// GuildComponentStatus.
func (gcsem *GuildComponentStatusEntityManager) Get(guildId uint, componentId uint) (*GuildComponentStatus, error) {
	cacheKey := gcsem.getComponentStatusCacheKey(guildId, componentId)
	comp, ok := cache.Get(gcsem.cache, cacheKey)

	if ok {
		return comp, nil
	}

	regComp := &GuildComponentStatus{}
	queryStr := ColumnGuild + " = ? AND " + ColumnComponent + " = ?"
	err := gcsem.em.GetFirstEntity(regComp, queryStr, guildId, componentId)
	if nil != err {
		return regComp, err
	}

	gcsem.Update(guildId, componentId, regComp)

	return regComp, nil
}

// GetDisplay returns the status of a component in a form
// that can be directly displayed in Discord.
func (gcsem *GuildComponentStatusEntityManager) GetDisplay(guildId uint, componentId uint) (string, error) {
	compState, err := gcsem.Get(guildId, componentId)
	if nil != err {
		return GuildComponentStatusDisabledDisplay, err
	}

	if compState.Enabled {
		return GuildComponentStatusEnabledDisplay, nil
	}

	return GuildComponentStatusDisabledDisplay, nil
}

// Update adds or updates a cached item in the GuildComponentStatus cache.
func (gcsem *GuildComponentStatusEntityManager) Update(guildId uint, componentId uint, component *GuildComponentStatus) {
	cache.Update(gcsem.cache, gcsem.getComponentStatusCacheKey(guildId, componentId), component)
}

// getComponentStatusCacheKey concatenates the passed guild and component ids to create
// a new unique cache key for the component status
func (gcsem *GuildComponentStatusEntityManager) getComponentStatusCacheKey(guildId uint, componentId uint) string {
	return fmt.Sprintf("%v_%v", guildId, componentId)
}
