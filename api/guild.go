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

package api

import (
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"github.com/lazybytez/jojo-discord-bot/api/database"
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

// guildCache caches the guilds the current instance is on.
// A Guild is cached for 10 minutes before the application needs to pull it
// again.
var guildCache = cache.New[uint64, Guild](10 * time.Minute)

// init ensure cache cleanup task is running
func init() {
	guildCache.EnableAutoCleanup(10 * time.Minute)
}

// GetGuild tries to get a Guild from the
// cache. If no cache entry is present, a request to the database will be made.
// If no Guild can be found, the function returns a new empty
// Guild.
func GetGuild(c *Component, guildId string) (*Guild, bool) {
	guildIdInt, err := strconv.ParseUint(guildId, 10, 64)
	if nil != err {
		c.Logger().Err(
			err,
			"Failed to get guild from database for id string \"%v\" due to integer conversion issues!",
			guildId)

		return &Guild{}, false
	}

	comp, ok := cache.Get(guildCache, guildIdInt)

	if ok {
		return comp, true
	}

	regComp := &Guild{}
	ok = database.GetFirstEntity(regComp, ColumnGuildId+" = ?", guildIdInt)

	UpdateGuild(c, guildId, regComp)

	return regComp, ok
}

// UpdateGuild adds or updates a cached item in the Guild cache.
func UpdateGuild(c *Component, guildId string, component *Guild) {
	guildIdInt, err := strconv.ParseUint(guildId, 10, 64)
	if nil != err {
		c.Logger().Err(
			err,
			"Failed to get guild from database for id string \"%v\" due to integer conversion issues!",
			guildId)

		return
	}

	cache.Update(guildCache, guildIdInt, component)
}
