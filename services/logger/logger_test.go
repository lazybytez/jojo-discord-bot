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
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rzajac/zltest"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type LogTestSuite struct {
	suite.Suite
}

func (suite *LogTestSuite) TestNew() {
	customLogger := zerolog.Logger{}

	tables := []struct {
		prefix         string
		logger         *zerolog.Logger
		expectedPrefix string
		expectedType   interface{}
	}{
		{"test", nil, "test", &zerolog.Logger{}},
		{"other test", &customLogger, "other test", &zerolog.Logger{}},
	}

	for _, table := range tables {
		result := New(table.prefix, table.logger)

		loggerType := "nil"
		// Logic in tests is not good, but in this case it is only for failed test output
		if nil != table.logger {
			loggerType = reflect.TypeOf(table.logger).Name()
		}

		suite.NotNilf(
			result,
			"Returned logger reference should never be nil! Arguments: %v, %v",
			table.prefix,
			loggerType)
		suite.Equalf(
			table.expectedPrefix,
			result.prefix,
			"Arguments: %v, %v",
			table.prefix,
			loggerType)
		suite.IsTypef(
			table.expectedType,
			result.loggerImpl,
			"Arguments: %v, %v",
			table.prefix,
			loggerType)
	}
}

func (suite *LogTestSuite) TestDebug() {
	tables := []struct {
		prefix string
		msg    string
		values []interface{}
	}{
		{"test_prefix", "test", make([]interface{}, 0)},
		{
			"another_test_prefix",
			"test with one string value: \"%s\"",
			append(make([]interface{}, 0), "Some Value!"),
		},
		{
			"second_another_test_prefix",
			"test with two string values: \"%s\", \"%s\"",
			append(make([]interface{}, 0), "Some Value!", "Second Value..."),
		},
		{
			"another_test_prefix",
			"test with one int value: \"%d\"",
			append(make([]interface{}, 0), 42),
		},
		{
			"second_another_test_prefix",
			"test with two int values: \"%d\", \"%d\"",
			append(make([]interface{}, 0), 32, 64),
		},
		{
			"second_another_test_prefix",
			"test with two mixed values: \"%s\", \"%d\"",
			append(make([]interface{}, 0), "One stack of items:", 64),
		},
		{
			"second_another_test_prefix",
			"test with values: \"%v\", \"%v\"",
			append(make([]interface{}, 0), map[string]int{"some": 5}, true),
		},
	}

	for _, table := range tables {
		tst := zltest.New(suite.T())

		zeroLogger := zerolog.New(tst).With().Timestamp().Logger()

		customLogger := Logger{
			prefix:     table.prefix,
			loggerImpl: &zeroLogger,
		}

		customLogger.Debug(table.msg, table.values...)

		logEntry := tst.LastEntry()
		logEntry.ExpLevel(zerolog.DebugLevel)
		logEntry.ExpStr(componentLogPrefix, table.prefix)
		logEntry.ExpMsg(fmt.Sprintf(table.msg, table.values...))
	}
}

func (suite *LogTestSuite) TestInfo() {
	tables := []struct {
		prefix string
		msg    string
		values []interface{}
	}{
		{"test_prefix", "test", make([]interface{}, 0)},
		{
			"another_test_prefix",
			"test with one string value: \"%s\"",
			append(make([]interface{}, 0), "Some Value!"),
		},
		{
			"second_another_test_prefix",
			"test with two string values: \"%s\", \"%s\"",
			append(make([]interface{}, 0), "Some Value!", "Second Value..."),
		},
		{
			"another_test_prefix",
			"test with one int value: \"%d\"",
			append(make([]interface{}, 0), 42),
		},
		{
			"second_another_test_prefix",
			"test with two int values: \"%d\", \"%d\"",
			append(make([]interface{}, 0), 32, 64),
		},
		{
			"second_another_test_prefix",
			"test with two mixed values: \"%s\", \"%d\"",
			append(make([]interface{}, 0), "One stack of items:", 64),
		},
		{
			"second_another_test_prefix",
			"test with values: \"%v\", \"%v\"",
			append(make([]interface{}, 0), map[string]int{"some": 5}, true),
		},
	}

	for _, table := range tables {
		tst := zltest.New(suite.T())

		zeroLogger := zerolog.New(tst).With().Timestamp().Logger()

		customLogger := Logger{
			prefix:     table.prefix,
			loggerImpl: &zeroLogger,
		}

		customLogger.Info(table.msg, table.values...)

		logEntry := tst.LastEntry()
		logEntry.ExpLevel(zerolog.InfoLevel)
		logEntry.ExpStr(componentLogPrefix, table.prefix)
		logEntry.ExpMsg(fmt.Sprintf(table.msg, table.values...))
	}
}

func (suite *LogTestSuite) TestWarn() {
	tables := []struct {
		prefix string
		msg    string
		values []interface{}
	}{
		{"test_prefix", "test", make([]interface{}, 0)},
		{
			"another_test_prefix",
			"test with one string value: \"%s\"",
			append(make([]interface{}, 0), "Some Value!"),
		},
		{
			"second_another_test_prefix",
			"test with two string values: \"%s\", \"%s\"",
			append(make([]interface{}, 0), "Some Value!", "Second Value..."),
		},
		{
			"another_test_prefix",
			"test with one int value: \"%d\"",
			append(make([]interface{}, 0), 42),
		},
		{
			"second_another_test_prefix",
			"test with two int values: \"%d\", \"%d\"",
			append(make([]interface{}, 0), 32, 64),
		},
		{
			"second_another_test_prefix",
			"test with two mixed values: \"%s\", \"%d\"",
			append(make([]interface{}, 0), "One stack of items:", 64),
		},
		{
			"second_another_test_prefix",
			"test with values: \"%v\", \"%v\"",
			append(make([]interface{}, 0), map[string]int{"some": 5}, true),
		},
	}

	for _, table := range tables {
		tst := zltest.New(suite.T())

		zeroLogger := zerolog.New(tst).With().Timestamp().Logger()

		customLogger := Logger{
			prefix:     table.prefix,
			loggerImpl: &zeroLogger,
		}

		customLogger.Warn(table.msg, table.values...)

		logEntry := tst.LastEntry()
		logEntry.ExpLevel(zerolog.WarnLevel)
		logEntry.ExpStr(componentLogPrefix, table.prefix)
		logEntry.ExpMsg(fmt.Sprintf(table.msg, table.values...))
	}
}

func (suite *LogTestSuite) TestErr() {
	tables := []struct {
		prefix string
		err    error
		msg    string
		values []interface{}
	}{
		{"test_prefix", fmt.Errorf("an error happened"), "test", make([]interface{}, 0)},
		{
			"another_test_prefix",
			fmt.Errorf("an error happened in XY1"),
			"test with one string value: \"%s\"",
			append(make([]interface{}, 0), "Some Value!"),
		},
		{
			"second_another_test_prefix",
			fmt.Errorf("an error happened in XY2"),
			"test with two string values: \"%s\", \"%s\"",
			append(make([]interface{}, 0), "Some Value!", "Second Value..."),
		},
		{
			"another_test_prefix",
			fmt.Errorf("an error happened in XY3"),
			"test with one int value: \"%d\"",
			append(make([]interface{}, 0), 42),
		},
		{
			"second_another_test_prefix",
			fmt.Errorf("an error happened in XY4"),
			"test with two int values: \"%d\", \"%d\"",
			append(make([]interface{}, 0), 32, 64),
		},
		{
			"second_another_test_prefix",
			fmt.Errorf("an error happened in XY5"),
			"test with two mixed values: \"%s\", \"%d\"",
			append(make([]interface{}, 0), "One stack of items:", 64),
		},
		{
			"second_another_test_prefix",
			fmt.Errorf("an error happened in XY6"),
			"test with values: \"%v\", \"%v\"",
			append(make([]interface{}, 0), map[string]int{"some": 5}, true),
		},
	}

	for _, table := range tables {
		tst := zltest.New(suite.T())

		zeroLogger := zerolog.New(tst).With().Timestamp().Logger()

		customLogger := Logger{
			prefix:     table.prefix,
			loggerImpl: &zeroLogger,
		}

		customLogger.Err(table.err, table.msg, table.values...)

		logEntry := tst.LastEntry()
		logEntry.ExpLevel(zerolog.ErrorLevel)
		logEntry.ExpStr(componentLogPrefix, table.prefix)
		logEntry.ExpMsg(fmt.Sprintf(table.msg, table.values...))
	}
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}
