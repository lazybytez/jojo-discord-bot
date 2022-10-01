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

package bot_core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/core_components/bot_core/command/module"
	"github.com/lazybytez/jojo-discord-bot/core_components/bot_core/command/sync_commands"
)

// jojoCommand holds the command configuration for the jojo command.
var jojoCommand *api.Command

// getModuleCommandChoices builds a slice containing all available modules
// as command option choices
func getModuleCommandChoices() []*discordgo.ApplicationCommandOptionChoice {
	availableModuleChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)

	for _, comp := range C.EntityManager().RegisteredComponent().GetAvailable() {
		if comp.IsCoreComponent() {
			continue
		}

		availableModuleChoices = append(availableModuleChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  comp.Name,
			Value: comp.Code,
		})
	}

	return availableModuleChoices
}

// initAndRegisterJojoCommand initializes the jojo command variable and registers the command
// in the command API
func initAndRegisterJojoCommand() {
	// Ensure the module package knows about the component
	module.C = &C
	sync_commands.C = &C

	jojoCommand = &api.Command{
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
									Choices:     getModuleCommandChoices(),
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
									Choices:     getModuleCommandChoices(),
								},
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
									Choices:     getModuleCommandChoices(),
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
			},
		},
		Handler: handleJojoCommand,
	}

	_ = C.SlashCommandManager().Register(jojoCommand)
}

// handleJojoCommand handles the parent JOJO command and delegates sub-command
// handling to the appropriate handlers
func handleJojoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	subCommands := map[string]func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
		option *discordgo.ApplicationCommandInteractionDataOption,
	){
		"module":        module.HandleModuleSubCommand,
		"sync-commands": sync_commands.HandleSyncCommandSubCommand,
	}

	api.ProcessSubCommands(
		s,
		i,
		nil,
		subCommands)
}
