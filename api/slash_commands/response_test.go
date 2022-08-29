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

package slash_commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/test/discordgo_mock"
	"github.com/lazybytez/jojo-discord-bot/test/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type ResponseTestSuite struct {
	suite.Suite
}

func (suite *ResponseTestSuite) TestGenerateInteractionResponseTemplate() {
	tables := []struct {
		name        string
		description string
	}{
		{name: "Joseph Joestar", description: "Like a fine wine, I guess I just get better with age."},
		{name: "", description: "Like a fine wine, I guess I just get better with age."},
		{name: "Joseph Joestar", description: ""},
		{name: "", description: ""},
	}

	for _, table := range tables {
		result := GenerateInteractionResponseTemplate(table.name, table.description)

		suite.NotNilf(result.Embeds, "Arguments: %v, %v", table.name, table.description)
		suite.IsTypef(
			[]*discordgo.MessageEmbed{},
			result.Embeds,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Len(
			result.Embeds,
			1,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			discordgo.MessageFlags(0),
			result.Flags,
			"Arguments: %v, %v",
			table.name,
			table.description)

		embed := result.Embeds[0]

		suite.NotNilf(
			embed,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			table.name,
			embed.Title,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			table.description,
			embed.Description,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			api.DefaultEmbedColor,
			embed.Color,
			"Arguments: %v, %v",
			table.name,
			table.description)

		suite.NotNilf(
			embed.Fields,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.IsTypef(
			[]*discordgo.MessageEmbedField{},
			embed.Fields,
			"Arguments: %v, %v",
			table.name,
			table.description)
	}
}

func (suite *ResponseTestSuite) TestGenerateEphemeralInteractionResponseTemplate() {
	tables := []struct {
		name        string
		description string
	}{
		{name: "Joseph Joestar", description: "Like a fine wine, I guess I just get better with age."},
		{name: "", description: "Like a fine wine, I guess I just get better with age."},
		{name: "Joseph Joestar", description: ""},
		{name: "", description: ""},
	}

	for _, table := range tables {
		result := GenerateEphemeralInteractionResponseTemplate(table.name, table.description)

		suite.NotNilf(result.Embeds, "Arguments: %v, %v", table.name, table.description)
		suite.IsTypef(
			[]*discordgo.MessageEmbed{},
			result.Embeds,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Len(
			result.Embeds,
			1,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			discordgo.MessageFlagsEphemeral,
			result.Flags,
			"Arguments: %v, %v",
			table.name,
			table.description)

		embed := result.Embeds[0]

		suite.NotNilf(
			embed,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			table.name,
			embed.Title,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			table.description,
			embed.Description,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.Equalf(
			api.DefaultEmbedColor,
			embed.Color,
			"Arguments: %v, %v",
			table.name,
			table.description)

		suite.NotNilf(
			embed.Fields,
			"Arguments: %v, %v",
			table.name,
			table.description)
		suite.IsTypef(
			[]*discordgo.MessageEmbedField{},
			embed.Fields,
			"Arguments: %v, %v",
			table.name,
			table.description)
	}
}

func (suite *ResponseTestSuite) TestRespondWithSuccess() {
	testResponseData := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Test",
				Description: "This is a test",
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}

	session, transport := discordgo_mock.MockSession()

	interactionCreate := &discordgo.InteractionCreate{}
	interactionCreate.Interaction = &discordgo.Interaction{
		ID:    "12345123451234512345",
		Token: "4z842ghh2908ghviu2gz908vh42f90824ph2h298zrf928fdh2gi",
	}

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	requestInteractionResponse := &discordgo.InteractionResponse{}

	transport.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		&http.Response{
			StatusCode: http.StatusCreated,
		}, nil)

	Respond(component, session, interactionCreate, testResponseData)

	transport.AssertExpectations(suite.T())

	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)
	loggerMock.AssertNotCalled(suite.T(), "Err", mock.Anything, mock.Anything, mock.Anything)

	suite.Equalf(
		discordgo.InteractionResponseChannelMessageWithSource,
		requestInteractionResponse.Type,
		"Received wrong interaction response type!")
	suite.NotNil(requestInteractionResponse.Data)
	suite.NotEqual(discordgo.MessageFlagsEphemeral, requestInteractionResponse.Data.Flags)
	suite.NotNil(requestInteractionResponse.Data.Embeds)
	suite.Len(requestInteractionResponse.Data.Embeds, 1)
	suite.Equal("Test", requestInteractionResponse.Data.Embeds[0].Title)
	suite.Equal("This is a test", requestInteractionResponse.Data.Embeds[0].Description)
	suite.Equal(api.DefaultEmbedColor, requestInteractionResponse.Data.Embeds[0].Color)
}

func (suite *ResponseTestSuite) TestRespondWithSuccessAndEphemeralMessage() {
	testResponseData := &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Test",
				Description: "This is a test",
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}

	session, transport := discordgo_mock.MockSession()

	interactionCreate := &discordgo.InteractionCreate{}
	interactionCreate.Interaction = &discordgo.Interaction{
		ID:    "12345123451234512345",
		Token: "4z842ghh2908ghviu2gz908vh42f90824ph2h298zrf928fdh2gi",
	}

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	requestInteractionResponse := &discordgo.InteractionResponse{}

	transport.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		&http.Response{
			StatusCode: http.StatusCreated,
		}, nil)

	Respond(component, session, interactionCreate, testResponseData)

	transport.AssertExpectations(suite.T())

	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)
	loggerMock.AssertNotCalled(suite.T(), "Err", mock.Anything, mock.Anything, mock.Anything)

	suite.Equalf(
		discordgo.InteractionResponseChannelMessageWithSource,
		requestInteractionResponse.Type,
		"Received wrong interaction response type!")
	suite.NotNil(requestInteractionResponse.Data)
	suite.Equal(discordgo.MessageFlagsEphemeral, requestInteractionResponse.Data.Flags)
	suite.NotNil(requestInteractionResponse.Data.Embeds)
	suite.Len(requestInteractionResponse.Data.Embeds, 1)
	suite.Equal("Test", requestInteractionResponse.Data.Embeds[0].Title)
	suite.Equal("This is a test", requestInteractionResponse.Data.Embeds[0].Description)
	suite.Equal(api.DefaultEmbedColor, requestInteractionResponse.Data.Embeds[0].Color)
}

func (suite *ResponseTestSuite) TestRespondWithError() {
	testResponseData := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Test",
				Description: "This is a test",
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}

	session, transport := discordgo_mock.MockSession()

	interactionCreate := &discordgo.InteractionCreate{}
	interactionCreate.Interaction = &discordgo.Interaction{
		ID:    "12345123451234512345",
		Token: "4z842ghh2908ghviu2gz908vh42f90824ph2h298zrf928fdh2gi",
	}

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	requestInteractionResponse := &discordgo.InteractionResponse{}

	expectedErr := fmt.Errorf("something failed while processing response")

	transport.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		nil, expectedErr)

	loggerMock.On(
		"Err",
		mock.AnythingOfType(reflect.TypeOf(&url.Error{}).Name()),
		mock.Anything,
		mock.Anything).Return(nil)

	Respond(component, session, interactionCreate, testResponseData)

	transport.AssertExpectations(suite.T())

	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)
}

//func (suite *ResponseTestSuite) TestEditResponseWithSuccess() {
//	testResponseData := &discordgo.WebhookEdit{
//		Embeds: &[]*discordgo.MessageEmbed{
//			{
//				Title:       "Test",
//				Description: "This is a test",
//				Color:       api.DefaultEmbedColor,
//				Fields:      []*discordgo.MessageEmbedField{},
//			},
//		},
//	}
//
//	session, transport := discordgo_mock.MockSession()
//
//	interactionCreate := &discordgo.InteractionCreate{}
//	interactionCreate.Interaction = &discordgo.Interaction{
//		ID:    "12345123451234512345",
//		Token: "4z842ghh2908ghviu2gz908vh42f90824ph2h298zrf928fdh2gi",
//	}
//
//	component := &api.Component{}
//	loggerMock := &log.LoggerMock{}
//	component.SetLogger(loggerMock)
//
//	webHookEditData := &discordgo.WebhookEdit{}
//
//	transport.OnRequestCaptureResult(http.MethodPost, webHookEditData).Once().Return(
//		&http.Response{
//			StatusCode: http.StatusCreated,
//		}, nil)
//
//	EditResponse(component, session, interactionCreate, webHookEditData)
//
//	transport.AssertExpectations(suite.T())
//
//	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)
//	loggerMock.AssertNotCalled(suite.T(), "Err", mock.Anything, mock.Anything, mock.Anything)
//
//	suite.NotNil(testResponseData.Embeds)
//	suite.Len(testResponseData.Embeds, 1)
//	suite.Equal("Test", (*testResponseData.Embeds)[0].Title)
//	suite.Equal("This is a test", (*testResponseData.Embeds)[0].Description)
//	suite.Equal(api.DefaultEmbedColor, (*testResponseData.Embeds)[0].Color)
//}

func TestResponse(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
