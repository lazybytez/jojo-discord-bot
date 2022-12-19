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
	"sync"
)

// botStatusManager is the DiscordGoStatusManager instance used across the bots' lifetime.
var botStatusManager *DiscordGoStatusManager

// DiscordGoStatusManager holds the available status of the bot
// and manages the cycling of these.
type DiscordGoStatusManager struct {
	mu      sync.RWMutex
	status  []SimpleBotStatus
	current int
}

// StatusManager manages the available status
// which are set fopr the bot.
type StatusManager interface {
	// AddStatusToRotation adds the given status to the list of
	// rotated status.
	AddStatusToRotation(status SimpleBotStatus)
	// Next works like next on an iterator which self resets automatically.
	Next() *SimpleBotStatus
}

// BotStatusManager returns the current StatusManager which
// allows to add additional status to the bot.
func (c *Component) BotStatusManager() StatusManager {
	return botStatusManager
}

func init() {
	botStatusManager = &DiscordGoStatusManager{
		mu:      sync.RWMutex{},
		status:  make([]SimpleBotStatus, 0),
		current: 0,
	}
}

// AddStatusToRotation adds the given status to the list of
// rotated status.
func (dgsm *DiscordGoStatusManager) AddStatusToRotation(status SimpleBotStatus) {
	dgsm.mu.Lock()
	defer dgsm.mu.Unlock()

	dgsm.status = append(dgsm.status, status)
}

// Next works like next on an iterator which self resets automatically.
func (dgsm *DiscordGoStatusManager) Next() *SimpleBotStatus {
	dgsm.mu.Lock()
	defer dgsm.mu.Unlock()

	statusCount := len(dgsm.status)
	if 0 == statusCount {
		return nil
	}

	if dgsm.current >= statusCount {
		dgsm.current = 0
	}

	status := &dgsm.status[dgsm.current]

	dgsm.current = dgsm.current + 1

	return status
}
