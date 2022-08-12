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
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestExtractTypeName(t *testing.T) {
	type TestStruct struct {
		Name string
	}

	tables := []struct {
		in  interface{}
		out string
	}{
		{discordgo.Session{}, "Session"},
		{&discordgo.Session{}, "Session"},
		{TestStruct{}, "TestStruct"},
		{&TestStruct{}, "TestStruct"},
	}

	for _, table := range tables {
		if r := ExtractTypeName(table.in); r != table.out {
			t.Errorf("output of to type name extraction with \"%v\" was incorrect, got: %v, want: %v.",
				table.in,
				r,
				table.out)
		}
	}
}
