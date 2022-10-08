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
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/test/dbmock"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EntityManagersTestSuite struct {
	suite.Suite
	dba    *dbmock.DatabaseAccessMock
	logger *logmock.LoggerMock
	em     EntityManager
}

func (suite *EntityManagersTestSuite) SetupTest() {
	dba := &dbmock.DatabaseAccessMock{}
	logger := &logmock.LoggerMock{}

	suite.dba = dba
	suite.logger = logger
	suite.em = NewEntityManager(suite.dba, suite.logger)
}

func (suite *EntityManagersTestSuite) TestGetGuildEntityManagerWithExistingGuildEntityManager() {
	guildEntityManager := &entities.GuildEntityManager{}

	suite.em.guild = guildEntityManager

	result := suite.em.Guilds()

	suite.NotNil(result)
	suite.Equal(guildEntityManager, result)
}

func (suite *EntityManagersTestSuite) TestGetGuildEntityManagerWithNoExistingGuildEntityManager() {
	result := suite.em.Guilds()
	result2 := suite.em.Guilds()

	// First call
	suite.NotNil(result)
	suite.IsType(&entities.GuildEntityManager{}, result)

	// Consecutive calls
	suite.Equal(result, result2)
}

func (suite *EntityManagersTestSuite) TestGetRegisteredComponentEntityManagerWithExistingRegisteredComponentEntityManager() {
	registeredComponentEntityManager := &entities.RegisteredComponentEntityManager{}

	suite.em.registeredComponentEntityManager = registeredComponentEntityManager

	result := suite.em.RegisteredComponent()

	suite.NotNil(result)
	suite.Equal(registeredComponentEntityManager, result)
}

func (suite *EntityManagersTestSuite) TestGetRegisteredComponentEntityManagerWithNoExistingRegisteredComponentEntityManager() {
	result := suite.em.RegisteredComponent()
	result2 := suite.em.RegisteredComponent()

	// First call
	suite.NotNil(result)
	suite.IsType(&entities.RegisteredComponentEntityManager{}, result)

	// Consecutive calls
	suite.Equal(result, result2)
}

func (suite *EntityManagersTestSuite) TestGetGlobalComponentStatusEntityManagerWithExistingGlobalComponentStatusEntityManager() {
	globalComponentStatusEntityManager := &entities.GlobalComponentStatusEntityManager{}

	suite.em.globalComponentStatusEntityManager = globalComponentStatusEntityManager

	result := suite.em.GlobalComponentStatus()

	suite.NotNil(result)
	suite.Equal(globalComponentStatusEntityManager, result)
}

func (suite *EntityManagersTestSuite) TestGetGlobalComponentStatusEntityManagerWithNoExistingGlobalComponentStatusEntityManager() {
	result := suite.em.GlobalComponentStatus()
	result2 := suite.em.GlobalComponentStatus()

	// First call
	suite.NotNil(result)
	suite.IsType(&entities.GlobalComponentStatusEntityManager{}, result)

	// Consecutive calls
	suite.Equal(result, result2)
}

func (suite *EntityManagersTestSuite) TestGetGuildComponentStatusEntityManagerWithExistingGuildComponentStatusEntityManager() {
	guildComponentStatusEntityManager := &entities.GuildComponentStatusEntityManager{}

	suite.em.guildComponentStatusEntityManager = guildComponentStatusEntityManager

	result := suite.em.GuildComponentStatus()

	suite.NotNil(result)
	suite.Equal(guildComponentStatusEntityManager, result)
}

func (suite *EntityManagersTestSuite) TestGetGuildComponentStatusEntityManagerWithNoExistingGuildComponentStatusEntityManager() {
	result := suite.em.GuildComponentStatus()
	result2 := suite.em.GuildComponentStatus()

	// First call
	suite.NotNil(result)
	suite.IsType(&entities.GuildComponentStatusEntityManager{}, result)

	// Consecutive calls
	suite.Equal(result, result2)
}

func TestEntityManagers(t *testing.T) {
	suite.Run(t, new(EntityManagersTestSuite))
}
