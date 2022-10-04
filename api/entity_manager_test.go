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
	"github.com/lazybytez/jojo-discord-bot/test/db"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EntitiesTestSuite struct {
	suite.Suite
	dba    *db.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     EntityManager
}

func (suite *EntitiesTestSuite) SetupTest() {
	dba := &db.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = NewEntityManager(suite.dba, suite.logger)
}

func (suite *EntitiesTestSuite) TestRegisterDefaultEntitiesWithSuccess() {
	for _, entity := range defaultEntities {
		suite.dba.On("RegisterEntity", entity).Once().Return(nil)
	}

	err := suite.em.RegisterDefaultEntities()

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
}

func (suite *EntitiesTestSuite) TestRegisterDefaultEntitiesWithFailure() {
	simulatedErr := fmt.Errorf("something bad happened")

	suite.dba.On("RegisterEntity", mock.Anything).Once().Return(simulatedErr)

	err := suite.em.RegisterDefaultEntities()

	suite.dba.AssertExpectations(suite.T())
	suite.Error(err)
}

func TestEntities(t *testing.T) {
	suite.Run(t, new(EntitiesTestSuite))
}
