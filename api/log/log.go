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
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package log

import "github.com/rs/zerolog/log"

// ComponentLogPrefix is used by the logger to prefix
// log messages with the component name.
const ComponentLogPrefix = "component"

// Debug logs a message using the applications logger.
// The log level of the logged message is debug.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Debug(componentName string, format string, v ...interface{}) {
	log.Debug().Str(ComponentLogPrefix, componentName).Msgf(format, v...)
}

// Info logs a message using the applications logger.
// The log level of the logged message is info.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Info(componentName string, format string, v ...interface{}) {
	log.Info().Str(ComponentLogPrefix, componentName).Msgf(format, v...)
}

// Warn logs a message using the applications logger.
// The log level of the logged message is warning.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Warn(componentName string, format string, v ...interface{}) {
	log.Warn().Str(ComponentLogPrefix, componentName).Msgf(format, v...)
}

// Err logs a message using the applications logger.
// The log level of the logged message is error.
// The function will append a component name to the message
// that can be used to distinguish log messages between modules.
func Err(componentName string, err error, format string, v ...interface{}) {
	log.Error().Err(err).Str(ComponentLogPrefix, componentName).Msgf(format, v...)
}
