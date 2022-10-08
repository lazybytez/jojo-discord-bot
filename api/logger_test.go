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
	"github.com/lazybytez/jojo-discord-bot/services/logger"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (suite *LoggerTestSuite) TestGetLoggerWithExistingLogger() {
	testLogger := &logmock.LoggerMock{}

	testComponent := Component{
		logger: testLogger,
	}

	result := testComponent.Logger()

	suite.NotNil(result)
	suite.Equal(testLogger, result)
}

func (suite *LoggerTestSuite) TestGetLoggerWithNoLogger() {
	testComponent := Component{}

	result := testComponent.Logger()
	result2 := testComponent.Logger()

	// First call
	suite.NotNil(result)
	suite.IsType(&logger.Logger{}, result)

	// Consecutive call(s)
	suite.Equal(result, result2)
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
