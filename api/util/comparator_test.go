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
	"testing"
)

func TestArraysEqualWithStrings(t *testing.T) {
	tables := []struct {
		inputA         *[]string
		inputB         *[]string
		expectedResult bool
	}{
		{nil, nil, true},
		{nil, &[]string{"item1", "item2"}, false},
		{&[]string{"item1", "item2"}, nil, false},
		{&[]string{"item1", "item2", "item3"}, &[]string{"item1", "item2"}, false},
		{&[]string{"item1", "item2"}, &[]string{"item1", "item2", "item3"}, false},
		{&[]string{"item1"}, &[]string{"item1"}, true},
		{&[]string{"item1", "item2"}, &[]string{"item1", "item2"}, true},
		{&[]string{"item2", "item1"}, &[]string{"item1", "item2"}, false},
		{&[]string{"item1", "item2"}, &[]string{"item2", "item1"}, false},
	}

	for _, table := range tables {
		if r := ArraysEqual(table.inputA, table.inputB); r != table.expectedResult {
			t.Errorf("[STRING ARRAY] output of array equals with \"%v\" and \"%v\" "+
				"was incorrect, got: %v, want: %v.",
				table.inputA,
				table.inputB,
				r,
				table.expectedResult)
		}
	}
}

func TestArraysEqualWithInt(t *testing.T) {
	tables := []struct {
		inputA         *[]int
		inputB         *[]int
		expectedResult bool
	}{
		{nil, nil, true},
		{nil, &[]int{42, 64}, false},
		{&[]int{12, 34}, nil, false},
		{&[]int{1, 2, 3}, &[]int{1, 2}, false},
		{&[]int{3, 2}, &[]int{3, 2, 1}, false},
		{&[]int{42}, &[]int{42}, true},
		{&[]int{32, 64}, &[]int{32, 64}, true},
		{&[]int{2, 1}, &[]int{1, 2}, false},
		{&[]int{1, 2}, &[]int{2, 1}, false},
	}

	for _, table := range tables {
		if r := ArraysEqual(table.inputA, table.inputB); r != table.expectedResult {
			t.Errorf("[INT ARRAY] output of array equals with \"%v\" and \"%v\" "+
				"was incorrect, got: %v, want: %v.",
				table.inputA,
				table.inputB,
				r,
				table.expectedResult)
		}
	}
}

func TestArraysEqualWithBool(t *testing.T) {
	tables := []struct {
		inputA         *[]bool
		inputB         *[]bool
		expectedResult bool
	}{
		{nil, nil, true},
		{nil, &[]bool{true, false}, false},
		{&[]bool{true, false}, nil, false},
		{&[]bool{true, false, true}, &[]bool{true, false}, false},
		{&[]bool{true, false}, &[]bool{true, false, true}, false},
		{&[]bool{true}, &[]bool{true}, true},
		{&[]bool{false}, &[]bool{false}, true},
		{&[]bool{true, false}, &[]bool{true, false}, true},
		{&[]bool{false, true}, &[]bool{true, false}, false},
		{&[]bool{true, false}, &[]bool{false, true}, false},
	}

	for _, table := range tables {
		if r := ArraysEqual(table.inputA, table.inputB); r != table.expectedResult {
			t.Errorf("[BOOL ARRAY] output of array equals with \"%v\" and \"%v\" was "+
				"incorrect, got: %v, want: %v.",
				table.inputA,
				table.inputB,
				r,
				table.expectedResult)
		}
	}
}

func TestMapsEqualEqualWithStringString(t *testing.T) {
	tables := []struct {
		inputA         *map[string]string
		inputB         *map[string]string
		expectedResult bool
	}{
		{nil, nil, true},
		{nil, &map[string]string{"a": "item1", "b": "item2"}, false},
		{&map[string]string{"a": "item1", "b": "item2"}, nil, false},
		{
			&map[string]string{"a": "item1", "b": "item2", "c": "item3"},
			&map[string]string{"a": "item1", "b": "item2"},
			false,
		},
		{
			&map[string]string{"a": "item1", "b": "item2"},
			&map[string]string{"a": "item1", "b": "item2", "c": "item3"},
			false,
		},
		{
			&map[string]string{"a": "item1"},
			&map[string]string{"a": "item1"},
			true,
		},
		{
			&map[string]string{"a": "item1", "b": "item2"},
			&map[string]string{"a": "item1", "b": "item2"},
			true,
		},
		{
			&map[string]string{"a": "item1", "b": "item2"},
			&map[string]string{"a": "item1", "b": "item2"},
			true,
		},
		{
			&map[string]string{"b": "item2", "a": "item1"},
			&map[string]string{"a": "item1", "b": "item2"},
			true,
		},
	}

	for _, table := range tables {
		if r := MapsEqual(table.inputA, table.inputB); r != table.expectedResult {
			t.Errorf("[STRING STRING MAP] output of array equals with \"%v\" and \"%v\" "+
				"was incorrect, got: %v, want: %v.",
				table.inputA,
				table.inputB,
				r,
				table.expectedResult)
		}
	}
}

func TestMapsEqualEqualWithStringInt(t *testing.T) {
	tables := []struct {
		inputA         *map[string]int
		inputB         *map[string]int
		expectedResult bool
	}{
		{nil, nil, true},
		{nil, &map[string]int{"a": 42, "b": 64}, false},
		{&map[string]int{"a": 42, "b": 64}, nil, false},
		{
			&map[string]int{"a": 1, "b": 2, "c": 3},
			&map[string]int{"a": 1, "b": 2},
			false,
		},
		{
			&map[string]int{"a": 1, "b": 2},
			&map[string]int{"a": 1, "b": 2, "c": 3},
			false,
		},
		{
			&map[string]int{"a": 42},
			&map[string]int{"a": 42},
			true,
		},
		{
			&map[string]int{"a": 32, "b": 64},
			&map[string]int{"a": 32, "b": 64},
			true,
		},
		{
			&map[string]int{"a": 32, "b": 64},
			&map[string]int{"a": 64, "b": 32},
			false,
		},
		{
			&map[string]int{"b": 64, "a": 32},
			&map[string]int{"a": 32, "b": 64},
			true,
		},
	}

	for _, table := range tables {
		if r := MapsEqual(table.inputA, table.inputB); r != table.expectedResult {
			t.Errorf("[STRING INT MAP] output of array equals with \"%v\" and \"%v\" "+
				"was incorrect, got: %v, want: %v.",
				table.inputA,
				table.inputB,
				r,
				table.expectedResult)
		}
	}
}

func TestMapsEqualEqualIntBoolean(t *testing.T) {
	tables := []struct {
		inputA         *map[int]bool
		inputB         *map[int]bool
		expectedResult bool
	}{
		{nil, nil, true},
		{nil, &map[int]bool{42: true, 64: false}, false},
		{&map[int]bool{42: true, 64: false}, nil, false},
		{
			&map[int]bool{5: true, 10: false, 15: true},
			&map[int]bool{5: true, 10: false},
			false,
		},
		{
			&map[int]bool{5: true, 10: false},
			&map[int]bool{5: true, 10: false, 15: true},
			false,
		},
		{
			&map[int]bool{42: true},
			&map[int]bool{42: true},
			true,
		},
		{
			&map[int]bool{42: false},
			&map[int]bool{42: false},
			true,
		},
		{
			&map[int]bool{2: true, 4: false},
			&map[int]bool{2: true, 4: false},
			true,
		},
		{
			&map[int]bool{8: true, 16: false},
			&map[int]bool{16: false, 8: true},
			true,
		},
		{
			&map[int]bool{16: false, 8: true},
			&map[int]bool{8: true, 16: false},
			true,
		},
	}

	for _, table := range tables {
		if r := MapsEqual(table.inputA, table.inputB); r != table.expectedResult {
			t.Errorf("[STRING INT MAP] output of array equals with \"%v\" and \"%v\" "+
				"was incorrect, got: %v, want: %v.",
				table.inputA,
				table.inputB,
				r,
				table.expectedResult)
		}
	}
}
