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

package main

import (
	_ "github.com/lazybytez/jojo-discord-bot/components"
	_ "github.com/lazybytez/jojo-discord-bot/core_components"

	"github.com/lazybytez/jojo-discord-bot/internal"
)

// Entrypoint of Go
// Call real internal.Bootstrap function of internal package
//
// OpenAPI / Swagger data
// @contact.name   Lazy Bytez
// @contact.url https://lazybytez.de/
// @contact.email  contact@lazybytez.de

// @license.name GNU Affero General Public License v3.0
// @license.url   https://www.gnu.org/licenses/agpl-3.0.html
func main() {
	internal.Bootstrap()
}
