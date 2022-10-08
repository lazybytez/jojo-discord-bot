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
	"fmt"
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"gorm.io/gorm"
	"time"
)

const GuildComponentStatusEnabledDisplay = ":white_check_mark:"
const GuildComponentStatusDisabledDisplay = ":x:"

// GuildComponentStatus holds the status of a component on a specific server
type GuildComponentStatus struct {
	gorm.Model
	GuildID     uint                `gorm:"index:idx_guild_component_status_guild_id;index:idx_guild_component_status_guild_id_component_id;"`
	Guild       Guild               `gorm:"constraint:OnDelete:CASCADE;"`
	ComponentID uint                `gorm:"index:idx_guild_component_status_component_id;index:idx_guild_component_status_guild_id_component_id;"`
	Component   RegisteredComponent `gorm:"constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// GuildComponentStatusEntityManager is the GuildCom specific entity manager
// that allows easy access to guilds in the entities.
type GuildComponentStatusEntityManager struct {
	EntityManager

	cache *cache.Cache[string, GuildComponentStatus]
}

// NewGuildComponentStatusEntityManager creates a new GuildComponentStatusEntityManager.
func NewGuildComponentStatusEntityManager(entityManager EntityManager) *GuildComponentStatusEntityManager {
	gem := &GuildComponentStatusEntityManager{
		entityManager,
		cache.New[string, GuildComponentStatus](10 * time.Minute),
	}

	err := gem.cache.EnableAutoCleanup(10 * time.Minute)
	if nil != err {
		entityManager.Logger().Err(err, "Failed to initialize periodic cache cleanup task "+
			"for GuildComponentStatus entity manager!")
	}

	return gem
}

// Get tries to get a GuildComponentStatus from the
// cache. If no cache entry is present, a request to the entities will be made.
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
	err := gcsem.DB().GetFirstEntity(regComp, queryStr, guildId, componentId)
	if nil != err {
		return regComp, err
	}

	cache.Update(gcsem.cache, gcsem.getComponentStatusCacheKey(regComp.GuildID, regComp.ComponentID), regComp)

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

// Create saves the passed Guild in the database.
// Use Update or Save to update an already existing Guild.
func (gcsem *GuildComponentStatusEntityManager) Create(guildComponentStatus *GuildComponentStatus) error {
	err := gcsem.DB().Create(guildComponentStatus)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(
		gcsem.cache,
		gcsem.getComponentStatusCacheKey(guildComponentStatus.GuildID, guildComponentStatus.ComponentID),
		guildComponentStatus)

	return nil
}

// Save updates the passed Guild in the database.
// This does a generic update, use Update to do a precise and more performant update
// of the entity when only updating a single field!
func (gcsem *GuildComponentStatusEntityManager) Save(guildComponentStatus *GuildComponentStatus) error {
	err := gcsem.DB().Save(guildComponentStatus)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(
		gcsem.cache,
		gcsem.getComponentStatusCacheKey(guildComponentStatus.GuildID, guildComponentStatus.ComponentID),
		guildComponentStatus)

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (gcsem *GuildComponentStatusEntityManager) Update(
	component *GuildComponentStatus,
	column string,
	value interface{},
) error {
	err := gcsem.DB().UpdateEntity(component, column, value)
	if nil != err {
		return err
	}

	cache.Update(gcsem.cache, gcsem.getComponentStatusCacheKey(component.GuildID, component.ComponentID), component)

	return nil
}

// getComponentStatusCacheKey concatenates the passed guild and component ids to create
// a new unique cache key for the component status
func (gcsem *GuildComponentStatusEntityManager) getComponentStatusCacheKey(guildId uint, componentId uint) string {
	return fmt.Sprintf("%v_%v", guildId, componentId)
}
