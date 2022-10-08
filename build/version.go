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

package build

import "fmt"

// These variables are injected during build time.
//
// Version holds the currently running version.
// This either matches the Git tag on which the build is based
// or "edge" for any other build.
//
// If Version is "edge", the CommitSHA should be populated too.
// It should always keep the SHA of the commit that was used to build
// this version.

var (
	Version   string
	CommitSHA string
)

// init handles the case when no values were given during build
func init() {
	if Version == "" {
		Version = "edge"
	}

	if CommitSHA == "" {
		CommitSHA = "unknown"
	}
}

// ComputeVersionString returns the current version in a displayable format
func ComputeVersionString() string {
	if Version == "edge" {
		return fmt.Sprintf("%s<%s>", Version, CommitSHA)
	}

	return Version
}
