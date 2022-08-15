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
)

// SlashCommand represents an available slash-command in the database
type SlashCommand struct {
	gorm.Model
	RegisteredComponentID uint
	RegisteredComponent   database.RegisteredComponent
	Name                  string `gorm:"uniqueIndex"`
}

// slashCommandCache is the cache used to reduce
// amount of database calls for the global component status.
var slashCommandCache = cache.New[string, SlashCommand](0)

// GetSlashCommand tries to get a SlashCommand from the
// cache. If no cache entry is present, a request to the database will be made.
// If no SlashCommand can be found, the function returns a new empty
// SlashCommand.
func GetSlashCommand(c *Component, slashCommand string) (*SlashCommand, bool) {
	comp, ok := cache.Get(slashCommandCache, slashCommand)

	if ok {
		return comp, true
	}

	regComp := &SlashCommand{}
	ok = database.GetFirstEntity(regComp, database.ColumnName+" = ?", slashCommand)

	UpdateSlashCommand(c, slashCommand, regComp)

	return regComp, ok
}

// UpdateSlashCommand adds or updates a cached item in the SlashCommand cache.
func UpdateSlashCommand(_ *Component, slashCommand string, component *SlashCommand) {
	cache.Update(slashCommandCache, slashCommand, component)
}
