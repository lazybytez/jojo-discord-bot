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

package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const tokenPrefix = "Bot "

var discord *discordgo.Session

// createSession creates a discordgo.Session,
// but does not open the connection yet.
//
// The token that is passed will be used to
// configure the session.
func createSession(token string) {
	if nil != discord {
		ExitFatal("DiscordGo session can be created only once!")
	}

	var err error
	discord, err = discordgo.New(tokenPrefix + token)
	if nil != err {
		ExitFatal(fmt.Sprintf("Failed to create discordgo session, error was: %v!", err.Error()))
	}
}

// startBot opens the connection and configures the bots intents.
func startBot() {
	if nil == discord {
		ExitFatal("Session must be first created before the bot can be started")
	}

	err := discord.Open()
	if nil != err {
		ExitFatal(fmt.Sprintf("Failed to open bot connection to Discord, error was: %v!", err.Error()))
	}

	updateIntents()
}

// stopBot tries to stop the bot.
// The bot is stopped by closing the discordgo.Session.
//
// If the bot has not been initialized until this point,
// the close function of the discordgo.Session won't be called.
//
// If closing the session throws an error, the error is ignored.
func stopBot() {
	if nil != discord {
		_ = discord.Close()
	}
}

// updateIntents to receive all necessary permissions
// for the bot
func updateIntents() {
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
}
