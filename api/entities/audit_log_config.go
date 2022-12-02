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
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"gorm.io/gorm"
	"strconv"
)

// AuditLogConfig holds the guild specific configuration for audit logging.
type AuditLogConfig struct {
	gorm.Model
	GuildID   uint  `gorm:"uniqueIndex;"`
	Guild     Guild `gorm:"constraint:OnDelete:CASCADE;"`
	ChannelId *uint64
	Enabled   bool
}

// AuditLogConfigEntityManager is the audit log config specific entity manager
// that allows easy access to audit log configurations.
type AuditLogConfigEntityManager struct {
	EntityManager
}

// NewAuditLogConfigEntityManager creates a new AuditLogConfigEntityManager.
func NewAuditLogConfigEntityManager(entityManager EntityManager) *AuditLogConfigEntityManager {
	alcem := &AuditLogConfigEntityManager{entityManager}

	return alcem
}

// GetByGuildId tries to get a AuditLogConfig by its guild ID.
// The function uses a cache and first tries to resolve a value from it.
// If no cache entry is present, a request to the entities will be made.
// If no AuditLogConfig can be found, the function returns a new empty
// AuditLogConfig.
func (alcem *AuditLogConfigEntityManager) GetByGuildId(guildId uint) (*AuditLogConfig, error) {
	cacheKey := alcem.getCacheKey(guildId)
	auditLogConfig := cache.Get(cacheKey, &AuditLogConfig{})

	if nil == auditLogConfig {
		return auditLogConfig, nil
	}

	auditLogConfig = &AuditLogConfig{}
	queryStr := ColumnGuild + " = ?"
	err := alcem.DB().GetFirstEntity(auditLogConfig, queryStr, guildId)
	if nil != err {
		return auditLogConfig, err
	}

	cache.Update(cacheKey, AuditLogConfig{})

	return auditLogConfig, nil
}

// Create saves the passed AuditLogConfig in the database.
// Use Update or Save to update an already existing AuditLogConfig.
func (alcem *AuditLogConfigEntityManager) Create(auditLogConfig *AuditLogConfig) error {
	err := alcem.DB().Create(auditLogConfig)
	if nil != err {
		return err
	}

	// Invalidate cache item (if present)
	cache.Invalidate(alcem.getCacheKey(auditLogConfig.GuildID), &AuditLogConfig{})

	return nil
}

// Save updates the passed Guild in the database.
// This does a generic update, use Update to do a precise and more performant update
// of the entity when only updating a single field!
func (alcem *AuditLogConfigEntityManager) Save(auditLogConfig *AuditLogConfig) error {
	err := alcem.DB().Save(auditLogConfig)
	if nil != err {
		return err
	}

	// Invalidate cache item (if present)
	cache.Invalidate(alcem.getCacheKey(auditLogConfig.GuildID), &AuditLogConfig{})

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (alcem *AuditLogConfigEntityManager) Update(auditLogConfig *AuditLogConfig, column string, value interface{}) error {
	err := alcem.DB().UpdateEntity(auditLogConfig, column, value)
	if nil != err {
		return err
	}

	// Invalidate cache item (if present)
	cache.Invalidate(alcem.getCacheKey(auditLogConfig.GuildID), &AuditLogConfig{})
	
	return nil
}

// getCacheKey returns the computed cache key used to cache
// AuditLogConfig objects.
func (alcem *AuditLogConfigEntityManager) getCacheKey(guildId uint) string {
	return strconv.FormatUint(uint64(guildId), 10)
}
