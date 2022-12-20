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
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/test/discordgo_mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ComponentTestSuite struct {
	suite.Suite
}

func (suite *ComponentTestSuite) TestLoadComponentWithSuccess() {
	hasCalled := false
	mockLoadComponentFunction := func(session *discordgo.Session) error {
		hasCalled = true

		return nil
	}

	testComponent := Component{
		Code:         "some-component",
		Name:         "Some Component",
		Description:  "This is a component!",
		LoadPriority: 100,
		State: &State{
			DefaultEnabled: true,
		},
		loadComponentFunction: mockLoadComponentFunction,
	}

	dgSession, _ := discordgo_mock.MockSession()

	err := testComponent.LoadComponent(dgSession)

	suite.NoError(err)
	suite.True(hasCalled)
	suite.Equal(dgSession, testComponent.discord)
	suite.True(testComponent.State.Loaded)
}

func (suite *ComponentTestSuite) TestLoadComponentWithFailure() {
	hasCalled := false

	mockLoadComponentFunction := func(session *discordgo.Session) error {
		hasCalled = true

		return fmt.Errorf("something bad happened")
	}

	testComponent := Component{
		Code:         "some-component",
		Name:         "Some Component",
		Description:  "This is a component!",
		LoadPriority: 100,
		State: &State{
			DefaultEnabled: true,
		},
		loadComponentFunction: mockLoadComponentFunction,
	}

	dgSession, _ := discordgo_mock.MockSession()

	err := testComponent.LoadComponent(dgSession)

	suite.Error(err)
	suite.True(hasCalled)
	suite.Equal(dgSession, testComponent.discord)
	suite.False(testComponent.State.Loaded)
}

func (suite *ComponentTestSuite) TestUnloadComponent() {
	testComponent := Component{
		Code:         "some_component",
		Name:         "Some Component",
		Description:  "This is a component!",
		LoadPriority: 100,
		State: &State{
			DefaultEnabled: true,
		},
	}

	mockHandlerManager := ComponentHandlerContainer{
		owner: &testComponent,
	}
	testComponent.handlerManager = &mockHandlerManager

	dgSession, _ := discordgo_mock.MockSession()
	handlerComponentMapping.handlers[("some_component_handler")] = &AssignedEventHandler{
		name:       "some_component_handler",
		component:  &testComponent,
		unregister: func() {},
	}

	err := testComponent.UnloadComponent(dgSession)

	suite.NoError(err)
	suite.Len(handlerComponentMapping.handlers, 0)
}

func (suite *ComponentTestSuite) TestIsCoreComponent() {
	tables := []struct {
		component      Component
		expectedResult bool
	}{
		{Component{Code: "some_random_component"}, false},
		{Component{Code: "123test"}, false},
		{Component{Code: "*some*test*"}, false},
		{Component{Code: "bot_core"}, true},
	}

	for _, table := range tables {
		result := IsCoreComponent(&table.component)

		suite.Equal(table.expectedResult, result)
	}
}

func (suite *ComponentTestSuite) TestComputeCategoriesWithOnlyInitialCategories() {
	baseCategories := Categories{CategoryInternal}

	testComponent := Component{
		Code:       "some-component",
		Categories: baseCategories,
	}

	suite.Equal(baseCategories, testComponent.Categories)
}

func (suite *ComponentTestSuite) TestComputeCategoriesWithOnlyCommandCategories() {
	expectedResult := Categories{CategoryAdministration, CategoryUtilities}

	testComponent := Component{
		Code:       "some-component",
		Categories: Categories{},
	}

	componentCommandMap = map[string]*Command{
		"a": {
			Category: CategoryAdministration,
			c:        &testComponent,
		},
		"b": {
			Category: CategoryUtilities,
			c:        &testComponent,
		},
		"c": {
			Category: CategoryAdministration,
			c:        &testComponent,
		},
	}

	testComponent.computeCategories()

	suite.Equal(expectedResult, testComponent.Categories)
}

func (suite *ComponentTestSuite) TestComputeCategoriesWithMixedCategorySource() {
	expectedResult := Categories{CategoryInternal, CategoryAdministration, CategoryUtilities}

	testComponent := Component{
		Code:       "some-component",
		Categories: Categories{CategoryInternal},
	}

	componentCommandMap = map[string]*Command{
		"a": {
			Category: CategoryAdministration,
			c:        &testComponent,
		},
		"b": {
			Category: CategoryUtilities,
			c:        &testComponent,
		},
		"c": {
			Category: CategoryAdministration,
			c:        &testComponent,
		},
	}

	testComponent.computeCategories()

	suite.Equal(expectedResult, testComponent.Categories)
}

func TestComponent(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}
