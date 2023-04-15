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

package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SlashCommandManagerTestSuite struct {
	suite.Suite
	owningComponent     *Component
	slashCommandManager SlashCommandManager
}

func (suite *SlashCommandManagerTestSuite) SetupTest() {
	suite.owningComponent = &Component{
		Code: "test_component",
		Name: "Test Component",
	}
	suite.slashCommandManager = SlashCommandManager{owner: suite.owningComponent}
}

func (suite *SlashCommandManagerTestSuite) TestComputeFullCommandStringFromInteractionData() {
	tables := []struct {
		input    discordgo.ApplicationCommandInteractionData
		expected string
	}{
		{
			input: discordgo.ApplicationCommandInteractionData{
				ID:   "123451234512345",
				Name: "dice",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type:  discordgo.ApplicationCommandOptionInteger,
						Name:  "dice-side-number",
						Value: 5,
					},
				},
			},
			expected: "dice",
		},
		{
			input: discordgo.ApplicationCommandInteractionData{
				ID:   "123451234512345",
				Name: "jojo",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type: discordgo.ApplicationCommandOptionSubCommand,
						Name: "sync-commands",
						Value: &discordgo.ApplicationCommandInteractionDataOption{
							Name: "sync-commands",
						},
					},
				},
			},
			expected: "jojo sync-commands",
		},
		{
			input: discordgo.ApplicationCommandInteractionData{
				ID:   "123451234512345",
				Name: "jojo",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type: discordgo.ApplicationCommandOptionSubCommand,
						Name: "sync-commands",
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Type:  discordgo.ApplicationCommandOptionString,
								Name:  "some-random-option",
								Value: "test",
							},
							{
								Type:  discordgo.ApplicationCommandOptionInteger,
								Name:  "some-second-option",
								Value: 2,
							},
						},
					},
				},
			},
			expected: "jojo sync-commands",
		},
		{
			input: discordgo.ApplicationCommandInteractionData{
				ID:   "123451234512345",
				Name: "jojo",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type: discordgo.ApplicationCommandOptionSubCommandGroup,
						Name: "module",
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Type: discordgo.ApplicationCommandOptionSubCommand,
								Name: "list",
							},
						},
					},
				},
			},
			expected: "jojo module list",
		},
		{
			input: discordgo.ApplicationCommandInteractionData{
				ID:   "123451234512345",
				Name: "jojo",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Type: discordgo.ApplicationCommandOptionSubCommandGroup,
						Name: "module",
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Type: discordgo.ApplicationCommandOptionSubCommand,
								Name: "enable",
								Options: []*discordgo.ApplicationCommandInteractionDataOption{
									{
										Type:  discordgo.ApplicationCommandOptionString,
										Name:  "module",
										Value: "ping_pong",
									},
								},
							},
						},
					},
				},
			},
			expected: "jojo module enable",
		},
	}

	for _, table := range tables {
		result := suite.slashCommandManager.computeFullCommandStringFromInteractionData(table.input)

		suite.Equal(table.expected, result)
	}
}

func (suite *SlashCommandManagerTestSuite) TestComputeConfiguredOptionsString() {
	tables := []struct {
		input    []*discordgo.ApplicationCommandInteractionDataOption
		expected string
	}{
		{
			input: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Type:  discordgo.ApplicationCommandOptionInteger,
					Name:  "dice-side-number",
					Value: 5,
				},
			},
			expected: "dice-side-number=5",
		},
		{
			input: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Name: "sync-commands",
					Value: &discordgo.ApplicationCommandInteractionDataOption{
						Name: "sync-commands",
					},
				},
			},
			expected: "",
		},
		{
			input: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Name: "sync-commands",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "some-random-option",
							Value: "test",
						},
						{
							Type:  discordgo.ApplicationCommandOptionInteger,
							Name:  "some-second-option",
							Value: 2,
						},
					},
				},
			},
			expected: "some-random-option=test; some-second-option=2",
		},
		{
			input: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
					Name: "module",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Type: discordgo.ApplicationCommandOptionSubCommand,
							Name: "list",
						},
					},
				},
			},
			expected: "",
		},
		{
			input: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
					Name: "module",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Type: discordgo.ApplicationCommandOptionSubCommand,
							Name: "enable",
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{
									Type:  discordgo.ApplicationCommandOptionString,
									Name:  "module",
									Value: "ping_pong",
								},
							},
						},
					},
				},
			},
			expected: "module=ping_pong",
		},
	}

	for _, table := range tables {
		result := suite.slashCommandManager.computeConfiguredOptionsString(table.input)

		suite.Equal(table.expected, result)
	}
}

func (suite *SlashCommandManagerTestSuite) TestGetCommandsForComponentWithoutMatchingCommands() {
	testComponentCode := entities.ComponentCode("no_commands_component")
	componentCommandMap = map[string]*Command{
		"a": {
			Category: CategoryAdministration,
			c:        suite.owningComponent,
		},
		"b": {
			Category: CategoryUtilities,
			c:        suite.owningComponent,
		},
		"c": {
			Category: CategoryAdministration,
			c:        suite.owningComponent,
		},
	}

	result := suite.slashCommandManager.GetCommandsForComponent(testComponentCode)

	suite.Equal([]*Command{}, result)
}

func (suite *SlashCommandManagerTestSuite) TestGetCommandsForComponentWithCommands() {
	testComponentCode := entities.ComponentCode("with_commands_component")

	testComponent := &Component{
		Code: testComponentCode,
	}

	foundCommandOne := &Command{
		Category: CategoryAdministration,
		c:        testComponent,
	}
	foundCommandTwo := &Command{
		Category: CategoryUtilities,
		c:        testComponent,
	}

	componentCommandMap = map[string]*Command{
		"a": foundCommandOne,
		"c": {
			Category: CategoryAdministration,
			c:        suite.owningComponent,
		},
		"b": foundCommandTwo,
	}

	result := suite.slashCommandManager.GetCommandsForComponent(testComponentCode)
	expected := []*Command{
		foundCommandOne,
		foundCommandTwo,
	}

	suite.Len(result, 2)
	suite.Equal(expected, result)
}

func TestSlashCommandManager(t *testing.T) {
	suite.Run(t, new(SlashCommandManagerTestSuite))
}
