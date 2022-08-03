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

import "github.com/lazybytez/jojo-discord-bot/api/log"

// ComponentLogger is a type that is used to hold
// the owner that keeps the information about the
// component used by the logging methods.
type ComponentLogger struct {
	owner *Component
}

// Logger provides useful methods that ease logging.
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Err(err error, format string, v ...interface{})
}

// Logger is used to obtain the Logger of a component
//
// On first call, this function initializes the private Component.logger
// field. On consecutive calls, the already present Logger will be used.
func (c *Component) Logger() Logger {
	if nil == c.logger {
		c.logger = &ComponentLogger{owner: c}
	}

	return c.logger
}

// Debug logs a message with level debug.
// This function appends the name of the Component from the receiver
// to the log message.
func (c *ComponentLogger) Debug(format string, v ...interface{}) {
	log.Debug(c.owner.Name, format, v...)
}

// Info logs a message with level info.
// This function appends the name of the Component from the receiver
// to the log message.
func (c *ComponentLogger) Info(format string, v ...interface{}) {
	log.Info(c.owner.Name, format, v...)
}

// Warn logs a message with level warnings.
// This function appends the name of the Component from the receiver
// to the log message.
func (c *ComponentLogger) Warn(format string, v ...interface{}) {
	log.Warn(c.owner.Name, format, v...)
}

// Err logs a message with level error.
// This function appends the name of the Component from the receiver
// to the log message.
//
// The supplied error will be applied to the log message.
func (c *ComponentLogger) Err(err error, format string, v ...interface{}) {
	log.Err(c.owner.Name, err, format, v...)
}
