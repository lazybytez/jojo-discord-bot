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

type ResponseTemplateGeneratorTestSuite struct {
	suite.Suite
}

func (suite *ResponseTemplateGeneratorTestSuite) TestGenerateInteractionResponseTemplate() {
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

func (suite *ResponseTemplateGeneratorTestSuite) TestGenerateEphemeralInteractionResponseTemplate() {
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

func TestResponseTemplateGenerator(t *testing.T) {
	suite.Run(t, new(ResponseTemplateGeneratorTestSuite))
}

type ResponseDispatcherTestSuite struct {
	suite.Suite
	responseData      *discordgo.InteractionResponseData
	interactionCreate *discordgo.InteractionCreate
	webhookEdit       *discordgo.WebhookEdit
	message           *discordgo.Message
}

func (suite *ResponseDispatcherTestSuite) SetupTest() {
	suite.responseData = &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Test",
				Description: "This is a test",
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}

	suite.interactionCreate = &discordgo.InteractionCreate{}
	suite.interactionCreate.Interaction = &discordgo.Interaction{
		ID:    "12345123451234512345",
		Token: "4z842ghh2908ghviu2gz908vh42f90824ph2h298zrf928fdh2gi",
	}

	suite.webhookEdit = &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Test",
				Description: "This is a test",
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}

	suite.message = &discordgo.Message{
		ID: "5678998765",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Test",
				Description: "This is a test",
				Color:       api.DefaultEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}
}

func (suite *ResponseDispatcherTestSuite) TestRespondWithSuccess() {
	session, transport := discordgo_mock.MockSession()

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	requestInteractionResponse := &discordgo.InteractionResponse{}

	transport.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		&http.Response{
			StatusCode: http.StatusCreated,
		}, nil)

	Respond(component, session, suite.interactionCreate, suite.responseData)

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
	suite.Equal(suite.responseData.Embeds[0].Title, requestInteractionResponse.Data.Embeds[0].Title)
	suite.Equal(suite.responseData.Embeds[0].Description, requestInteractionResponse.Data.Embeds[0].Description)
	suite.Equal(suite.responseData.Embeds[0].Color, requestInteractionResponse.Data.Embeds[0].Color)
}

func (suite *ResponseDispatcherTestSuite) TestRespondWithSuccessAndEphemeralMessage() {
	// Add ephemeral flag
	suite.responseData.Flags = discordgo.MessageFlagsEphemeral

	session, transport := discordgo_mock.MockSession()

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	requestInteractionResponse := &discordgo.InteractionResponse{}

	transport.OnRequestCaptureResult(http.MethodPost, requestInteractionResponse).Once().Return(
		&http.Response{
			StatusCode: http.StatusCreated,
		}, nil)

	Respond(component, session, suite.interactionCreate, suite.responseData)

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
	suite.Equal(suite.responseData.Embeds[0].Title, requestInteractionResponse.Data.Embeds[0].Title)
	suite.Equal(suite.responseData.Embeds[0].Description, requestInteractionResponse.Data.Embeds[0].Description)
	suite.Equal(suite.responseData.Embeds[0].Color, requestInteractionResponse.Data.Embeds[0].Color)
}

func (suite *ResponseDispatcherTestSuite) TestRespondWithError() {
	session, transport := discordgo_mock.MockSession()

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

	Respond(component, session, suite.interactionCreate, suite.responseData)

	transport.AssertExpectations(suite.T())

	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)
}

func (suite *ResponseDispatcherTestSuite) TestEditResponseWithSuccess() {
	session, transport := discordgo_mock.MockSession()

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	webHookEditData := &discordgo.WebhookEdit{}

	call := transport.OnRequestCaptureResult(http.MethodPatch, webHookEditData).Once()
	_, err := transport.RespondWith(call, suite.message)
	suite.NoError(err)

	msg := EditResponse(component, session, suite.interactionCreate, suite.webhookEdit)

	transport.AssertExpectations(suite.T())

	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)
	loggerMock.AssertNotCalled(suite.T(), "Err", mock.Anything, mock.Anything, mock.Anything)

	suite.NotNil(webHookEditData.Embeds)
	suite.Len(*webHookEditData.Embeds, 1)
	suite.Equal((*suite.webhookEdit.Embeds)[0].Title, (*webHookEditData.Embeds)[0].Title)
	suite.Equal((*suite.webhookEdit.Embeds)[0].Description, (*webHookEditData.Embeds)[0].Description)
	suite.Equal((*suite.webhookEdit.Embeds)[0].Color, (*webHookEditData.Embeds)[0].Color)

	suite.NotNil(msg.Embeds)
	suite.Len(msg.Embeds, 1)
	suite.Equal(suite.message.Embeds[0].Title, msg.Embeds[0].Title)
	suite.Equal(suite.message.Embeds[0].Description, msg.Embeds[0].Description)
	suite.Equal(suite.message.Embeds[0].Color, msg.Embeds[0].Color)
}

func (suite *ResponseDispatcherTestSuite) TestEditResponseWithFailure() {
	session, transport := discordgo_mock.MockSession()

	component := &api.Component{}
	loggerMock := &log.LoggerMock{}
	component.SetLogger(loggerMock)

	webHookEditData := &discordgo.WebhookEdit{}

	expectedErr := fmt.Errorf("something failed while processing response")

	transport.OnRequestCaptureResult(http.MethodPatch, webHookEditData).Once().Return(nil, expectedErr)

	loggerMock.On(
		"Err",
		mock.AnythingOfType(reflect.TypeOf(&url.Error{}).Name()),
		mock.Anything,
		mock.Anything).Return(nil)

	msg := EditResponse(component, session, suite.interactionCreate, suite.webhookEdit)

	transport.AssertExpectations(suite.T())

	loggerMock.AssertNotCalled(suite.T(), "Warn", mock.Anything, mock.Anything)

	suite.Nil(msg)
}

func TestResponse(t *testing.T) {
	suite.Run(t, new(ResponseDispatcherTestSuite))
}
