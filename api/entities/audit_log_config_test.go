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

type AuditLogConfigEntityManagerTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     entity_manager_mock.EntityManagerMock
	gem    *AuditLogConfigEntityManager
}

func (suite *AuditLogConfigEntityManagerTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = entity_manager_mock.EntityManagerMock{}
	suite.gem = &AuditLogConfigEntityManager{
		&suite.em,
		cache.New[uint64, AuditLogConfig](5 * time.Second),
	}
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestNewAuditLogConfigEntityManager() {
	testEntityManager := entity_manager_mock.EntityManagerMock{}

	gem := NewAuditLogConfigEntityManager(&testEntityManager)

	suite.NotNil(gem)
	suite.NotNil(gem.cache)
	suite.Equal(&testEntityManager, gem.EntityManager)

	gem.cache.DisableAutoCleanup()
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestGetByGuildID() {
	testGuildId := uint64(12345123451234512345)

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&AuditLogConfig{}).Name()),
		[]interface{}{ColumnGuildId + " = ?", testGuildId},
	).Run(func(args mock.Arguments) {
		switch v := args.Get(0).(type) {
		case *AuditLogConfig:
			v.GuildID = testGuildId
		}
	}).Return(nil).Once()

	result, err := suite.gem.GetByGuildId(testGuildId)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testGuildId, result.GuildID)

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.True(ok)
	suite.Equal(result, cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestGetByGuildIDWithCache() {
	testGuildId := uint64(12345123451234512345)

	testAuditLogConfig := &AuditLogConfig{
		GuildID: testGuildId,
	}

	cache.Update(suite.gem.cache, testGuildId, testAuditLogConfig)

	// Do not expect call of GetFirstEntity or DB calls
	// When GetFirstEntity is called, test will fail as call is unexpected

	result, err := suite.gem.GetByGuildId(testGuildId)

	suite.dba.AssertExpectations(suite.T())
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(testGuildId, result.GuildID)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestGetByGuildIDWithError() {
	testGuildId := uint64(12345123451234512345)

	expectedError := fmt.Errorf("something bad happened during database read")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On(
		"GetFirstEntity",
		mock.AnythingOfType(reflect.TypeOf(&AuditLogConfig{}).Name()),
		[]interface{}{ColumnGuildId + " = ?", testGuildId},
	).Return(expectedError).Once()

	result, err := suite.gem.GetByGuildId(testGuildId)

	suite.dba.AssertExpectations(suite.T())
	suite.Error(err)
	suite.NotNil(result)
	suite.Equal(*result, AuditLogConfig{})

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.False(ok)
	suite.Nil(cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestCreate() {
	testGuildId := uint64(12345123451234512345)
	testAuditLogConfig := AuditLogConfig{
		GuildID: testGuildId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testAuditLogConfig).Return(nil).Once()

	err := suite.gem.Create(&testAuditLogConfig)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.True(ok)
	suite.Equal(&testAuditLogConfig, cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestCreateWithError() {
	testGuildId := uint64(12345123451234512345)
	testAuditLogConfig := AuditLogConfig{
		GuildID: testGuildId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testAuditLogConfig).Return(expectedErr).Once()

	err := suite.gem.Create(&testAuditLogConfig)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.False(ok)
	suite.Nil(cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestSave() {
	testGuildId := uint64(12345123451234512345)
	testAuditLogConfig := AuditLogConfig{
		GuildID: testGuildId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testAuditLogConfig).Return(nil).Once()

	err := suite.gem.Save(&testAuditLogConfig)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.True(ok)
	suite.Equal(&testAuditLogConfig, cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestSaveWithError() {
	testGuildID := uint64(12345123451234512345)
	testAuditLogConfig := AuditLogConfig{
		GuildID: testGuildID,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testAuditLogConfig).Return(expectedErr).Once()

	err := suite.gem.Save(&testAuditLogConfig)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildID)
	suite.False(ok)
	suite.Nil(cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestUpdate() {
	testGuildId := uint64(12345123451234512345)
	testAuditLogConfig := AuditLogConfig{
		GuildID: testGuildId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testAuditLogConfig, testColumn, testValue).Return(nil).Once()

	err := suite.gem.Update(&testAuditLogConfig, testColumn, testValue)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.True(ok)
	suite.Equal(&testAuditLogConfig, cachedRegisteredComponent)
}

func (suite *AuditLogConfigEntityManagerTestSuite) TestUpdateWithError() {
	testGuildId := uint64(12345123451234512345)
	testAuditLogConfig := AuditLogConfig{
		GuildID: testGuildId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testAuditLogConfig, testColumn, testValue).Return(expectedErr).Once()

	err := suite.gem.Update(&testAuditLogConfig, testColumn, testValue)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())

	cachedRegisteredComponent, ok := cache.Get(suite.gem.cache, testGuildId)
	suite.False(ok)
	suite.Nil(cachedRegisteredComponent)
}

func TestAuditLogConfigEntityManager(t *testing.T) {
	suite.Run(t, new(AuditLogConfigEntityManagerTestSuite))
}