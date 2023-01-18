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

// ArraysEqual checks if two arrays are equal and returns the result.
// Note that this function requires comparable types!
//
// The arguments are passed by reference to allow nil
// comparisons.
// When nil is passed, true will be returned if both slice references are nil.
// If only one slice is nil, the result will be false.
//
// Note that the order of elements is also checked.
func ArraysEqual[T comparable](a *[]T, b *[]T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	aV := *a
	bV := *b
	if len(aV) != len(bV) {
		return false
	}

	for k, val := range aV {
		if val != bV[k] {
			return false
		}
	}

	return true
}

// MapsEqual checks if two maps have the same content
// Note that this function requires comparable types for key and value!
//
// The arguments are passed by reference to allow nil
// comparisons.
// When nil is passed, true will be returned if both maps references are nil.
// If only one map reference is nil, the result will be false.
//
// Note that the order of elements is ignored.
func MapsEqual[K comparable, V comparable](a *map[K]V, b *map[K]V) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	aV := *a
	bV := *b
	if len(aV) != len(bV) {
		return false
	}

	for k, valA := range aV {
		if valB, ok := bV[k]; !ok || valA != valB {
			return false
		}
	}

	return true
}

// PointerValuesEqual compares two comparable variables which values are pointers.
// The function first checks for nil and then dereferences the pointers.
func PointerValuesEqual[T comparable](a *T, b *T) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}
