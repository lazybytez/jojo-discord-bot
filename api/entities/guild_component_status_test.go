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
	"github.com/lazybytez/jojo-discord-bot/api/cache"
	"github.com/lazybytez/jojo-discord-bot/test/dbmock"
	"github.com/lazybytez/jojo-discord-bot/test/entity_manager_mock"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type GuildComponentStatusEntityManagerTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     entity_manager_mock.EntityManagerMock
	gem    *GuildComponentStatusEntityManager
}

func (suite *GuildComponentStatusEntityManagerTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = entity_manager_mock.EntityManagerMock{}
	suite.gem = &GuildComponentStatusEntityManager{
		&suite.em,
		cache.New[string, GuildComponentStatus](5 * time.Second),
	}
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestGetCacheKey() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)

	expectedCacheKey := "65835858358583_48688742646283"

	result := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	suite.Equal(expectedCacheKey, result)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestNewGuildComponentStatusEntityManager() {
	testEntityManager := entity_manager_mock.EntityManagerMock{}

	gem := NewGuildComponentStatusEntityManager(&testEntityManager)

	suite.NotNil(gem)
	suite.NotNil(gem.cache)
	suite.Equal(&testEntityManager, gem.EntityManager)

	gem.cache.DisableAutoCleanup()
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestGet() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&GuildComponentStatus{}).Name()),
		[]interface{}{ColumnGuild + " = ? AND " + ColumnComponent + " = ?", guildId, componentId},
	).Run(func(args mock.Arguments) {
		switch v := args.Get(0).(type) {
		case *GuildComponentStatus:
			v.GuildID = guildId
			v.ComponentID = componentId
		}
	}).Return(nil).Once()

	result, err := suite.gem.Get(guildId, componentId)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(guildId, result.GuildID)
	suite.Equal(componentId, result.ComponentID)

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)

	suite.True(ok)
	suite.Equal(result, cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestGetWithCache() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)

	testGuildComponentStatus := &GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cache.Update(suite.gem.cache, cacheKey, testGuildComponentStatus)

	// Do not expect call of GetFirstEntity or DB calls
	// When GetFirstEntity is called, test will fail as call is unexpected

	result, err := suite.gem.Get(guildId, componentId)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(guildId, result.GuildID)
	suite.Equal(componentId, result.ComponentID)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestGetWithError() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)

	expectedError := fmt.Errorf("something bad happened during database read")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&GuildComponentStatus{}).Name()),
		[]interface{}{ColumnGuild + " = ? AND " + ColumnComponent + " = ?", guildId, componentId},
	).Return(expectedError).Once()

	result, err := suite.gem.Get(guildId, componentId)

	suite.dba.AssertExpectations(suite.T())
	suite.Error(err)
	suite.NotNil(result)
	suite.Equal(*result, GuildComponentStatus{})

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.False(ok)
	suite.Nil(cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestCreate() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)

	testGuildComponentStatus := GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testGuildComponentStatus).Return(nil).Once()

	err := suite.gem.Create(&testGuildComponentStatus)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.True(ok)
	suite.Equal(&testGuildComponentStatus, cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestCreateWithError() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)
	testGuildComponentStatus := GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testGuildComponentStatus).Return(expectedErr).Once()

	err := suite.gem.Create(&testGuildComponentStatus)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.False(ok)
	suite.Nil(cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestSave() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)
	testGuildComponentStatus := GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testGuildComponentStatus).Return(nil).Once()

	err := suite.gem.Save(&testGuildComponentStatus)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.True(ok)
	suite.Equal(&testGuildComponentStatus, cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestSaveWithError() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)
	testGuildComponentStatus := GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testGuildComponentStatus).Return(expectedErr).Once()

	err := suite.gem.Save(&testGuildComponentStatus)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.False(ok)
	suite.Nil(cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestUpdate() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)
	testGuildComponentStatus := GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testGuildComponentStatus, testColumn, testValue).Return(nil).Once()

	err := suite.gem.Update(&testGuildComponentStatus, testColumn, testValue)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.True(ok)
	suite.Equal(&testGuildComponentStatus, cachedGuildComponentStatus)
}

func (suite *GuildComponentStatusEntityManagerTestSuite) TestUpdateWithError() {
	guildId := uint(65835858358583)
	componentId := uint(48688742646283)
	testGuildComponentStatus := GuildComponentStatus{
		GuildID:     guildId,
		ComponentID: componentId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testGuildComponentStatus, testColumn, testValue).Return(expectedErr).Once()

	err := suite.gem.Update(&testGuildComponentStatus, testColumn, testValue)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cacheKey := suite.gem.getComponentStatusCacheKey(guildId, componentId)

	cachedGuildComponentStatus, ok := cache.Get(suite.gem.cache, cacheKey)
	suite.False(ok)
	suite.Nil(cachedGuildComponentStatus)
}

func TestGuildComponentStatusEntityManager(t *testing.T) {
	suite.Run(t, new(GuildComponentStatusEntityManagerTestSuite))
}
