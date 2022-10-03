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
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Guild represents a single Discord guild
// that the bot is currently on.
//
// Note that the guild name is just stored for convenience when
// manually searching the DB for a guild.
type Guild struct {
	gorm.Model
	GuildID uint64 `gorm:"uniqueIndex"`
	Name    string
}

// GuildEntityManager is the Guild specific entity manager
// that allows easy access to guilds in the entities.
type GuildEntityManager struct {
	EntityManager
	cache *cache.Cache[uint64, Guild]
}

// NewGuildEntityManager creates a new GuildEntityManager with
// the given EntityManager.
func NewGuildEntityManager(em EntityManager) *GuildEntityManager {
	gem := &GuildEntityManager{
		em,
		cache.New[uint64, Guild](10 * time.Minute),
	}

	err := gem.cache.EnableAutoCleanup(10 * time.Minute)
	if nil != err {
		em.Logger().Err(err, "Failed to initialize periodic cache cleanup task for Guild entity manager!")
	}

	return gem
}

// Get tries to get a Guild from the
// cache. If no cache entry is present, a request to the entities will be made.
// If no Guild can be found, the function returns a new empty
// Guild.
func (gem *GuildEntityManager) Get(guildId string) (*Guild, error) {
	guildIdInt, err := strconv.ParseUint(guildId, 10, 64)
	if nil != err {
		return &Guild{}, err
	}

	comp, ok := cache.Get(gem.cache, guildIdInt)

	if ok {
		return comp, nil
	}

	guild := &Guild{}
	err = gem.DB().GetFirstEntity(guild, ColumnGuildId+" = ?", guildIdInt)
	if nil != err {
		return &Guild{}, err
	}

	cache.Update(gem.cache, guild.GuildID, guild)

	return guild, err
}

// Count returns the number of all guilds stored in the entities
func (gem *GuildEntityManager) Count() (int64, error) {
	var count int64 = 0
	db := gem.DB().WorkOn([]Guild{}).Count(&count)

	return count, db.Error
}

// Create saves the passed Guild in the database.
// Use Update or Save to update an already existing Guild.
func (gem *GuildEntityManager) Create(guild *Guild) error {
	err := gem.DB().Create(guild)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(gem.cache, guild.GuildID, guild)

	return nil
}

// Save updates the passed Guild in the database.
// This does a generic update, use Update to do a precise and more performant update
// of the entity when only updating a single field!
func (gem *GuildEntityManager) Save(guild *Guild) error {
	err := gem.DB().Save(guild)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(gem.cache, guild.GuildID, guild)

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (gem *GuildEntityManager) Update(guild *Guild, column string, value interface{}) error {
	err := gem.DB().UpdateEntity(guild, column, value)
	if nil != err {
		return err
	}

	cache.Update(gem.cache, guild.GuildID, guild)

	return nil
}
