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
	"github.com/stretchr/testify/mock"
	"net/http"
)

// RoundTripper is a custom http.RoundTripper set for the target discordgo.Session.
// It is special in the manner that it is a mock and its http.RoundTrip function
// acts as a mocked method.
// This allows in tests to completely intercept the way discordgo handles requests.
type RoundTripper struct {
	mock.Mock
}

// RoundTrip handles the request by letting testifies mock implementation
// handle the processing.
func (roundTripper *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := roundTripper.Mock.Called(req)

	resp := args.Get(0)
	err := args.Error(1)

	switch resp := resp.(type) {
	case *http.Response:
		return resp, err
	case http.Response:
		return &resp, err
	default:
		return nil, err
	}
}

// newRoundTripper creates a new RoundTripper to pass to discordgo.Session
func newRoundTripper() *RoundTripper {
	return &RoundTripper{}
}
