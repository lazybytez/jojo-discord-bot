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
	"github.com/lazybytez/jojo-discord-bot/components"
	"github.com/rs/zerolog/log"
	"os"
)

// ExitGracefully shutdowns the bot gracefully.
// Using this functions results in the process exiting with
// exit code 0.
//
// The function tries to ensure that all allocated resources like
// connections or locks are freed/closed correctly.
func ExitGracefully(reason string) {
	releaseResources()

	log.Info().Msg(reason)
	os.Exit(0)
}

// ExitFatal shutdowns the application ungracefully with a
// non-zero exit code.
// The routines to properly stop the application are not applied on
// a fatal exit. The function should be only called when the application cannot
// recover, which is typically the case when core connections cannot be established
// or an initialization routine fails.
func ExitFatal(reason string) {
	log.Fatal().Msg(reason)
}

// ExitFatalGracefully shutdowns the application gracefully with a
// non-zero exit code.
// The routines to properly stop the application are not applied on
// a ExitFatal exit. The function should be only called when the application cannot
// recover, which is typically the case when core connections cannot be established
// or an initialization routine fails.
func ExitFatalGracefully(reason string) {
	releaseResources()

	log.Fatal().Msg(reason)
}

// releaseResources ensures that all allocated resources, locks
// and connections are freed before the application terminates
// gracefully.
func releaseResources() {
	components.UnloadComponents(discord)
	stopBot()
}
