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
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"strconv"
	"time"
)

// GuildEntityManager is the Guild specific entity manager
// that allows easy access to guilds in the database.
type GuildEntityManager struct {
	em    *EntityManager
	cache *cache.Cache[uint64, Guild]
}

// GuildDBAccess is the interface that defines the capabilities
// of the GuildEntityManager
type GuildDBAccess interface {
	// Get tries to get a Guild from the
	// cache. If no cache entry is present, a request to the database will be made.
	// If no Guild can be found, the function returns a new empty
	// Guild.
	Get(guildId string) (*Guild, error)
	// Update adds or updates a cached item in the Guild cache.
	Update(guildId string, guild *Guild) error
}

// Guilds returns the GuildEntityManager that is currently active,
// which can be used to do Guild specific database actions.
func (em *EntityManager) Guilds() *GuildEntityManager {
	if nil != em.entityManagers.guild {
		gem := &GuildEntityManager{
			em,
			cache.New[uint64, Guild](10 * time.Minute),
		}
		em.entityManagers.guild = gem

		gem.cache.EnableAutoCleanup(10 * time.Minute)
	}

	return em.entityManagers.guild
}

// Get tries to get a Guild from the
// cache. If no cache entry is present, a request to the database will be made.
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
	err = gem.em.GetFirstEntity(guild, ColumnGuildId+" = ?", guildIdInt)
	if nil != err {
		return &Guild{}, err
	}

	err = gem.Update(guildId, guild)

	return guild, err
}

// Update adds or updates a cached item in the Guild cache.
func (gem *GuildEntityManager) Update(guildId string, guild *Guild) error {
	guildIdInt, err := strconv.ParseUint(guildId, 10, 64)
	if nil != err {
		return err
	}

	cache.Update(gem.cache, guildIdInt, guild)
	return nil
}
