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

package helper

import (
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ComponentHelperTestSuite struct {
	suite.Suite
}

func (suite *ComponentHelperTestSuite) TestTestIfComponentIsRegisteredWithNoComponentFound() {
	testComponent := &api.Component{
		Code: "some-component",
	}

	result := TestIfComponentIsRegistered(testComponent)

	assert.False(suite.T(), result)
}

func (suite *ComponentHelperTestSuite) TestTestIfComponentIsRegisteredWithComponentFound() {
	testComponent := &api.Component{
		Code: "some-component",
	}

	api.Components = append(api.Components, testComponent)

	result := TestIfComponentIsRegistered(testComponent)

	assert.True(suite.T(), result)
}

func TestComponentHelper(t *testing.T) {
	suite.Run(t, new(ComponentHelperTestSuite))
}
