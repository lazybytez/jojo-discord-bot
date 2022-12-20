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

// InArrayOrSlice checks if the provided element is in
// the provided sliceOrArray.
func InArrayOrSlice[T comparable](sliceOrArray []T, element T) bool {
	for _, elem := range sliceOrArray {
		if elem == element {
			return true
		}
	}

	return false
}

// UniqueArrayOrSlice returns the passed slice or array without any duplicates in it.
func UniqueArrayOrSlice[T comparable](sliceOrArray []T) []T {
	unique := make([]T, 0)

	for _, elem := range sliceOrArray {
		if !InArrayOrSlice(unique, elem) {
			unique = append(unique, elem)
		}
	}

	return unique
}
