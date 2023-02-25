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

type GuildEntityManagerTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     entity_manager_mock.EntityManagerMock
	gem    *GuildEntityManager
}

func (suite *GuildEntityManagerTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = entity_manager_mock.EntityManagerMock{}
	suite.gem = &GuildEntityManager{
		&suite.em,
	}

	cache.Init(cache.ModeMemory, 10*time.Minute, "")
}

func (suite *GuildEntityManagerTestSuite) TestNewGuildEntityManager() {
	testEntityManager := entity_manager_mock.EntityManagerMock{}

	gem := NewGuildEntityManager(&testEntityManager)

	suite.NotNil(gem)
	suite.Equal(&testEntityManager, gem.EntityManager)
}

func (suite *GuildEntityManagerTestSuite) TestGet() {
	testIdString := "652658256236529525"
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&Guild{}).Name()),
		[]interface{}{ColumnGuildId + " = ?", testId},
	).Run(func(args mock.Arguments) {
		switch v := args.Get(0).(type) {
		case *Guild:
			v.GuildID = testId
		}
	}).Return(nil).Once()

	result, err := suite.gem.Get(testIdString)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testId, result.GuildID)

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.True(ok)
	suite.Equal(*result, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestGetWithCache() {
	testIdString := "652658256236529525"
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := &Guild{
		GuildID: testId,
	}

	err := cache.Update(testCacheKey, *testGuild)
	suite.NoError(err)

	// Do not expect call of GetFirstEntity or DB calls
	// When GetFirstEntity is called, test will fail as call is unexpected

	result, err := suite.gem.Get(testIdString)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testId, result.GuildID)
}

func (suite *GuildEntityManagerTestSuite) TestGetWithError() {
	testIdString := "652658256236529525"
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"

	expectedError := fmt.Errorf("something bad happened during database read")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&Guild{}).Name()),
		[]interface{}{ColumnGuildId + " = ?", testId},
	).Return(expectedError).Once()

	result, err := suite.gem.Get(testIdString)

	suite.dba.AssertExpectations(suite.T())
	suite.Error(err)
	suite.NotNil(result)
	suite.Equal(*result, Guild{})

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestCreate() {
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := Guild{
		GuildID: testId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testGuild).Return(nil).Once()

	err := suite.gem.Create(&testGuild)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestCreateWithError() {
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := Guild{
		GuildID: testId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testGuild).Return(expectedErr).Once()

	err := suite.gem.Create(&testGuild)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestSave() {
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := Guild{
		GuildID: testId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testGuild).Return(nil).Once()

	err := suite.gem.Save(&testGuild)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestSaveWithError() {
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := Guild{
		GuildID: testId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testGuild).Return(expectedErr).Once()

	err := suite.gem.Save(&testGuild)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestUpdate() {
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := Guild{
		GuildID: testId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testGuild, testColumn, testValue).Return(nil).Once()

	err := suite.gem.Update(&testGuild, testColumn, testValue)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func (suite *GuildEntityManagerTestSuite) TestUpdateWithError() {
	testId := uint64(652658256236529525)
	testCacheKey := "652658256236529525"
	testGuild := Guild{
		GuildID: testId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testGuild, testColumn, testValue).Return(expectedErr).Once()

	err := suite.gem.Update(&testGuild, testColumn, testValue)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedGuild, ok := cache.Get(testCacheKey, Guild{})
	suite.False(ok)
	suite.Equal(Guild{}, cachedGuild)
}

func TestGuildEntityManager(t *testing.T) {
	suite.Run(t, new(GuildEntityManagerTestSuite))
}
