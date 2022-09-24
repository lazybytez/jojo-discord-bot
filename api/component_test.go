package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/test/discordgo_mock"
	"github.com/lazybytez/jojo-discord-bot/test/handler_manager_mock"
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
	mockHandlerManager := handler_manager_mock.HandlerManagerMock{}

	testComponent := Component{
		Code:         "some-component",
		Name:         "Some Component",
		Description:  "This is a component!",
		LoadPriority: 100,
		State: &State{
			DefaultEnabled: true,
		},
		handlerManager: &mockHandlerManager,
	}

	dgSession, _ := discordgo_mock.MockSession()

	mockHandlerManager.On("UnregisterAll").Once()

	err := testComponent.UnloadComponent(dgSession)

	suite.NoError(err)
	mockHandlerManager.AssertExpectations(suite.T())
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

func TestComponent(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}
