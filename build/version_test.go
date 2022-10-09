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

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type VersionTestSuite struct {
	suite.Suite
}

func (suite *VersionTestSuite) TestComputeVersionStringWithEdgeVersionAndSHA() {
	Version = "edge"
	CommitSHA = "12345543321"

	result := ComputeVersionString()

	suite.Equal("edge<12345543321>", result)
}

func (suite *VersionTestSuite) TestComputeVersionStringWithTaggedVersion() {
	Version = "v1.5.3"

	result := ComputeVersionString()

	suite.Equal("v1.5.3", result)
}

func TestVersion(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}
