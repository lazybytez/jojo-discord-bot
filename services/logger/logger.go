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

package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ComponentLogPrefix is used by the logger to prefix
// log messages with the component name.
const ComponentLogPrefix = "component"

// Logger is a type that is used to hold metadata
// that should be used in the loggers methods.
type Logger struct {
	loggerImpl *zerolog.Logger
	prefix     string
}

// New creates a new logger with the passed prefix
// and logger implementation.
// When loggerImpl is nil, the default logger of zerolog will be used.
func New(prefix string, loggerImpl *zerolog.Logger) *Logger {
	if nil == loggerImpl {
		loggerImpl = &log.Logger
	}

	return &Logger{
		prefix:     prefix,
		loggerImpl: loggerImpl,
	}
}

// Debug logs a message with level debug.
// This function appends the name of the Component from the receiver
// to the log message.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.loggerImpl.Debug().Str(ComponentLogPrefix, l.prefix).Msgf(format, v...)
}

// Info logs a message with level info.
// This function appends the name of the Component from the receiver
// to the log message.
func (l *Logger) Info(format string, v ...interface{}) {
	l.loggerImpl.Info().Str(ComponentLogPrefix, l.prefix).Msgf(format, v...)
}

// Warn logs a message with level warnings.
// This function appends the name of the Component from the receiver
// to the log message.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.loggerImpl.Warn().Str(ComponentLogPrefix, l.prefix).Msgf(format, v...)
}

// Err logs a message with level error.
// This function appends the name of the Component from the receiver
// to the log message.
//
// The supplied error will be applied to the log message.
func (l *Logger) Err(err error, format string, v ...interface{}) {
	l.loggerImpl.Error().Err(err).Str(ComponentLogPrefix, l.prefix).Msgf(format, v...)
}

// Logger returns the real zerolog.Logger used by this logger
// implementation. Except for internal usage (e.g. internal or api packages),
// the wrapper functions should be used.
func (l *Logger) Logger() *zerolog.Logger {
	return l.loggerImpl
}
