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
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/test/discordgo_test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
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

	session, transport := discordgo_test.MockSession()

	interactionCreate := &discordgo.InteractionCreate{}
	interactionCreate.Interaction = &discordgo.Interaction{
		ID:    "12345123451234512345",
		Token: "4z842ghh2908ghviu2gz908vh42f90824ph2h298zrf928fdh2gi",
	}

	component := &api.Component{}

	requestInteractionResponse := &discordgo.InteractionResponse{}
	method := ""

	transport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		if nil == req {
			return false
		}

		method = req.Method

		err := json.NewDecoder(req.Body).Decode(&requestInteractionResponse)

		return nil == err
	})).Once().Return(&http.Response{
		StatusCode: http.StatusCreated,
	}, nil)

	Respond(component, session, interactionCreate, testResponseData)

	transport.AssertExpectations(suite.T())

	suite.Equalf(http.MethodPost, method, "Request was done with wrong HTTP method!")
	suite.Equalf(
		discordgo.InteractionResponseChannelMessageWithSource,
		requestInteractionResponse.Type,
		"Received wrong interaction response type!")
}

func TestResponse(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
