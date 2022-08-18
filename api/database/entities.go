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

import "C"
import (
	"gorm.io/gorm"
)

// entityManagers is a struct embedded by EntityManager
// that holds the instances of the entity specific entity managers
type entityManagers struct {
	guild                              *GuildEntityManager
	globalComponentStatusEntityManager *GlobalComponentStatusEntityManager
	registeredComponentEntityManager   *RegisteredComponentEntityManager
	guildComponentStatusEntityManager  *GuildComponentStatusEntityManager
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
	// GuildComponentStatus returns the GuildComponentStatusEntityManager that is currently active,
	// which can be used to do GuildComponentStatus specific database actions.
	GuildComponentStatus() *GuildComponentStatusEntityManager
}

// RegisteredComponent represents a single component that is or was known
// to the system.
//
// Single purpose of this struct is to provide a database
// table with which relations can be build to ensure integrity
// of the GuildComponentStatus and GlobalComponentStatus tables.
type RegisteredComponent struct {
	gorm.Model
	Code           string `gorm:"uniqueIndex"`
	Name           string
	Description    string
	DefaultEnabled bool
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
	ComponentID uint                `gorm:"index:idx_global_component_status_component_id;"`
	Component   RegisteredComponent `gorm:"constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// GuildComponentStatus holds the status of a component on a specific server
type GuildComponentStatus struct {
	gorm.Model
	GuildID     uint                `gorm:"index:idx_guild_component_status_guild_id;index:idx_guild_component_status_guild_id_component_id;"`
	Guild       Guild               `gorm:"constraint:OnDelete:CASCADE;"`
	ComponentID uint                `gorm:"index:idx_guild_component_status_component_id;index:idx_guild_component_status_guild_id_component_id;"`
	Component   RegisteredComponent `gorm:"constraint:OnDelete:CASCADE;"`
	Enabled     bool
}

// SlashCommand represents an available slash-command in the database
type SlashCommand struct {
	gorm.Model
	RegisteredComponentID uint                `gorm:"index:idx_slash_command_component_id;"`
	RegisteredComponent   RegisteredComponent `gorm:"constraint:OnDelete:CASCADE;"`
	Name                  string              `gorm:"uniqueIndex"`
}

// registerDefaultEntities takes care of letting gorm
// know about all entities in this file.
func registerDefaultEntities(em *EntityManager) error {
	// Guild related entities
	err := em.RegisterEntity(&Guild{})
	if nil != err {
		return err
	}

	// Component related entities
	err = em.RegisterEntity(&RegisteredComponent{})
	if nil != err {
		return err
	}
	err = em.RegisterEntity(&GlobalComponentStatus{})
	if nil != err {
		return err
	}
	err = em.RegisterEntity(&GuildComponentStatus{})
	if nil != err {
		return err
	}

	// Slash-command related entities
	err = em.RegisterEntity(&SlashCommand{})
	if nil != err {
		return err
	}
	err = em.RegisterEntity(&ActiveSlashCommand{})
	if nil != err {
		return err
	}

	return nil
}
