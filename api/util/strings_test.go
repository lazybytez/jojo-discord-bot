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

import "testing"

func TestStringToSnakeCase(t *testing.T) {
	tables := []struct {
		in  string
		out string
	}{
		{"A simple test", "a_simple_test"},
		{"Another More Complex Test", "another_more_complex_test"},
		{"sOmE wEiRd TeST", "s_om_e_w_ei_rd_te_st"},
		{"Also with_underscores", "also_with_underscores"},
		{"5 test in a ROW", "5_test_in_a_row"},
	}

	for _, table := range tables {
		if r := StringToSnakeCase(table.in); r != table.out {
			t.Errorf("output of to snake case with \"%v\" was incorrect, got: %v, want: %v.",
				table.in,
				r,
				table.out)
		}
	}
}
