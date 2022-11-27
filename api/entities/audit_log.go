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
	"gorm.io/gorm"
)

// AuditLog holds the audit log of the bot.
// The audit log contains information about different administrative actions.
// The most important purpose of the audit log is to ensure that any change done to the bots configuration
// done by the users of the bot can be tracked. This can be actions on guilds or actions performed in private
// messages.
type AuditLog struct {
	gorm.Model
	GuildID               uint                `gorm:"index:idx_audit_log_guild_id;index:idx_audit_log_guild_id_user_id;index:idx_audit_log_guild_id_component_id_user_id;"`
	Guild                 Guild               `gorm:"constraint:OnDelete:CASCADE;"`
	RegisteredComponentID uint                `gorm:"index:idx_audit_log_guild_id;index:idx_audit_log_guild_id_user_id;index:idx_audit_log_guild_id_component_id_user_id;"`
	RegisteredComponent   RegisteredComponent `gorm:"constraint:OnDelete:CASCADE;"`
	UserID                uint64              `gorm:"index:idx_audit_log_user_id;index:idx_audit_log_guild_id_user_id;index:idx_audit_log_guild_id_component_id_user_id;"`
	Message               string
}

// AuditLogEntityManager is the audit log specific entity manager
// that allows easy access to guilds in the entities.
type AuditLogEntityManager struct {
	EntityManager
}

// NewAuditLogEntityManager creates a new AuditLogEntityManager.
func NewAuditLogEntityManager(entityManager EntityManager) *AuditLogEntityManager {
	alem := &AuditLogEntityManager{
		entityManager,
	}

	return alem
}

// Create saves the passed AuditLog in the database.
// Use Update or Save to update an already existing AuditLog.
func (alem *AuditLogEntityManager) Create(guild *AuditLog) error {
	err := alem.DB().Create(guild)
	if nil != err {
		return err
	}

	return nil
}

// Save updates the passed AuditLog in the database.
// This does a generic update, use Update to do a precise and more performant update
// of the entity when only updating a single field!
func (alem *AuditLogEntityManager) Save(guild *AuditLog) error {
	err := alem.DB().Save(guild)
	if nil != err {
		return err
	}

	return nil
}

// Update updates the defined field on the entity and saves it in the database.
func (alem *AuditLogEntityManager) Update(auditLog *AuditLog, column string, value interface{}) error {
	err := alem.DB().UpdateEntity(auditLog, column, value)
	if nil != err {
		return err
	}

	return nil
}
