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

package log

import "github.com/stretchr/testify/mock"

// LoggerMock is a custom logger embedding
// mock.Mock and allows to do expectations on logging methods.
type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Debug(format string, v ...interface{}) {
	l.Called(format, v)
}

func (l *LoggerMock) Info(format string, v ...interface{}) {
	l.Called(format, v)
}

func (l *LoggerMock) Warn(format string, v ...interface{}) {
	l.Called(format, v)
}

func (l *LoggerMock) Err(err error, format string, v ...interface{}) {
	l.Called(err, format, v)
}
