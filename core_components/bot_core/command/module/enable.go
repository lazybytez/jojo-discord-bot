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

package module

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

// handleModuleEnable enables the targeted module.
func handleModuleEnable(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	option *discordgo.ApplicationCommandInteractionDataOption,
) {
	resp := slash_commands.GenerateEphemeralInteractionResponseTemplate("Enable Module", "")

	regComp := findComponent(option)
	if nil == regComp || regComp.IsCoreComponent() {
		respondWithMissingComponent(s, i, resp, option.Options[0].Value)

		return
	}

	guild, err := C.EntityManager().Guilds().Get(i.GuildID)
	if nil != err {
		respondWithMissingComponent(s, i, resp, regComp.Name)

		return
	}

	if isModuleToggleRateLimited(guild) {
		respondWithRateLimited(s, i, resp)

		return
	}
	increaseRateLimitCount(guild)

	if !enableComponentForGuild(s, i, guild, regComp, resp) {
		respondWithAlreadyEnabled(s, i, resp, regComp.Name)

		return
	}

	respondWithTogglingComponent(
		s,
		i,
		resp,
		regComp.Name,
		UserActionEnable)
	C.SlashCommandManager().SyncApplicationComponentCommands(s, i.GuildID)
	finishWithModuleEnableSuccessfulEmbedField(s, i, resp, regComp)

	dgoGuild, err := s.Guild(i.GuildID)
	if nil != err {
		C.Logger().Err(err, "Failed to get guild with id \"%S\" to create "+
			"bot audit log when enabling a module on guild!",
			i.GuildID)

		return
	}

	user := i.User
	if nil == user {
		user = i.Member.User
	}

	C.BotAuditLogger().Log(
		dgoGuild,
		user,
		fmt.Sprintf("The component `%s` has been enabled", regComp.Name),
		true)
}

// enableComponentForGuild enables the specified component
// for the specified guild, if not already enabled.
//
// Returns true if the component has been enabled and was not
// enabled before.
func enableComponentForGuild(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	guild *entities.Guild,
	regComp *entities.RegisteredComponent,
	resp *discordgo.InteractionResponseData,
) bool {
	em := C.EntityManager()

	guildSpecificStatus, err := em.GuildComponentStatus().Get(guild.ID, regComp.ID)
	if nil != err {
		guildSpecificStatus.Component = *regComp
		guildSpecificStatus.Guild = *guild
		guildSpecificStatus.Enabled = true

		err = em.GuildComponentStatus().Create(guildSpecificStatus)
		if nil != err {
			C.Logger().Warn("Could not create guild component status for component \"%v\" on guild \"%v\"",
				regComp.Code,
				guild.GuildID)

			return false
		}

		return true
	}

	if guildSpecificStatus.Enabled {
		return false
	}

	guildSpecificStatus.Enabled = true
	err = em.GuildComponentStatus().Save(guildSpecificStatus)
	if nil != err {
		C.Logger().Warn("Could not update guild component status for component \"%v\" on guild \"%v\"",
			regComp.Code,
			guild.GuildID)

		return false
	}

	return true
}

// finishWithModuleEnableSuccessfulEmbedField updates the previously send
// processing message with a success message, that indicates
// that the module could be enabled properly.
func finishWithModuleEnableSuccessfulEmbedField(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	comp *entities.RegisteredComponent,
) {
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

	slash_commands.EditResponse(C, s, i, &discordgo.WebhookEdit{
		Embeds: &resp.Embeds,
	})
}

// respondWithAlreadyEnabled fills the passed discordgo.InteractionResponseData
// with an embed field that indicates that the specified component is already enabled.
func respondWithAlreadyEnabled(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
	componentName string,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name:  ":x: Error",
			Value: fmt.Sprintf("Module with name \"%v\" is already enabled!", componentName),
		},
	}

	resp.Embeds[0].Fields = embeds

	slash_commands.Respond(C, s, i, resp)
}
