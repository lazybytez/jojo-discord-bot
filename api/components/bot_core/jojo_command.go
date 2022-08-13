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
	"github.com/lazybytez/jojo-discord-bot/api/database"
)

// jojoCommand holds the command configuration for the jojo command.
var jojoCommand *api.Command

// getModuleCommandChoices builds a slice containing all available modules
// as command option choices
func getModuleCommandChoices() []*discordgo.ApplicationCommandOptionChoice {
	availableModuleChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)

	for _, comp := range api.Components {
		if api.IsCoreComponent(comp) {
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
	jojoCommand = &api.Command{
		Cmd: &discordgo.ApplicationCommand{
			Name:        "jojo",
			Description: "Manage modules and core settings of the bot!",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Name:        "module",
					Description: "Manage which modules should be enabled / disabled on your server!",
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
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
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
	handleModuleList(s, i)
}

// handleModuleList prints out a list of all commands and their status.
func handleModuleList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	compNames := ""
	compStatus := ""

	for _, comp := range api.Components {
		if api.IsCoreComponent(comp) {
			continue
		}

		if "" != compNames {
			compNames += "\n"
		}
		compNames += comp.Name

		regComp, ok := database.GetRegisteredComponent(C, comp.Code)
		if !ok {
			continue
		}

		if "" != compStatus {
			compStatus += "\n"
		}

		globalStatus, ok := database.GetGlobalComponentStatus(C, regComp.ID)
		if !ok {
			continue
		}

		if !globalStatus.Enabled {
			compStatus += ":no_entry:"

			continue
		}

		guild, ok := database.GetGuild(C, i.GuildID)
		if !ok {
			continue
		}

		guildSpecificStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
		if !ok || !guildSpecificStatus.Enabled {
			compStatus += ":x:"

			continue
		}

		compStatus += ":white_check_mark:"
	}

	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Module Status",
				Description: "Overview of all modules and whether they are enabled or not",
				Color:       api.DefaultEmbedColor,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Module",
						Value:  compNames + "\n\n",
						Inline: true,
					},
					{
						Name:   "Status",
						Value:  compStatus,
						Inline: true,
					},
					{
						Name: "Legend",
						Value: ":white_check_mark: - Enabled\n" +
							":x: - Disabled\n" +
							":no_entry: - Globally disabled (Maintenance)",
						Inline: false,
					},
				},
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
