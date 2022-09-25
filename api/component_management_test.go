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
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/suite"
	"testing"
)

// RegisterComponentTestSuite tests RegisterComponent explicitly
// and sortComponents implicitly!
type RegisterComponentTestSuite struct {
	suite.Suite
}

func (suite *RegisterComponentTestSuite) SetupTest() {
	// Ensure empty map before every test
	Components = []*Component{}
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithSingleFeatureComponent() {
	testComponent := Component{
		// Metadata
		Code: "component1",
	}

	hasCalled := false
	loadFunc := func(session *discordgo.Session) error {
		hasCalled = true

		return nil
	}

	RegisterComponent(&testComponent, loadFunc)

	suite.Len(Components, 1)

	resultComp := Components[0]

	suite.NotNil(resultComp)
	suite.Equal(&testComponent, resultComp)

	err := resultComp.loadComponentFunction(nil)
	suite.NoError(err)
	suite.True(hasCalled)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleUnsortedFeatureComponents() {
	testComponent1 := Component{
		// Metadata
		Code: "component1",
	}
	testComponent2 := Component{
		// Metadata
		Code: "component2",
	}
	testComponent3 := Component{
		// Metadata
		Code: "component3",
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)

	suite.Len(Components, 3)

	// Fist component
	resultComp1 := Components[0]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component
	resultComp2 := Components[1]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component
	resultComp3 := Components[2]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleSortedFeatureComponents() {
	testComponent1 := Component{
		// Metadata
		Code:         "component1",
		LoadPriority: -1000,
	}
	testComponent2 := Component{
		// Metadata
		Code: "component2",
	}
	testComponent3 := Component{
		// Metadata
		Code:         "component3",
		LoadPriority: 1000,
	}
	testComponent4 := Component{
		// Metadata
		Code: "component4",
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}
	hasCalled4 := false
	loadFunc4 := func(session *discordgo.Session) error {
		hasCalled4 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)
	RegisterComponent(&testComponent4, loadFunc4)

	suite.Len(Components, 4)

	// Fist component -> Third component
	resultComp1 := Components[3]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component
	resultComp2 := Components[1]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component -> First component
	resultComp3 := Components[0]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)

	// Fourth component
	resultComp4 := Components[2]
	suite.NotNil(resultComp4)
	suite.Equal(&testComponent4, resultComp4)
	err4 := resultComp4.loadComponentFunction(nil)
	suite.NoError(err4)
	suite.True(hasCalled4)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithSingleCoreComponent() {
	testComponent := Component{
		// Metadata
		Code: "bot_component1",
	}

	hasCalled := false
	loadFunc := func(session *discordgo.Session) error {
		hasCalled = true

		return nil
	}

	RegisterComponent(&testComponent, loadFunc)

	suite.Len(Components, 1)

	resultComp := Components[0]

	suite.NotNil(resultComp)
	suite.Equal(&testComponent, resultComp)

	err := resultComp.loadComponentFunction(nil)
	suite.NoError(err)
	suite.True(hasCalled)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleUnsortedCoreComponents() {
	testComponent1 := Component{
		// Metadata
		Code: "bot_component1",
	}
	testComponent2 := Component{
		// Metadata
		Code: "bot_component2",
	}
	testComponent3 := Component{
		// Metadata
		Code: "bot_component3",
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)

	suite.Len(Components, 3)

	// Fist component
	resultComp1 := Components[0]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component
	resultComp2 := Components[1]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component
	resultComp3 := Components[2]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleSortedCoreComponents() {
	testComponent1 := Component{
		// Metadata
		Code:         "bot_component1",
		LoadPriority: -1000,
	}
	testComponent2 := Component{
		// Metadata
		Code: "bot_component2",
	}
	testComponent3 := Component{
		// Metadata
		Code:         "bot_component3",
		LoadPriority: 1000,
	}
	testComponent4 := Component{
		// Metadata
		Code: "bot_component4",
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}
	hasCalled4 := false
	loadFunc4 := func(session *discordgo.Session) error {
		hasCalled4 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)
	RegisterComponent(&testComponent4, loadFunc4)

	suite.Len(Components, 4)

	// Fist component -> Fourth component
	resultComp1 := Components[3]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component -> Second component
	resultComp2 := Components[1]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component -> First component
	resultComp3 := Components[0]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)

	// Fourth component -> Third component
	resultComp4 := Components[2]
	suite.NotNil(resultComp4)
	suite.Equal(&testComponent4, resultComp4)
	err4 := resultComp4.loadComponentFunction(nil)
	suite.NoError(err4)
	suite.True(hasCalled4)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleUnsortedMixedComponents() {
	testComponent1 := Component{
		// Metadata
		Code: "component1",
	}
	testComponent2 := Component{
		// Metadata
		Code: "bot_component2",
	}
	testComponent3 := Component{
		// Metadata
		Code: "component3",
	}
	testComponent4 := Component{
		// Metadata
		Code: "bot_component4",
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}
	hasCalled4 := false
	loadFunc4 := func(session *discordgo.Session) error {
		hasCalled4 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)
	RegisterComponent(&testComponent4, loadFunc4)

	suite.Len(Components, 4)

	// Fist component -> Third component
	resultComp1 := Components[2]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component -> First component
	resultComp2 := Components[0]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component -> Fourth component
	resultComp3 := Components[3]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)

	// Fourth component -> Second component
	resultComp4 := Components[1]
	suite.NotNil(resultComp4)
	suite.Equal(&testComponent4, resultComp4)
	err4 := resultComp4.loadComponentFunction(nil)
	suite.NoError(err4)
	suite.True(hasCalled4)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleSortedBotMixedComponents() {
	testComponent1 := Component{
		// Metadata
		Code: "bot_component1",
	}
	testComponent2 := Component{
		// Metadata
		Code:         "bot_component2",
		LoadPriority: -999,
	}
	testComponent3 := Component{
		// Metadata
		Code: "component3",
	}
	testComponent4 := Component{
		// Metadata
		Code:         "bot_component4",
		LoadPriority: 9999,
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}
	hasCalled4 := false
	loadFunc4 := func(session *discordgo.Session) error {
		hasCalled4 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)
	RegisterComponent(&testComponent4, loadFunc4)

	suite.Len(Components, 4)

	// Fist component -> Second component
	resultComp1 := Components[1]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component -> Fourth component
	resultComp2 := Components[2]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component -> Third component
	resultComp3 := Components[3]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)

	// Fourth component -> First component
	resultComp4 := Components[0]
	suite.NotNil(resultComp4)
	suite.Equal(&testComponent4, resultComp4)
	err4 := resultComp4.loadComponentFunction(nil)
	suite.NoError(err4)
	suite.True(hasCalled4)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleSortedFeatureMixedComponents() {
	testComponent1 := Component{
		// Metadata
		Code: "component1",
	}
	testComponent2 := Component{
		// Metadata
		Code:         "component2",
		LoadPriority: -999,
	}
	testComponent3 := Component{
		// Metadata
		Code: "bot_component3",
	}
	testComponent4 := Component{
		// Metadata
		Code:         "component4",
		LoadPriority: 9999,
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}
	hasCalled4 := false
	loadFunc4 := func(session *discordgo.Session) error {
		hasCalled4 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)
	RegisterComponent(&testComponent4, loadFunc4)

	suite.Len(Components, 4)

	// Fist component -> Second component
	resultComp1 := Components[2]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component -> Fourth component
	resultComp2 := Components[3]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component -> Third component
	resultComp3 := Components[0]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)

	// Fourth component -> First component
	resultComp4 := Components[1]
	suite.NotNil(resultComp4)
	suite.Equal(&testComponent4, resultComp4)
	err4 := resultComp4.loadComponentFunction(nil)
	suite.NoError(err4)
	suite.True(hasCalled4)
}

func (suite *RegisterComponentTestSuite) TestRegisterComponentWithMultipleSortedMixedComponents() {
	testComponent1 := Component{
		// Metadata
		Code:         "bot_component1",
		LoadPriority: -99,
	}
	testComponent2 := Component{
		// Metadata
		Code:         "component2",
		LoadPriority: -999,
	}
	testComponent3 := Component{
		// Metadata
		Code: "bot_component3",
	}
	testComponent4 := Component{
		// Metadata
		Code:         "component4",
		LoadPriority: 9999,
	}
	testComponent5 := Component{
		// Metadata
		Code:         "bot_component5",
		LoadPriority: 9999,
	}
	testComponent6 := Component{
		// Metadata
		Code: "component6",
	}

	hasCalled1 := false
	loadFunc1 := func(session *discordgo.Session) error {
		hasCalled1 = true
		return nil
	}
	hasCalled2 := false
	loadFunc2 := func(session *discordgo.Session) error {
		hasCalled2 = true
		return nil
	}
	hasCalled3 := false
	loadFunc3 := func(session *discordgo.Session) error {
		hasCalled3 = true
		return nil
	}
	hasCalled4 := false
	loadFunc4 := func(session *discordgo.Session) error {
		hasCalled4 = true
		return nil
	}
	hasCalled5 := false
	loadFunc5 := func(session *discordgo.Session) error {
		hasCalled5 = true
		return nil
	}
	hasCalled6 := false
	loadFunc6 := func(session *discordgo.Session) error {
		hasCalled6 = true
		return nil
	}

	RegisterComponent(&testComponent1, loadFunc1)
	RegisterComponent(&testComponent2, loadFunc2)
	RegisterComponent(&testComponent3, loadFunc3)
	RegisterComponent(&testComponent4, loadFunc4)
	RegisterComponent(&testComponent5, loadFunc5)
	RegisterComponent(&testComponent6, loadFunc6)

	suite.Len(Components, 6)

	// Fist component -> Third component
	resultComp1 := Components[2]
	suite.NotNil(resultComp1)
	suite.Equal(&testComponent1, resultComp1)
	err1 := resultComp1.loadComponentFunction(nil)
	suite.NoError(err1)
	suite.True(hasCalled1)

	// Second component -> Fith component
	resultComp2 := Components[5]
	suite.NotNil(resultComp2)
	suite.Equal(&testComponent2, resultComp2)
	err2 := resultComp2.loadComponentFunction(nil)
	suite.NoError(err2)
	suite.True(hasCalled2)

	// Third component -> Second component
	resultComp3 := Components[1]
	suite.NotNil(resultComp3)
	suite.Equal(&testComponent3, resultComp3)
	err3 := resultComp3.loadComponentFunction(nil)
	suite.NoError(err3)
	suite.True(hasCalled3)

	// Fourth component -> Fourth component
	resultComp4 := Components[3]
	suite.NotNil(resultComp4)
	suite.Equal(&testComponent4, resultComp4)
	err4 := resultComp4.loadComponentFunction(nil)
	suite.NoError(err4)
	suite.True(hasCalled4)

	// Fifth component -> First component
	resultComp5 := Components[0]
	suite.NotNil(resultComp5)
	suite.Equal(&testComponent5, resultComp5)
	err5 := resultComp5.loadComponentFunction(nil)
	suite.NoError(err5)
	suite.True(hasCalled5)

	// Sixth component -> Fourth component
	resultComp6 := Components[4]
	suite.NotNil(resultComp6)
	suite.Equal(&testComponent6, resultComp6)
	err6 := resultComp6.loadComponentFunction(nil)
	suite.NoError(err6)
	suite.True(hasCalled6)
}

func TestRegisterComponent(t *testing.T) {
	suite.Run(t, new(RegisterComponentTestSuite))
}

type GetGuildIdFromEventInterfaceTestSuite struct {
	suite.Suite
}

func (suite *GetGuildIdFromEventInterfaceTestSuite) TestGetGuildIdFromEventInterfaceWithPossibleOptions() {
	tables := []struct {
		input          interface{}
		expectedOutput string
	}{
		{discordgo.GuildCreate{Guild: &discordgo.Guild{ID: "this-is-some-id"}}, "this-is-some-id"},
		{&discordgo.GuildCreate{Guild: &discordgo.Guild{ID: "this-is-some-id"}}, "this-is-some-id"},
		{discordgo.MessageCreate{Message: &discordgo.Message{GuildID: "this-is-some-id"}}, "this-is-some-id"},
		{&discordgo.MessageCreate{Message: &discordgo.Message{GuildID: "this-is-some-id"}}, "this-is-some-id"},
		{
			struct {
				SomeField string
			}{SomeField: "some-test"},
			"",
		},
	}

	for _, table := range tables {
		result := getGuildIdFromEventInterface(table.input)

		suite.Equal(table.expectedOutput, result)
	}
}

func TestGetGuildIdFromEventInterface(t *testing.T) {
	suite.Run(t, new(GetGuildIdFromEventInterfaceTestSuite))
}
