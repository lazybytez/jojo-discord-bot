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
	"github.com/lazybytez/jojo-discord-bot/api"
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
	GuildID int `gorm:"index"`
	Name    string
}

// guildCache caches the guilds the current instance is on.
// A Guild is cached for 10 minutes before the application needs to pull it
// again.
var guildCache = cache.New[int, Guild](10 * time.Minute)

// init ensure cache cleanup task is running
func init() {
	guildCache.EnableAutoCleanup(10 * time.Minute)
}

// GetGuild tries to get a Guild from the
// cache. If no cache entry is present, a request to the database will be made.
// If no Guild can be found, the function returns a new empty
// Guild.
func GetGuild(c *api.Component, guildId string) (*Guild, bool) {
	guildIdInt, err := strconv.Atoi(guildId)
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
	ok = GetFirstEntity(c, regComp, ColumnGuildId+" = ?", guildIdInt)

	UpdateGuild(c, guildId, regComp)

	return regComp, ok
}

// UpdateGuild adds or updates a cached item in the Guild cache.
func UpdateGuild(c *api.Component, guildId string, component *Guild) {
	guildIdInt, err := strconv.Atoi(guildId)
	if nil != err {
		c.Logger().Err(
			err,
			"Failed to get guild from database for id string \"%v\" due to integer conversion issues!",
			guildId)

		return
	}

	cache.Update(guildCache, guildIdInt, component)
}
