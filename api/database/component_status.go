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

import "gorm.io/gorm"

const ColumnGuild = "guild"
const ColumnComponent = "component"
const ColumnEnabled = "enabled"

// ComponentStatus holds the status of a component on a specific server
type ComponentStatus struct {
	gorm.Model
	Guild     Guild  `gorm:"foreignKey:ID;index:idx_guild_component"`
	Component string `gorm:"foreignKey:ID;index:idx_guild_component;index:idx_component"`
	Enabled   bool
}

// GlobalComponentStatus holds the status of a component in the global context.
// This allows disabling a bugging component globally if necessary.
type GlobalComponentStatus struct {
	gorm.Model
	Component string `gorm:"foreignKey:ID;index:idx_guild_component;index:idx_component"`
	Enabled   bool
}
