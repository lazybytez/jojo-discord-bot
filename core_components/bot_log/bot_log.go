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

package bot_log

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// Core component that handles logging of essential
// events that happen during the bot lifecycle.

var C = api.Component{
	// Metadata
	Code:       "bot_log",
	Name:       "Bot Log",
	Categories: api.Categories{api.CategoryInternal},
	Description: "This component prints out some basic information in the " +
		"log when bot is ready or added to guilds.",

	State: &api.State{
		DefaultEnabled: true,
	},
}

func init() {
	api.RegisterComponent(&C, LoadComponent)
}

func LoadComponent(_ *discordgo.Session) error {
	_, _ = C.HandlerManager().RegisterOnce("botready", onBotReady)

	_, _ = C.HandlerManager().Register("guild_join", onGuildJoin)
	_, _ = C.HandlerManager().Register("guild_leave", onGuildLeave)

	return nil
}

// onBotReady acts as a Discord ready event handler.
//
// It prints out the name and discriminator of the bot and
// count of guilds the bot is on.
func onBotReady(s *discordgo.Session, _ *discordgo.Ready) {
	C.Logger().Info("Logged in as: \"%v#%v\"!", s.State.User.Username, s.State.User.Discriminator)
	C.Logger().Info("The bot is registered on \"%v\" guilds!", len(s.State.Guilds))
}

// onGuildJoin is triggered when the bot joins a guild.
//
// It provides information about the guild that has been joined.
func onGuildJoin(_ *discordgo.Session, g *discordgo.GuildCreate) {
	C.Logger().Info("The bot joined the guild \"%v\" with ID \"%v\"", g.Name, g.ID)
}

// onGuildJoin is triggered when the bot leaves a guild.
//
// It provides information about the guild that has been left.
func onGuildLeave(_ *discordgo.Session, g *discordgo.GuildDelete) {
	guild, err := C.EntityManager().Guilds().Get(g.ID)
	guildName := ""
	if nil == err {
		guildName = guild.Name
	}

	C.Logger().Info("The bot left the guild \"%v\" with ID \"%v\"", guildName, g.ID)
}
