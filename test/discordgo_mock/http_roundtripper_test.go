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
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type HttpRoundTripperTestSuite struct {
	suite.Suite
}

func (suite *HttpRoundTripperTestSuite) TestRoundTripWithResponse() {
	roundTripper := RoundTripper{}

	req := http.Request{
		Method: http.MethodPatch,
	}

	// We always expect the function to return what ever has been passed
	// to mock.Return
	resp := http.Response{
		StatusCode: http.StatusAccepted,
	}

	roundTripper.On("RoundTrip", &req).Return(&resp, nil)

	resultResp, resultErr := roundTripper.RoundTrip(&req)

	roundTripper.AssertExpectations(suite.T())
	suite.Nil(resultErr)
	suite.Equal(&resp, resultResp)
}

func (suite *HttpRoundTripperTestSuite) TestRoundTripWithValueResponse() {
	roundTripper := RoundTripper{}

	req := http.Request{
		Method: http.MethodPatch,
	}

	// We always expect the function to return what ever has been passed
	// to mock.Return
	resp := http.Response{
		StatusCode: http.StatusAccepted,
	}

	roundTripper.On("RoundTrip", &req).Return(resp, nil)

	resultResp, resultErr := roundTripper.RoundTrip(&req)

	roundTripper.AssertExpectations(suite.T())
	suite.Nil(resultErr)
	suite.Equal(&resp, resultResp)
}

func (suite *HttpRoundTripperTestSuite) TestRoundTripWithResponseAndError() {
	roundTripper := RoundTripper{}

	req := http.Request{
		Method: http.MethodPatch,
	}

	// We always expect the function to return what ever has been passed
	// to mock.Return
	resp := http.Response{
		StatusCode: http.StatusAccepted,
	}
	err := fmt.Errorf("this is some test error")

	roundTripper.On("RoundTrip", &req).Return(resp, err)

	resultResp, resultErr := roundTripper.RoundTrip(&req)

	roundTripper.AssertExpectations(suite.T())
	suite.EqualValues(err, resultErr)
	suite.Equal(&resp, resultResp)
}

func (suite *HttpRoundTripperTestSuite) TestRoundTripWithNoResponseAndWithError() {
	roundTripper := RoundTripper{}

	req := http.Request{
		Method: http.MethodPatch,
	}

	err := fmt.Errorf("this is some test error")

	roundTripper.On("RoundTrip", &req).Return(nil, err)

	resultResp, resultErr := roundTripper.RoundTrip(&req)

	roundTripper.AssertExpectations(suite.T())
	suite.EqualValues(err, resultErr)
	suite.Nil(resultResp)
}

func (suite *HttpRoundTripperTestSuite) TestRoundTripWithNoRequest() {
	roundTripper := RoundTripper{}

	err := fmt.Errorf("request is nil")

	roundTripper.On("RoundTrip", mock.Anything).Return(nil, err)

	resultResp, resultErr := roundTripper.RoundTrip(nil)

	roundTripper.AssertExpectations(suite.T())
	suite.EqualValues(err, resultErr)
	suite.Nil(resultResp)
}

func (suite *HttpRoundTripperTestSuite) TestNewRoundTripper() {
	roundTripper := newRoundTripper()

	suite.NotNil(roundTripper)
	suite.IsType(&RoundTripper{}, roundTripper)
}

func TestHttpRoundTripper(t *testing.T) {
	suite.Run(t, new(HttpRoundTripperTestSuite))
}
