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

package util

import (
	"regexp"
	"strings"
)

// snakeCaseFirstCap matches the first char and can be used to turn it into lowercase
// and add an underscore between a leading special character and the rest of the string
var snakeCaseFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)

// snakeCaseAllCap matches the capitals coming after a capital character or digit.
var snakeCaseAllCap = regexp.MustCompile(`([a-z\\d])([A-Z])`)

// StringToSnakeCase turns the given string into snake-case.
//
// Examples:
//   - ThisIsATest => this_is_a_test
//   - $ThisIsATest => $_this_is_a_test
func StringToSnakeCase(str string) string {
	snake := snakeCaseFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = snakeCaseAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}
