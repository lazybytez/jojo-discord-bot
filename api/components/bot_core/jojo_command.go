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
	"fmt"
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
	subCommands := map[string]func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
		option *discordgo.ApplicationCommandInteractionDataOption,
	){
		"module": handleModuleSubCommand,
	}

	success := api.ProcessSubCommands(
		s,
		i,
		nil,
		subCommands)

	if !success {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The executed (sub)command is invalid or does not exist!",
			},
		})
	}
}

// handleModuleSubCommand delegates the sub-commands of the module sub-command
// to their dedicated handlers.
func handleModuleSubCommand(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	subCommands := map[string]func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
		option *discordgo.ApplicationCommandInteractionDataOption,
	){
		"list":    handleModuleList,
		"show":    handleModuleShow,
		"enable":  handleModuleEnable,
		"disable": handleModuleDisable,
	}

	success := api.ProcessSubCommands(
		s,
		i,
		option,
		subCommands)

	if !success {
		if !success {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The executed (sub)command is invalid or does not exist!",
				},
			})
		}
	}
}

// handleModuleList prints out a list of all commands and their status.
func handleModuleList(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
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
						Value:  compNames,
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

// handleModuleShow prints out a list of all commands and their status.
func handleModuleShow(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Module Information",
				Color:  api.DefaultEmbedColor,
				Fields: []*discordgo.MessageEmbedField{},
			},
		},
	}

	var comp *api.Component
	for _, c := range api.Components {
		if c.Code == option.Options[0].Value {
			comp = c
			break
		}
	}

	if nil == comp || api.IsCoreComponent(comp) {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	regComp, ok := database.GetRegisteredComponent(C, comp.Code)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	globalStatus, ok := database.GetGlobalComponentStatus(C, regComp.ID)
	globalStatusOutput := ":white_check_mark:"
	if !ok {
		globalStatusOutput = ":no_entry:"
	}

	if !globalStatus.Enabled {
		globalStatusOutput += ":no_entry:"
	}

	guild, ok := database.GetGuild(C, i.GuildID)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guildSpecificStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
	guildSpecificStatusOutput := ":white_check_mark:"
	if !ok || !guildSpecificStatus.Enabled {
		guildSpecificStatusOutput += ":x:"
	}

	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Name",
			Value:  comp.Name,
			Inline: false,
		},
		{
			Name:   "Description",
			Value:  comp.Description,
			Inline: false,
		},
		{
			Name:   "Guild Status",
			Value:  guildSpecificStatusOutput,
			Inline: true,
		},
		{
			Name:   "Global Status",
			Value:  globalStatusOutput,
			Inline: true,
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

// handleModuleEnable enables the targeted module.
func handleModuleEnable(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Enable Module",
				Color:  api.DefaultEmbedColor,
				Fields: []*discordgo.MessageEmbedField{},
			},
		},
	}

	var comp *api.Component
	for _, c := range api.Components {
		if c.Code == option.Options[0].Value {
			comp = c
			break
		}
	}

	if nil == comp || api.IsCoreComponent(comp) {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	regComp, ok := database.GetRegisteredComponent(C, comp.Code)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guild, ok := database.GetGuild(C, i.GuildID)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guildSpecificStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
	if !ok {
		guildSpecificStatus.Component = *regComp
		guildSpecificStatus.Guild = *guild
		guildSpecificStatus.Enabled = true

		database.Create(guildSpecificStatus)

		resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Module",
				Value:  comp.Name,
				Inline: false,
			},
			{
				Name:   "Status",
				Value:  ":white_check_mark: - The module has been enabled!",
				Inline: false,
			},
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: resp,
		})

		return
	}

	if guildSpecificStatus.Enabled {
		respondWithAlreadyEnabled(s, i, resp, comp.Name)

		return
	}

	guildSpecificStatus.Enabled = true
	database.Save(guildSpecificStatus)

	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Module",
			Value:  comp.Name,
			Inline: false,
		},
		{
			Name:   "Status",
			Value:  ":white_check_mark: - The module has been enabled!",
			Inline: false,
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

// handleModuleDisable enables the targeted module.
func handleModuleDisable(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Disable Module",
				Color:  api.DefaultEmbedColor,
				Fields: []*discordgo.MessageEmbedField{},
			},
		},
	}

	var comp *api.Component
	for _, c := range api.Components {
		if c.Code == option.Options[0].Value {
			comp = c
			break
		}
	}

	if nil == comp || api.IsCoreComponent(comp) {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	regComp, ok := database.GetRegisteredComponent(C, comp.Code)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guild, ok := database.GetGuild(C, i.GuildID)
	if !ok {
		respondWithMissingComponent(s, i, resp, comp.Name)

		return
	}

	guildSpecificStatus, ok := database.GetComponentStatus(C, guild.ID, regComp.ID)
	if !ok {
		guildSpecificStatus.Component = *regComp
		guildSpecificStatus.Guild = *guild
		guildSpecificStatus.Enabled = false

		database.Create(guildSpecificStatus)

		resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Module",
				Value:  comp.Name,
				Inline: false,
			},
			{
				Name:   "Status",
				Value:  ":x: - The module has been disabled!",
				Inline: false,
			},
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: resp,
		})

		return
	}

	if !guildSpecificStatus.Enabled {
		respondWithAlreadyDisabled(s, i, resp, comp.Name)

		return
	}

	guildSpecificStatus.Enabled = false
	database.Save(guildSpecificStatus)

	resp.Embeds[0].Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Module",
			Value:  comp.Name,
			Inline: false,
		},
		{
			Name:   "Status",
			Value:  ":x: - The module has been disabled!",
			Inline: false,
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

// respondWithMissingComponent fills the passed discordgo.InteractionResponseData
// with an embed field that indicates that the specified component could not be found.
func respondWithMissingComponent(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName interface{},
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Error",
			Value: fmt.Sprintf("No module with name \"%v\" could be found!", componentName),
		},
	}

	resp.Embeds[0].Fields = embeds

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

// respondWithAlreadyEnabled fills the passed discordgo.InteractionResponseData
// with an embed field that indicates that the specified component is already enabled.
func respondWithAlreadyEnabled(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName interface{},
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Error",
			Value: fmt.Sprintf("Module with name \"%v\" is already enabled!", componentName),
		},
	}

	resp.Embeds[0].Fields = embeds

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

// respondWithAlreadyDisabled fills the passed discordgo.InteractionResponseData
// with an embed field that indicates that the specified component is already disabled.
func respondWithAlreadyDisabled(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName interface{},
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Error",
			Value: fmt.Sprintf("Module with name \"%v\" is already disabled!", componentName),
		},
	}

	resp.Embeds[0].Fields = embeds

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
