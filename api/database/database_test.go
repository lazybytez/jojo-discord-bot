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

package database

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lazybytez/jojo-discord-bot/test/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestGetEntityManager(t *testing.T) {
	entityManagerDummy := GormEntityManager{
		nil,
		&log.LoggerMock{},
		entityManagers{},
	}

	entityManager = entityManagerDummy

	assert.Equal(t, entityManagerDummy, entityManager)
}

type DatabaseTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	gormDB  *gorm.DB
	logger  *log.LoggerMock
	em      EntityManager
}

func (suite *DatabaseTestSuite) SetupTest() {
	sqlMock, mock, err := sqlmock.New()
	suite.NoError(err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlMock,
		PreferSimpleProtocol: true,
	})

	dbMock, err := gorm.Open(dialector)
	suite.NoError(err)

	suite.sqlMock = mock
	suite.gormDB = dbMock

	loggerMock := &log.LoggerMock{}
	suite.logger = loggerMock

	suite.em = &GormEntityManager{
		dbMock,
		loggerMock,
		entityManagers{},
	}
}

type TestEntity struct {
	gorm.Model
	Name string
}

func (suite *DatabaseTestSuite) TestRegisterEntityWithSuccess() {
	testEntity := &TestEntity{
		Model: gorm.Model{ID: 42},
		Name:  "test entity",
	}

	// Passively check whether auto-migrate is executed or not
	// We do not have checks for when table already exists, because
	// we don't want to deeply test GORM. We just want some insights whether
	// auto-migrate is executed or not. And these statements are a good
	// indicator that AutoMigrate has been called.
	suite.sqlMock.ExpectQuery("SELECT (.+) information_schema\\.tables (.+)").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).FromCSVString("0"))
	suite.sqlMock.ExpectExec("CREATE TABLE \"test_entities\" (.+)").WillReturnResult(
		sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectExec(
		"CREATE INDEX IF NOT EXISTS \"idx_test_entities_deleted_at\" (.+)",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	suite.logger.On("Info", mock.AnythingOfType("string"), []interface{}{"TestEntity"})

	err := suite.em.RegisterEntity(testEntity)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestRegisterEntityWithFail() {
	testEntity := &TestEntity{
		Model: gorm.Model{ID: 42},
		Name:  "test entity",
	}

	expectedErr := fmt.Errorf("your next line is \"I already expected that error to happen!\"")

	suite.sqlMock.ExpectQuery("SELECT (.+) information_schema\\.tables (.+)").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).FromCSVString("0"))
	suite.sqlMock.ExpectExec("CREATE TABLE \"test_entities\" (.+)").WillReturnError(expectedErr)

	suite.logger.On("Err", expectedErr, mock.AnythingOfType("string"), []interface{}{"TestEntity"})

	err := suite.em.RegisterEntity(testEntity)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
