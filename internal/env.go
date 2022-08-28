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
	"github.com/joho/godotenv"
	"os"
	"path"
)

const envFile = ".env"

const token = "TOKEN"
const sqlMode = "DB_MODE"
const sqlDsn = "DB_DSN"

// JojoBotConfig represents the entire environment variable based configuration
// of the application. Globally available values are public,
// sensitive values are private and can only being read
// from within the current package.
//
// When new configuration options get added to the application, they should be appended
// to the structure.
type JojoBotConfig struct {
	token   string
	sqlMode string
	sqlDsn  string
}

// Config holds the currently loaded configuration
// of the application.
//
// It should be used to retrieve the set environment variables
// instead of directly calling os.Getenv.
//
// The variable is initialized during the init function.
var Config JojoBotConfig

// initEnv initializes environment with local .env file
// This will load the environment variables defined in the specified
// env file and merge them into os.Environ.
func initEnv() {
	pwd, err := os.Getwd()
	if nil != err {
		ExitFatal("Failed to get current working directory to load env file from!")
	}

	envFilePath := path.Join(pwd, envFile)
	coreLogger.Info("Trying to load env file from \"%v\"...", envFilePath)

	err = godotenv.Load(envFilePath)
	if nil == err {
		coreLogger.Info("Sucessfully loaded env file!")
	}

	Config = JojoBotConfig{
		token:   os.Getenv(token),
		sqlMode: os.Getenv(sqlMode),
		sqlDsn:  os.Getenv(sqlDsn),
	}
	coreLogger.Info("Sucessfully loaded environment configuration!")

	cleanUpSensitiveValues()
}

// cleanUpSensitiveValues throws out sensitive values from the
// environment configuration to prevent other packages from using those
// without reloading them. This reduces the attack surface from accidentally
// printing out the env to a user of the application.
//
// Note that variables loaded using godotenv.Load could be loaded again at any time
func cleanUpSensitiveValues() {
	err := os.Setenv(token, "")

	if nil != err {
		ExitFatal("Failed to cleanup sensitive environment values!")
	}
}
