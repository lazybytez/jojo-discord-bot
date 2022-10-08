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

package internal

import (
	"fmt"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/services/database"
	serviceLogger "github.com/lazybytez/jojo-discord-bot/services/logger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ModeSQLite = "SQLite"
const ModePostgres = "Postgres"

const DatabaseLoggerPrefix = "database"

var gormDB *gorm.DB

// initGorm initializes GORM to allow the use of databases in the application
func initGorm() {
	var dial *gorm.Dialector

	coreLogger.Info("Setting up database connection...")
	switch Config.sqlMode {
	case ModeSQLite:
		coreLogger.Info("Using SQLite as database driver!")
		dial = getSQLiteDialector()
	case ModePostgres:
		coreLogger.Info("Using PostgreSQL as database driver!")
		dial = getPostgresDialector()
	default:
		ExitFatal(fmt.Sprintf("The database mode \"%v\" is not valid!", Config.sqlMode))
	}

	coreLogger.Info("Open GORM instance...")
	var err error
	gormDB, err = gorm.Open(*dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if nil != err {
		ExitFatal(fmt.Sprintf("Failed to initialize database subsystem! Error: \"%v\"", err.Error()))
	}

	coreLogger.Info("Database subsystem has been initialized successfully!")
}

// CreateEntityManager creates a new api.EntityManager with the default gorm.DB instance
func CreateEntityManager() api.EntityManager {
	em := api.NewEntityManager(database.New(gormDB), serviceLogger.New(DatabaseLoggerPrefix, nil))

	return em
}

// getSQLiteDialector creates a SQLite gormDB and returns the gorm.Dialector instance
//
// The gormDB is created either:
//  1. when DB_DSN is present under the specified path
//  2. when DB_DSN is empty under the current working directory
func getSQLiteDialector() *gorm.Dialector {
	dsn := Config.sqlDsn
	if "" == dsn {
		dsn = "jojo.db"
	}
	sqlDialector := sqlite.Open(dsn)

	return &sqlDialector
}

// getPostgresDialector creates a gormDB connection and returns the gorm.Dialector instance
//
// The defined DB_DSN will be used to establish the connection.
// If the DSN is empty, the application will crash!
func getPostgresDialector() *gorm.Dialector {
	dsn := Config.sqlDsn
	if "" == dsn {
		ExitFatal(fmt.Sprintf("A DB_DSN must be defined when using DB_MODE \"%v\"", ModePostgres))
	}

	sqlDialector := postgres.Open(dsn)

	return &sqlDialector
}
