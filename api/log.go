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

import "github.com/lazybytez/jojo-discord-bot/api/log"

// Logger is used to obtain the Logging of a component
//
// On first call, this function initializes the private Component.logger
// field. On consecutive calls, the already present Logging will be used.
func (c *Component) Logger() *log.Logger {
	if nil == c.logger {
		c.logger = log.New(c.Name, nil)
	}

	return c.logger
}
