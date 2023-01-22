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

package bot_webapi

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/stretchr/testify/suite"
	"testing"
)

// dummyCommands holds some test commands
// that can be used to check the behaviour of the
// options to DTO converter
var optionDummyCommands = []*api.Command{
	{
		Cmd: &discordgo.ApplicationCommand{
			Name:        "jojo",
			Description: "Manage modules and core settings of the bot!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "module",
					Description: "Manage which modules should be enabled / disabled on your server!",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "list",
							Description: "List all modules and their status",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "show",
							Description: "Show information about a specific module",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "module",
									Description: "The name of the module to show information about",
									Required:    true,
									Type:        discordgo.ApplicationCommandOptionString,
									Choices: []*discordgo.ApplicationCommandOptionChoice{
										{
											Name:  "Test",
											Value: "test",
										},
										{
											Name:  "Other",
											Value: "other",
										},
										{
											Name:  "Value",
											Value: "value",
										},
									},
								},
							},
						},
						{
							Name:        "enable",
							Description: "Enable a module for the guild",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "module",
									Description: "The name of the module to enable",
									Required:    true,
									Type:        discordgo.ApplicationCommandOptionString,
									Choices: []*discordgo.ApplicationCommandOptionChoice{
										{
											Name:  "Test",
											Value: "test",
										},
										{
											Name:  "Other",
											Value: "other",
										},
										{
											Name:  "Value",
											Value: "value",
										},
									}},
							},
						},
						{
							Name:        "disable",
							Description: "Disable a module for the guild",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "module",
									Description: "The name of the module to disable",
									Required:    true,
									Type:        discordgo.ApplicationCommandOptionString,
									Choices: []*discordgo.ApplicationCommandOptionChoice{
										{
											Name:  "Test",
											Value: "test",
										},
										{
											Name:  "Other",
											Value: "other",
										},
										{
											Name:  "Value",
											Value: "value",
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "sync-commands",
					Description: "Trigger a re-synchronisation of slash-commands for the " +
						"guild to tackle inconsistencies",
					Type: discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "auditlog",
					Description: "Manage settings of the bot audit log!",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "status",
							Description: "Show the status of the current bot audit log configuration",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "enable",
							Description: "Enable printing the bot audit log to the configured channel",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "channel",
									Description: "The channel where audit log messages should be send to",
									Required:    true,
									Type:        discordgo.ApplicationCommandOptionChannel,
									ChannelTypes: []discordgo.ChannelType{
										discordgo.ChannelTypeGuildText,
									},
								},
							},
						},
						{
							Name:        "disable",
							Description: "Disable printing the bot audit log to the configured channel",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
				},
			},
		},
		Category: api.CategoryAdministration,
		Handler:  func(s *discordgo.Session, i *discordgo.InteractionCreate) {},
	},
	{
		Cmd: &discordgo.ApplicationCommand{
			Name:        "dice",
			Description: "throw one or more dice of your wished type.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "dice-sites-number",
					Description: "The number of how many sites the die has, default is 6",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
				{
					Name:        "number-dice",
					Description: "How many dice you want to throw, default is 1",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
			},
		},
		Global:   true,
		Category: api.CategoryFun,
		Handler:  func(s *discordgo.Session, i *discordgo.InteractionCreate) {},
	},
	{
		Cmd: &discordgo.ApplicationCommand{
			Name:        "stats",
			Description: "Show information of the bot and runtime statistics.",
		},
		Global:   true,
		Category: api.CategoryUtilities,
		Handler:  func(s *discordgo.Session, i *discordgo.InteractionCreate) {},
	},
	{
		Cmd: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Play ping pong with the bot!",
		},
		Category: api.CategoryFun,
		Handler:  func(s *discordgo.Session, i *discordgo.InteractionCreate) {},
	},
	{
		Cmd: &discordgo.ApplicationCommand{
			Name:        "pong",
			Description: "Play ping pong with the bot!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "test",
					Description: "Use this to test second level attributes with options",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "Some random channel you want to use",
							Required:    true,
							Type:        discordgo.ApplicationCommandOptionChannel,
							ChannelTypes: []discordgo.ChannelType{
								discordgo.ChannelTypeGuildText,
							},
						},
					},
				},
			},
		},
		Category: api.CategoryFun,
		Handler:  func(s *discordgo.Session, i *discordgo.InteractionCreate) {},
	},
}

type CommandOptionApiTypesTestSuite struct {
	suite.Suite
}

func (suite *CommandOptionApiTypesTestSuite) TestComputeCommandOptionDTOsForCommand() {
	tables := []struct {
		id       string
		expected []CommandOptionDTO
	}{
		{
			id:       "ping",
			expected: []CommandOptionDTO{},
		},
		{
			id:       "jojo_sync-commands",
			expected: []CommandOptionDTO{},
		},
		{
			id:       "jojo_module_list",
			expected: []CommandOptionDTO{},
		},
		{
			id: "dice",
			expected: []CommandOptionDTO{
				{
					Owner:   "dice",
					Name:    "dice-sites-number",
					Type:    int(discordgo.ApplicationCommandOptionInteger),
					Choices: nil,
				},
				{
					Owner:   "dice",
					Name:    "number-dice",
					Type:    int(discordgo.ApplicationCommandOptionInteger),
					Choices: nil,
				},
			},
		},
		{
			id: "jojo_module_enable",
			expected: []CommandOptionDTO{
				{
					Owner: "jojo_module_enable",
					Name:  "module",
					Type:  int(discordgo.ApplicationCommandOptionString),
					Choices: []CommandOptionChoiceDTO{
						{
							Name:  "Test",
							Value: "test",
						},
						{
							Name:  "Other",
							Value: "other",
						},
						{
							Name:  "Value",
							Value: "value",
						},
					},
				},
			},
		},
		{
			id: "pong_test",
			expected: []CommandOptionDTO{
				{
					Owner:   "pong_test",
					Name:    "channel",
					Type:    int(discordgo.ApplicationCommandOptionChannel),
					Choices: nil,
				},
			}},
	}

	for _, table := range tables {
		result, err := computeCommandOptionDTOsForCommand(optionDummyCommands, table.id)

		suite.NoErrorf(
			err,
			"did not expect error for command id \"%s\"",
			table.id)
		suite.EqualValuesf(
			table.expected,
			result,
			"got wrong result for test with command id \"%s\"",
			table.id)
	}
}

func (suite *CommandOptionApiTypesTestSuite) TestComputeCommandOptionDTOsForCommandWithMissingCommands() {
	tables := []struct {
		id string
	}{
		{
			id: "somenotexistingcommand",
		},
		{
			id: "somenot_existingcommand",
		},
		{
			id: "somenot_existing_command",
		},
	}

	for _, table := range tables {
		result, err := computeCommandOptionDTOsForCommand(optionDummyCommands, table.id)

		suite.Errorf(
			err,
			"did not expect error for command id \"%s\"",
			table.id)
		suite.EqualValuesf(
			[]CommandOptionDTO{},
			result,
			"got wrong result for test with command id \"%s\"",
			table.id)
	}
}

func TestCommandOptionApiTypes(t *testing.T) {
	suite.Run(t, new(CommandOptionApiTypesTestSuite))
}
