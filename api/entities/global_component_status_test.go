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

package entities

import (
	"fmt"
	"github.com/lazybytez/jojo-discord-bot/services/cache"
	"github.com/lazybytez/jojo-discord-bot/test/dbmock"
	"github.com/lazybytez/jojo-discord-bot/test/entity_manager_mock"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type GlobalComponentStatusEntityManagerTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     entity_manager_mock.EntityManagerMock
	gem    *GlobalComponentStatusEntityManager
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = entity_manager_mock.EntityManagerMock{}
	suite.gem = &GlobalComponentStatusEntityManager{
		&suite.em,
	}

	cache.Init(cache.ModeMemory, 10*time.Minute, "")
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestNewGlobalComponentStatusEntityManager() {
	testEntityManager := entity_manager_mock.EntityManagerMock{}

	gem := NewGlobalComponentStatusEntityManager(&testEntityManager)

	suite.NotNil(gem)
	suite.Equal(&testEntityManager, gem.EntityManager)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestGet() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&GlobalComponentStatus{}).Name()),
		[]interface{}{ColumnComponent + " = ?", testId},
	).Run(func(args mock.Arguments) {
		switch v := args.Get(0).(type) {
		case *GlobalComponentStatus:
			v.ComponentID = testId
		}
	}).Return(nil).Once()

	result, err := suite.gem.Get(testId)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testId, result.ComponentID)

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.True(ok)
	suite.Equal(*result, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestGetWithCache() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := &GlobalComponentStatus{
		ComponentID: testId,
	}

	cache.Update(testCacheKey, *testGlobalComponentStatus)

	// Do not expect call of GetFirstEntity or DB calls
	// When GetFirstEntity is called, test will fail as call is unexpected

	result, err := suite.gem.Get(testId)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testId, result.ComponentID)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestGetWithError() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	expectedError := fmt.Errorf("something bad happened during database read")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&GlobalComponentStatus{}).Name()),
		[]interface{}{ColumnComponent + " = ?", testId},
	).Return(expectedError).Once()

	result, err := suite.gem.Get(testId)

	suite.dba.AssertExpectations(suite.T())
	suite.Error(err)
	suite.NotNil(result)
	suite.Equal(*result, GlobalComponentStatus{})

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestCreate() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := GlobalComponentStatus{
		ComponentID: testId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testGlobalComponentStatus).Return(nil).Once()

	err := suite.gem.Create(&testGlobalComponentStatus)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestCreateWithError() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := GlobalComponentStatus{
		ComponentID: testId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testGlobalComponentStatus).Return(expectedErr).Once()

	err := suite.gem.Create(&testGlobalComponentStatus)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestSave() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := GlobalComponentStatus{
		ComponentID: testId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testGlobalComponentStatus).Return(nil).Once()

	err := suite.gem.Save(&testGlobalComponentStatus)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestSaveWithError() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := GlobalComponentStatus{
		ComponentID: testId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testGlobalComponentStatus).Return(expectedErr).Once()

	err := suite.gem.Save(&testGlobalComponentStatus)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestUpdate() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := GlobalComponentStatus{
		ComponentID: testId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testGlobalComponentStatus, testColumn, testValue).Return(nil).Once()

	err := suite.gem.Update(&testGlobalComponentStatus, testColumn, testValue)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func (suite *GlobalComponentStatusEntityManagerTestSuite) TestUpdateWithError() {
	testId := uint(65835858358583)
	testCacheKey := "65835858358583"
	testGlobalComponentStatus := GlobalComponentStatus{
		ComponentID: testId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testGlobalComponentStatus, testColumn, testValue).Return(expectedErr).Once()

	err := suite.gem.Update(&testGlobalComponentStatus, testColumn, testValue)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedGlobalComponentStatus, ok := cache.Get(testCacheKey, GlobalComponentStatus{})
	suite.False(ok)
	suite.Equal(GlobalComponentStatus{}, cachedGlobalComponentStatus)
}

func TestGlobalComponentStatusEntityManager(t *testing.T) {
	suite.Run(t, new(GlobalComponentStatusEntityManagerTestSuite))
}
