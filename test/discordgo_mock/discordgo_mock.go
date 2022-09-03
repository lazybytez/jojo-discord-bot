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
	"net/http"
)

// This provides a custom http.RoundTripper that
// embeds a mock. When creating a new discordgo.Session
// using functions of this package, it is possible to expect
// specific requests and return specific responses.
// This can be used to create integration and some-kind of unit tests.
// Note that unit tests using this helper tend to be a discordgo integration test,
// as the entire discordgo library is run through, except that a fake http.RounTripper is used.

// MockSession returns a new discordgo.Session that is modified
// with custom transport that can be used with testify.
func MockSession() (*discordgo.Session, *RoundTripper) {
	mockRoundTripper := newRoundTripper()

	session := &discordgo.Session{
		SyncEvents:  true,
		Ratelimiter: discordgo.NewRatelimiter(),
		Client: &http.Client{
			Transport: mockRoundTripper,
		},
	}

	return session, mockRoundTripper
}
