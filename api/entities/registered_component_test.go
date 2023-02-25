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

type RegisteredComponentEntityManagerTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     entity_manager_mock.EntityManagerMock
	gem    *RegisteredComponentEntityManager
}

func (suite *RegisteredComponentEntityManagerTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = entity_manager_mock.EntityManagerMock{}
	suite.gem = &RegisteredComponentEntityManager{
		&suite.em,
		[]ComponentCode{},
	}

	err := cache.Init(cache.ModeMemory, 10*time.Minute, "")
	suite.NoError(err)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestNewRegisteredComponentEntityManager() {
	testEntityManager := entity_manager_mock.EntityManagerMock{}

	gem := NewRegisteredComponentEntityManager(&testEntityManager)

	suite.NotNil(gem)
	suite.Equal(&testEntityManager, gem.EntityManager)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestGet() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&RegisteredComponent{}).Name()),
		[]interface{}{ColumnCode + " = ?", testCode},
	).Run(func(args mock.Arguments) {
		switch v := args.Get(0).(type) {
		case *RegisteredComponent:
			v.Code = testCode
		}
	}).Return(nil).Once()

	result, err := suite.gem.Get(testCode)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testCode, result.Code)

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.True(ok)
	suite.Equal(*result, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestGetWithCache() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"

	testRegisteredComponent := &RegisteredComponent{
		Code: testCode,
	}

	err := cache.Update(testCacheKey, *testRegisteredComponent)
	suite.NoError(err)

	// Do not expect call of GetFirstEntity or DB calls
	// When GetFirstEntity is called, test will fail as call is unexpected

	result, err := suite.gem.Get(testCode)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testCode, result.Code)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestGetWithError() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"

	expectedError := fmt.Errorf("something bad happened during database read")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&RegisteredComponent{}).Name()),
		[]interface{}{ColumnCode + " = ?", testCode},
	).Return(expectedError).Once()

	result, err := suite.gem.Get(testCode)

	suite.dba.AssertExpectations(suite.T())
	suite.Error(err)
	suite.NotNil(result)
	suite.Equal(*result, RegisteredComponent{})

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestCreate() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"
	testRegisteredComponent := RegisteredComponent{
		Code: testCode,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testRegisteredComponent).Return(nil).Once()

	err := suite.gem.Create(&testRegisteredComponent)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestCreateWithError() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"
	testRegisteredComponent := RegisteredComponent{
		Code: testCode,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testRegisteredComponent).Return(expectedErr).Once()

	err := suite.gem.Create(&testRegisteredComponent)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestSave() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"
	testRegisteredComponent := RegisteredComponent{
		Code: testCode,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testRegisteredComponent).Return(nil).Once()

	err := suite.gem.Save(&testRegisteredComponent)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestSaveWithError() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"
	testRegisteredComponent := RegisteredComponent{
		Code: testCode,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testRegisteredComponent).Return(expectedErr).Once()

	err := suite.gem.Save(&testRegisteredComponent)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestUpdate() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"
	testRegisteredComponent := RegisteredComponent{
		Code: testCode,
	}
	testColumn := "some_column"
	testValue := "some_value"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testRegisteredComponent, testColumn, testValue).Return(nil).Once()

	err := suite.gem.Update(&testRegisteredComponent, testColumn, testValue)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func (suite *RegisteredComponentEntityManagerTestSuite) TestUpdateWithError() {
	testCode := ComponentCode("very_important_component")
	testCacheKey := "very_important_component"
	testRegisteredComponent := RegisteredComponent{
		Code: testCode,
	}
	testColumn := "some_column"
	testValue := "some_value"

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testRegisteredComponent, testColumn, testValue).Return(expectedErr).Once()

	err := suite.gem.Update(&testRegisteredComponent, testColumn, testValue)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(testCacheKey, RegisteredComponent{})
	suite.False(ok)
	suite.Equal(RegisteredComponent{}, cachedRegisteredComponent)
}

func TestRegisteredComponentEntityManager(t *testing.T) {
	suite.Run(t, new(RegisteredComponentEntityManagerTestSuite))
}
