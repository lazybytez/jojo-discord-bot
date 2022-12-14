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

package statistics

import (
	"os"
	"time"
)

var (
	lastGuildCountUpdate    time.Time
	cachedGuildCount        int
	cachedSlashCommandCount = -1
	cachedClusterId         string
)

// collectGuildCount returns the currently cached guild count
// or recomputes the guild count.
//
// The guild count is cached for 10 minutes.
func collectGuildCount() int {
	if time.Since(lastGuildCountUpdate) > 10*time.Minute {
		cachedGuildCount = C.DiscordApi().GuildCount()

		lastGuildCountUpdate = time.Now()
	}

	return cachedGuildCount
}

// collectSlashCommandCount returns the current count of registered slash commands.
// As slash-commands should be only registered during startup, we just cache
// the resulting value infinite and return it on later calls.
func collectSlashCommandCount() int {
	if -1 == cachedSlashCommandCount {
		cachedSlashCommandCount = C.SlashCommandManager().GetCommandCount()
	}

	return cachedSlashCommandCount
}

// collectClusterId returns the cluster id of the current instance.
// The cluster id is equal to the hostname of the current system.
func collectClusterId() string {
	if "" == cachedClusterId {
		var err error
		cachedClusterId, err = os.Hostname()
		if nil != err {
			cachedClusterId = "Error"
		}
	}

	return cachedClusterId
}
