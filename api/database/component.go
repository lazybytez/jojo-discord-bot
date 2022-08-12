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

const ColumnName = "name"

// RegisteredComponent represents a single component that is or was known
// to the system.
//
// Single purpose of this struct is to provide a database
// table with which relations can be build to ensure integrity
// of the ComponentStatus and GlobalComponentStatus tables.
type RegisteredComponent struct {
	gorm.Model
	name string `gorm:"uniqueIndex"`
}
