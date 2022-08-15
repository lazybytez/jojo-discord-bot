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

// entityManagers is a struct embedded by EntityManager
// that holds the instances of the entity specific entity managers
type entityManagers struct {
	guild                              *GuildEntityManager
	globalComponentStatusEntityManager *GlobalComponentStatusEntityManager
	registeredComponentEntityManager   *RegisteredComponentEntityManager
}

// entitySpecificManagerAccess contains methods that allow to retrieve
// entity specific entity managers that provide caching and dedicated functions
// to work with specific entities.
type entitySpecificManagerAccess interface {
	// Guilds returns the GuildEntityManager that is currently active,
	// which can be used to do Guild specific database actions.
	Guilds() *GuildEntityManager
	// GlobalComponentStatus returns the GlobalComponentStatusEntityManager that is currently active,
	// which can be used to do GlobalComponentStatus specific database actions.
	GlobalComponentStatus() *GlobalComponentStatusEntityManager
	// RegisteredComponent returns the RegisteredComponentEntityManager that is currently active,
	// which can be used to do RegisteredComponent specific database actions.
	RegisteredComponent() *RegisteredComponentEntityManager
}

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

// GlobalComponentStatus holds the status of a component in the global context.
// This allows disabling a bugging component globally if necessary.
type GlobalComponentStatus struct {
	gorm.Model
	ComponentID uint
	Component   RegisteredComponent `gorm:"index:idx_guild_component;index:idx_component;constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// RegisteredComponent represents a single component that is or was known
// to the system.
//
// Single purpose of this struct is to provide a database
// table with which relations can be build to ensure integrity
// of the ComponentStatus and GlobalComponentStatus tables.
type RegisteredComponent struct {
	gorm.Model
	Code string `gorm:"uniqueIndex"`
}
