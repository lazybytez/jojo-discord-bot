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

package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExtractTypeNameTestSuite struct {
	suite.Suite
}

func (suite *ExtractTypeNameTestSuite) TestExtractTypeName() {
	var stringVar = "other_string"
	tables := []struct {
		in  interface{}
		out string
	}{
		{discordgo.Session{}, "Session"},
		{&discordgo.Session{}, "Session"},
		{"some_string", "string"},
		{&stringVar, "string"},
	}

	for _, table := range tables {
		result := ExtractTypeName(table.in)

		suite.Equalf(table.out, result, "Arguments: %v", table.in)
	}
}

func TestReflect(t *testing.T) {
	suite.Run(t, new(ExtractTypeNameTestSuite))
}
