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
	"github.com/lazybytez/jojo-discord-bot/test/dbmock"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type APITestSuite struct {
	suite.Suite
}

func (suite *APITestSuite) TestInit() {
	databaseMock := &dbmock.DatabaseAccessMock{}
	loggerMock := &logmock.LoggerMock{}

	testEntityManager := EntityManager{
		database: databaseMock,
		logger:   loggerMock,
	}

	databaseMock.On("RegisterEntity", mock.Anything).Return(nil).Times(len(defaultEntities))
	loggerMock.On("Info",
		mock.AnythingOfType("string"),
		mock.AnythingOfType(reflect.TypeOf([]interface{}{}).Name()))

	err := Init(testEntityManager)

	databaseMock.AssertExpectations(suite.T())
	suite.NoError(err)
}

func (suite *APITestSuite) TestInitWithFailure() {
	databaseMock := &dbmock.DatabaseAccessMock{}
	loggerMock := &logmock.LoggerMock{}

	testError := fmt.Errorf("something really bad happened")

	testEntityManager := EntityManager{
		database: databaseMock,
		logger:   loggerMock,
	}

	databaseMock.On("RegisterEntity", mock.Anything).Return(testError).Once()
	loggerMock.On("Err",
		testError,
		mock.AnythingOfType("string"),
		mock.AnythingOfType(reflect.TypeOf([]interface{}{}).Name()))

	err := Init(testEntityManager)

	databaseMock.AssertExpectations(suite.T())
	suite.Error(err)
}

func TestApi(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
