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
	"sync"
	"testing"
)

type StatusManagerSuite struct {
	suite.Suite
}

func (suite *StatusManagerSuite) SetupTest() {
	botStatusManager = &DiscordGoStatusManager{
		mu:      sync.RWMutex{},
		status:  make([]SimpleBotStatus, 0),
		current: 0,
	}
}

func (suite *StatusManagerSuite) TestNextWithNoStatus() {
	for i := 0; i < 10; i++ {
		suite.Nil(botStatusManager.Next())
	}
}

func (suite *StatusManagerSuite) TestNextWithOneStatus() {
	firstStatus := SimpleBotStatus{
		ActivityType: discordgo.ActivityTypeGame,
		Content:      "Test",
	}

	botStatusManager.AddStatusToRotation(firstStatus)

	// Cycle five times
	for i := 0; i < 5; i++ {
		suite.Equal(firstStatus, *botStatusManager.Next())
	}
}

func (suite *StatusManagerSuite) TestNextWithMultipleStatus() {
	firstStatus := SimpleBotStatus{
		ActivityType: discordgo.ActivityTypeGame,
		Content:      "Test",
	}

	secondStatus := SimpleBotStatus{
		ActivityType: discordgo.ActivityTypeStreaming,
		Url:          "https://localhost:8080/",
	}

	thirdStatus := SimpleBotStatus{
		ActivityType: discordgo.ActivityTypeListening,
		Content:      "Roundabout",
	}

	botStatusManager.AddStatusToRotation(firstStatus)
	botStatusManager.AddStatusToRotation(secondStatus)
	botStatusManager.AddStatusToRotation(thirdStatus)

	// First cycle
	suite.Equal(firstStatus, *botStatusManager.Next())
	suite.Equal(secondStatus, *botStatusManager.Next())
	suite.Equal(thirdStatus, *botStatusManager.Next())
	// Second cycle
	suite.Equal(firstStatus, *botStatusManager.Next())
	suite.Equal(secondStatus, *botStatusManager.Next())
	suite.Equal(thirdStatus, *botStatusManager.Next())
	// Third cycle
	suite.Equal(firstStatus, *botStatusManager.Next())
	suite.Equal(secondStatus, *botStatusManager.Next())
	suite.Equal(thirdStatus, *botStatusManager.Next())
	// Fourth cycle
	suite.Equal(firstStatus, *botStatusManager.Next())
	suite.Equal(secondStatus, *botStatusManager.Next())
	suite.Equal(thirdStatus, *botStatusManager.Next())
}

func TestStatusManager(t *testing.T) {
	suite.Run(t, new(StatusManagerSuite))
}
