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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInArrayOrSlice(t *testing.T) {
	tables := []struct {
		arrayOrSlice []string
		needle       string
		result       bool
	}{
		{[]string{}, "test", false},
		{[]string{"test"}, "test", true},
		{[]string{"test", "other"}, "test", true},
		{[]string{"other", "test"}, "test", true},
		{[]string{"other", "with", "many", "test", "elements"}, "test", true},
		{[]string{"other", "with", "many", "elements"}, "test", false},
	}

	for _, table := range tables {
		result := InArrayOrSlice(table.arrayOrSlice, table.needle)

		assert.Equalf(t,
			table.result,
			result,
			"Failed to assert InArrayOrSlice with array \"%v\" and needle \"%s\"",
			table.arrayOrSlice,
			table.needle)
	}
}

func TestUniqueArrayOrSlice(t *testing.T) {
	tables := []struct {
		arrayOrSlice []string
		result       []string
	}{
		{[]string{}, []string{}},
		{[]string{"test"}, []string{"test"}},
		{[]string{"test", "test"}, []string{"test"}},
		{[]string{"test", "other"}, []string{"test", "other"}},
		{[]string{"test", "other", "test"}, []string{"test", "other"}},
		{[]string{"test", "other", "test", "other"}, []string{"test", "other"}},
		{[]string{"other", "with", "many", "test", "elements"}, []string{"other", "with", "many", "test", "elements"}},
		{[]string{"other", "with", "test", "many", "test", "elements"}, []string{"other", "with", "test", "many", "elements"}},
	}

	for _, table := range tables {
		result := UniqueArrayOrSlice(table.arrayOrSlice)

		assert.Equalf(t,
			table.result,
			result,
			"Failed to assert UniqueArrayOrSlice with array \"%v\"",
			table.arrayOrSlice)
	}
}
