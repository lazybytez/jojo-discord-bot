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

package api

import (
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SlashCommandManagerTestSuite struct {
	suite.Suite
	owningComponent     *Component
	slashCommandManager SlashCommandManager
}

func (suite *SlashCommandManagerTestSuite) SetupTest() {
	suite.owningComponent = &Component{
		Code: "test_component",
		Name: "Test Component",
	}
	suite.slashCommandManager = SlashCommandManager{owner: suite.owningComponent}
}

func (suite *SlashCommandManagerTestSuite) TestGetCommandsForComponentWithoutMatchingCommands() {
	testComponentCode := entities.ComponentCode("no_commands_component")
	componentCommandMap = map[string]*Command{
		"a": {
			Category: CategoryAdministration,
			c:        suite.owningComponent,
		},
		"b": {
			Category: CategoryUtilities,
			c:        suite.owningComponent,
		},
		"c": {
			Category: CategoryAdministration,
			c:        suite.owningComponent,
		},
	}

	result := suite.slashCommandManager.GetCommandsForComponent(testComponentCode)

	suite.Equal([]*Command{}, result)
}

func (suite *SlashCommandManagerTestSuite) TestGetCommandsForComponentWithCommands() {
	testComponentCode := entities.ComponentCode("no_commands_component")

	testComponent := &Component{
		Code: testComponentCode,
	}

	foundCommandOne := &Command{
		Category: CategoryAdministration,
		c:        testComponent,
	}
	foundCommandTwo := &Command{
		Category: CategoryUtilities,
		c:        testComponent,
	}

	componentCommandMap = map[string]*Command{
		"a": foundCommandOne,
		"c": {
			Category: CategoryAdministration,
			c:        suite.owningComponent,
		},
		"b": foundCommandTwo,
	}

	result := suite.slashCommandManager.GetCommandsForComponent(testComponentCode)

	suite.Len(result, 2)

	suite.Equal(foundCommandOne.c.Code, result[0].c.Code)
	suite.Equal(foundCommandOne.Category, result[0].Category)

	suite.Equal(foundCommandTwo.c.Code, result[1].c.Code)
	suite.Equal(foundCommandTwo.Category, result[1].Category)
}

func TestSlashCommandManager(t *testing.T) {
	suite.Run(t, new(SlashCommandManagerTestSuite))
}
