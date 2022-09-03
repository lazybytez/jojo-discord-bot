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

package discordgo_mock

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/suite"
	"net/http"
	"reflect"
	"testing"
)

type DiscordGoMockTestSuite struct {
	suite.Suite
}

func (suite *DiscordGoMockTestSuite) TestMockSession() {
	session, transport := MockSession()

	suite.NotNil(session, "Did not expect to receive a nil session!")
	suite.True(session.SyncEvents, "Expected dummy session to use sync events!")

	suite.NotNil(session.Ratelimiter, "Did not expect session.RateLimiter to be nil!")
	suite.IsTypef(
		&discordgo.RateLimiter{},
		session.Ratelimiter,
		"Expected RateLimiter on session of type \"%v\"",
		reflect.TypeOf(discordgo.RateLimiter{}).Name())

	suite.NotNil(session.Client, "Did not expect session.Client to be nil!")
	suite.IsTypef(&http.Client{}, session.Client, "Expected session.Client to be of type *http.Client!")

	suite.NotNil(session.Client.Transport, "Expected session.Client.Transport not to be nil!")
	suite.IsType(
		&RoundTripper{}, session.Client.Transport,
		"Expected session.Client.Transport to be of type \"%v\"!",
		reflect.TypeOf(&RoundTripper{}).Name())
	suite.Equal(
		session.Client.Transport,
		transport,
		"Expected transport reference of client to equal returned transport reference of MockSession()")
}

func TestDiscordGoMock(t *testing.T) {
	suite.Run(t, new(DiscordGoMockTestSuite))
}
