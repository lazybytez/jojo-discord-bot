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
// events that happen during the bots lifecycle.

var C *api.Component

func init() {
	C = &api.Component{
		// Metadata
		Code: "bot_log",
		Name: "Bot Log",
		Description: "This component prints out some basic information in the " +
			"log when bot is ready or added to guilds.",

		State: api.State{
			Enabled: true,
		},

		Lifecycle: api.LifecycleHooks{
			LoadComponent: LoadComponent,
		},
	}
}

func LoadComponent(discord *discordgo.Session) error {
	_, _ = C.HandlerManager().RegisterOnce("botready", onBotReady)

	_, _ = C.HandlerManager().Register("guild_join", onGuildJoin)
	_, _ = C.HandlerManager().Register("guild_leave", onGuildLeave)

	return nil
}

// onBotReady acts as a Discord ready event handler.
//
// It prints out the name and discriminator of the bot and
// count of guilds the bot is on.
func onBotReady(s *discordgo.Session, r *discordgo.Ready) {
	C.Logger().Info("Logged in as: \"%v#%v\"!", s.State.User.Username, s.State.User.Discriminator)
	C.Logger().Info("The bot is registered on \"%v\" guilds!", len(s.State.Guilds))
}

// onGuildJoin is triggered when the bot joins a guild.
//
// It provides information about the guild that has been joined.
func onGuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	C.Logger().Info("The bot joined the guild \"%v\" with ID \"%v\"", g.Name, g.ID)
}

// onGuildJoin is triggered when the bot leaves a guild.
//
// It provides information about the guild that has been left.
func onGuildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	C.Logger().Info("The bot left the guild \"%v\" with ID \"%v\"", g.Name, g.ID)
}
