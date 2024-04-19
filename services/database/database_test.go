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
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lazybytez/jojo-discord-bot/services"
	"github.com/lazybytez/jojo-discord-bot/test/logmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseTestSuite struct {
	suite.Suite
	sqlMock        sqlmock.Sqlmock
	gormDB         *gorm.DB
	logger         *logmock.LoggerMock
	databaseAccess services.DatabaseAccess
}

func (suite *DatabaseTestSuite) SetupTest() {
	sqlDbMock, sqlMock, err := sqlmock.New()
	suite.NoError(err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDbMock,
		PreferSimpleProtocol: true,
	})

	dbMock, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	suite.NoError(err)

	suite.sqlMock = sqlMock
	suite.gormDB = dbMock

	loggerMock := &logmock.LoggerMock{}
	suite.logger = loggerMock

	suite.databaseAccess = &GormDatabaseAccess{
		dbMock,
	}
}

type TestEntity struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (suite *DatabaseTestSuite) TestNew() {
	fakeGormDB := &gorm.DB{}

	db := New(fakeGormDB)

	suite.NotNil(db)
	suite.Equal(fakeGormDB, db.database)
}

func (suite *DatabaseTestSuite) TestRegisterEntityWithSuccess() {
	testEntity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	// Passively check whether auto-migrate is executed or not
	// We do not have checks for when table already exists, because
	// we don't want to deeply test GORM. We just want some insights whether
	// auto-migrate is executed or not. And these statements are a good
	// indicator that AutoMigrate has been called.
	suite.sqlMock.ExpectQuery("SELECT (.+) FROM information_schema\\.tables (.+)").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).FromCSVString("0"))
	suite.sqlMock.ExpectExec("CREATE TABLE \"test_entities\" (.+)").WillReturnResult(
		sqlmock.NewResult(1, 1))

	suite.logger.On("Info", mock.AnythingOfType("string"), []interface{}{"TestEntity"})

	err := suite.databaseAccess.RegisterEntity(testEntity)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestRegisterEntityWithFail() {
	testEntity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("your next line is \"I already expected that error to happen!\"")

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM information_schema\\.tables (.+)").WillReturnRows(
		sqlmock.NewRows([]string{"COUNT(*)"}).FromCSVString("0"))
	suite.sqlMock.ExpectExec("CREATE TABLE \"test_entities\" (.+)").WillReturnError(expectedErr)

	suite.logger.On("Err", expectedErr, mock.AnythingOfType("string"), []interface{}{"TestEntity"})

	err := suite.databaseAccess.RegisterEntity(testEntity)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestGetFirstEntityWithSuccess() {
	expectedResult := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM \"test_entities\" (.+) LIMIT \\D+").
		WithArgs(42, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(expectedResult.ID, expectedResult.Name))

	resultEntity := &TestEntity{}

	err := suite.databaseAccess.GetFirstEntity(resultEntity, "id = ?", 42)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())

	suite.Equal(expectedResult.ID, resultEntity.ID)
	suite.Equal(expectedResult.Name, resultEntity.Name)
}

func (suite *DatabaseTestSuite) TestGetFirstEntityWithFailure() {
	expectedResult := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("no records found")

	suite.sqlMock.ExpectQuery(
		"SELECT (.+) FROM \"test_entities\" (.+) LIMIT \\D+",
	).WithArgs(42, 1).WillReturnError(expectedErr)

	resultEntity := &TestEntity{}

	err := suite.databaseAccess.GetFirstEntity(resultEntity, "id = ?", 42)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())

	suite.NotEqual(expectedResult.ID, resultEntity.ID)
	suite.NotEqual(expectedResult.Name, resultEntity.Name)
}

func (suite *DatabaseTestSuite) TestGetLastEntityWithSuccess() {
	expectedResult := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM \"test_entities\" (.+) LIMIT \\D+").
		WithArgs(42, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(expectedResult.ID, expectedResult.Name))

	resultEntity := &TestEntity{}

	err := suite.databaseAccess.GetLastEntity(resultEntity, "id = ?", 42)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())

	suite.Equal(expectedResult.ID, resultEntity.ID)
	suite.Equal(expectedResult.Name, resultEntity.Name)
}

func (suite *DatabaseTestSuite) TestGetLastEntityWithFailure() {
	expectedResult := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("no records found")

	suite.sqlMock.ExpectQuery(
		"SELECT (.+) FROM \"test_entities\" (.+) LIMIT \\D+",
	).WithArgs(42, 1).WillReturnError(expectedErr)

	resultEntity := &TestEntity{}

	err := suite.databaseAccess.GetLastEntity(resultEntity, "id = ?", 42)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())

	suite.NotEqual(expectedResult.ID, resultEntity.ID)
	suite.NotEqual(expectedResult.Name, resultEntity.Name)
}

func (suite *DatabaseTestSuite) TestGetEntitiesWithSuccess() {
	expectedResult := &TestEntity{
		ID:   42,
		Name: "test entity",
	}
	expectedResult2 := &TestEntity{
		ID:   64,
		Name: "another test entity",
	}

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM \"test_entities\" (.+)").
		WithArgs(42, 64).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(expectedResult.ID, expectedResult.Name).
				AddRow(expectedResult2.ID, expectedResult2.Name))

	resultEntities := make([]*TestEntity, 0)

	err := suite.databaseAccess.GetEntities(&resultEntities, "id = ? OR id = ?", 42, 64)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())

	suite.Len(resultEntities, 2)

	// Checks for result entity one
	suite.Equal(expectedResult.ID, resultEntities[0].ID)
	suite.Equal(expectedResult.Name, resultEntities[0].Name)

	// Checks for result entity two
	suite.Equal(expectedResult2.ID, resultEntities[1].ID)
	suite.Equal(expectedResult2.Name, resultEntities[1].Name)
}

func (suite *DatabaseTestSuite) TestGetEntitiesWithFailure() {
	expectedErr := fmt.Errorf("no records found")

	suite.sqlMock.ExpectQuery("SELECT (.+) FROM \"test_entities\" (.+)").
		WithArgs(42, 64).
		WillReturnError(expectedErr)

	resultEntities := make([]*TestEntity, 0)

	err := suite.databaseAccess.GetEntities(&resultEntities, "id = ? OR id = ?", 42, 64)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())

	suite.Len(resultEntities, 0)
}

func (suite *DatabaseTestSuite) TestCreateWithSuccess() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery("INSERT INTO \"test_entities\" (.+)").
		WithArgs(entity.Name, entity.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(entity.ID))
	suite.sqlMock.ExpectCommit()

	err := suite.databaseAccess.Create(entity)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestCreateWithFailure() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("permission denied")

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectQuery("INSERT INTO \"test_entities\" (.+)").
		WithArgs(entity.Name, entity.ID).
		WillReturnError(expectedErr)
	suite.sqlMock.ExpectRollback()

	err := suite.databaseAccess.Create(entity)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestSaveWithSuccess() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("UPDATE \"test_entities\" SET (.+)").
		WithArgs(entity.Name, entity.ID).
		WillReturnResult(sqlmock.NewResult(42, 1))
	suite.sqlMock.ExpectCommit()

	err := suite.databaseAccess.Save(entity)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestSaveWithFailure() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("permission denied")

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("UPDATE \"test_entities\" SET (.+)").
		WithArgs(entity.Name, entity.ID).
		WillReturnError(expectedErr)
	suite.sqlMock.ExpectRollback()

	err := suite.databaseAccess.Save(entity)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestUpdateEntityWithSuccess() {
	entity := &TestEntity{
		ID:   42,
		Name: "some first test",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("UPDATE \"test_entities\" SET (.+)").
		WithArgs("some second test", entity.ID).
		WillReturnResult(sqlmock.NewResult(42, 1))
	suite.sqlMock.ExpectCommit()

	err := suite.databaseAccess.UpdateEntity(entity, "name", "some second test")

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
	suite.Equal("some second test", entity.Name)
}

func (suite *DatabaseTestSuite) TestUpdateEntityWithFailure() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("permission denied")

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("UPDATE \"test_entities\" SET (.+)").
		WithArgs("some second test", entity.ID).
		WillReturnError(expectedErr)
	suite.sqlMock.ExpectRollback()

	err := suite.databaseAccess.UpdateEntity(entity, "name", "some second test")

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestDeleteEntityWithSuccess() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("DELETE FROM \"test_entities\" WHERE (.+)").
		WithArgs(entity.ID).
		WillReturnResult(sqlmock.NewResult(42, 1))
	suite.sqlMock.ExpectCommit()

	err := suite.databaseAccess.DeleteEntity(entity)

	suite.NoError(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestDeleteEntityWithFailure() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	expectedErr := fmt.Errorf("permission denied")

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec("DELETE FROM \"test_entities\" WHERE (.+)").
		WithArgs(entity.ID).
		WillReturnError(expectedErr)
	suite.sqlMock.ExpectRollback()

	err := suite.databaseAccess.DeleteEntity(entity)

	suite.Error(err)
	suite.NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *DatabaseTestSuite) TestWorkOn() {
	entity := &TestEntity{
		ID:   42,
		Name: "test entity",
	}

	db := suite.databaseAccess.WorkOn(entity)

	suite.NotNil(db)
	suite.IsType(&gorm.DB{}, db)
	suite.NotEqual(suite.gormDB, db)
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
