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
