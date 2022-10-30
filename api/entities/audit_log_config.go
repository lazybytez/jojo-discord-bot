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
	"time"
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

	cache *cache.Cache[uint, AuditLogConfig]
}

// NewAuditLogConfigEntityManager creates a new AuditLogConfigEntityManager.
func NewAuditLogConfigEntityManager(entityManager EntityManager) *AuditLogConfigEntityManager {
	alcem := &AuditLogConfigEntityManager{
		entityManager,
		cache.New[uint, AuditLogConfig](10 * time.Minute),
	}

	err := alcem.cache.EnableAutoCleanup(10 * time.Minute)
	if nil != err {
		entityManager.Logger().Err(err, "Failed to initialize periodic cache cleanup task "+
			"for AuditLogConfig entity manager!")
	}

	return alcem
}

// GetByGuildId tries to get a AuditLogConfig by its guild ID.
// The function uses a cache and first tries to resolve a value from it.
// If no cache entry is present, a request to the entities will be made.
// If no AuditLogConfig can be found, the function returns a new empty
// AuditLogConfig.
func (alcem *AuditLogConfigEntityManager) GetByGuildId(guildId uint) (*AuditLogConfig, error) {
	auditLogConfig, ok := cache.Get(alcem.cache, guildId)

	if ok {
		return auditLogConfig, nil
	}

	auditLogConfig = &AuditLogConfig{}
	queryStr := ColumnGuild + " = ?"
	err := alcem.DB().GetFirstEntity(auditLogConfig, queryStr, guildId)
	if nil != err {
		return auditLogConfig, err
	}

	cache.Update(alcem.cache, auditLogConfig.GuildID, auditLogConfig)

	return auditLogConfig, nil
}

// Create saves the passed AuditLogConfig in the database.
// Use Update or Save to update an already existing AuditLogConfig.
func (alcem *AuditLogConfigEntityManager) Create(auditLogConfig *AuditLogConfig) error {
	err := alcem.DB().Create(auditLogConfig)
	if nil != err {
		return err
	}

	// Ensure entity is in cache when just updated
	cache.Update(alcem.cache, auditLogConfig.GuildID, auditLogConfig)

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

	// Ensure entity is in cache when just updated
	cache.Update(alcem.cache, auditLogConfig.GuildID, auditLogConfig)

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (alcem *AuditLogConfigEntityManager) Update(auditLogConfig *AuditLogConfig, column string, value interface{}) error {
	err := alcem.DB().UpdateEntity(auditLogConfig, column, value)
	if nil != err {
		return err
	}

	cache.Update(alcem.cache, auditLogConfig.GuildID, auditLogConfig)

	return nil
}
