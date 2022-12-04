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
	"github.com/lazybytez/jojo-discord-bot/test/dbmock"
	"github.com/lazybytez/jojo-discord-bot/test/entity_manager_mock"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuditLogEntityManagerTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     entity_manager_mock.EntityManagerMock
	gem    *AuditLogEntityManager
}

func (suite *AuditLogEntityManagerTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = entity_manager_mock.EntityManagerMock{}
	suite.gem = &AuditLogEntityManager{
		&suite.em,
	}
}

func (suite *AuditLogEntityManagerTestSuite) TestNewAuditLogConfigEntityManager() {
	testEntityManager := entity_manager_mock.EntityManagerMock{}

	gem := NewAuditLogConfigEntityManager(&testEntityManager)

	suite.NotNil(gem)
	suite.NotNil(gem.cache)
	suite.Equal(&testEntityManager, gem.EntityManager)

	gem.cache.DisableAutoCleanup()
}

func (suite *AuditLogEntityManagerTestSuite) TestCreate() {
	testGuildId := uint(123123)
	testAuditLogConfig := AuditLog{
		GuildID: testGuildId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testAuditLogConfig).Return(nil).Once()

	err := suite.gem.Create(&testAuditLogConfig)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())
}

func (suite *AuditLogEntityManagerTestSuite) TestCreateWithError() {
	testGuildId := uint(123123)
	testAuditLogConfig := AuditLog{
		GuildID: testGuildId,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Create", &testAuditLogConfig).Return(expectedErr).Once()

	err := suite.gem.Create(&testAuditLogConfig)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())
}

func (suite *AuditLogEntityManagerTestSuite) TestSave() {
	testGuildId := uint(123123)
	testAuditLogConfig := AuditLog{
		GuildID: testGuildId,
	}

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testAuditLogConfig).Return(nil).Once()

	err := suite.gem.Save(&testAuditLogConfig)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())
}

func (suite *AuditLogEntityManagerTestSuite) TestSaveWithError() {
	testGuildID := uint(123123)
	testAuditLogConfig := AuditLog{
		GuildID: testGuildID,
	}

	expectedErr := fmt.Errorf("something happened during update")

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("Save", &testAuditLogConfig).Return(expectedErr).Once()

	err := suite.gem.Save(&testAuditLogConfig)

	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.dba.AssertExpectations(suite.T())
}

func (suite *AuditLogEntityManagerTestSuite) TestUpdate() {
	testGuildId := uint(123123)
	testAuditLogConfig := AuditLog{
		GuildID: testGuildId,
	}
	testColumn := "some_column"
	testValue := "some_value"

	suite.em.On("DB").Return(suite.dba)
	suite.dba.On("UpdateEntity", &testAuditLogConfig, testColumn, testValue).Return(nil).Once()

	err := suite.gem.Update(&testAuditLogConfig, testColumn, testValue)

	suite.NoError(err)
	suite.dba.AssertExpectations(suite.T())
}

func (suite *AuditLogEntityManagerTestSuite) TestUpdateWithError() {
	testGuildId := uint(123123)
	testAuditLogConfig := AuditLog{
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
}

func TestAuditLogEntityManager(t *testing.T) {
	suite.Run(t, new(AuditLogEntityManagerTestSuite))
}
