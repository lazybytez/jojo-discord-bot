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
	"github.com/joho/godotenv"
	"os"
	"path"
)

const envFile = ".env"

// Available environment variables
const (
	token      = "TOKEN"
	sqlMode    = "DB_MODE"
	sqlDsn     = "DB_DSN"
	webApiMode = "WEBAPI_MODE"
	webApiHost = "WEBAPI_HOST"
)

// JojoBotConfig represents the entire environment variable based configuration
// of the application. Globally available values are public,
// sensitive values are private and can only being read
// from within the current package.
//
// When new configuration options get added to the application, they should be appended
// to the structure.
type JojoBotConfig struct {
	token      string
	sqlMode    string
	sqlDsn     string
	webApiMode string
	webApiHost string
}

// Config holds the currently loaded configuration
// of the application.
//
// It should be used to retrieve the set environment variables
// instead of directly calling os.Getenv.
//
// The variable is initialized during the init function.
var Config JojoBotConfig

// getEnvOrDefault tries to get an environment variable from the environment.
// If value could not be found, the passed default value will be used.
func getEnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if "" == val {
		return defaultValue
	}

	return val
}

// getEnvOrDefault tries to get an environment variable from the environment.
// If value could not be found, the application will exit with a fatal crash.
func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if "" == val {
		ExitFatal(fmt.Sprintf("The required environment variable \"%s\" has not been configured!", key))
	}

	return val
}

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
		coreLogger.Info("Successfully loaded env file!")
	}

	Config = JojoBotConfig{
		token:      getEnvOrFail(token),
		sqlMode:    getEnvOrFail(sqlMode),
		sqlDsn:     getEnvOrDefault(sqlDsn, ""),
		webApiMode: getEnvOrDefault(webApiMode, DefaultWebApiMode),
		webApiHost: getEnvOrDefault(webApiHost, DefaultWebApiHost),
	}
	coreLogger.Info("Successfully loaded environment configuration!")
}
