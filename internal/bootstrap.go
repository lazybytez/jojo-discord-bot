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
    "os/signal"
    "syscall"
)

// Bootstrap hadles the start of the application.
// It is responsible to execute the startup sequence
// and get the application up and running properly.
func Bootstrap() {
    startBot(Config.token)

    components.RegisterComponents(discord)

    waitForTerminate()
}

// waitForTerminate blocks the console and waits
// for a termination signal.
//
// When a sigterm is received, the application is stopped
// gracefully. This means all open connections or used resources are
// freed/closed before exit.
func waitForTerminate() {
    log.Info().Msg("Bot is running.  Press CTRL-C to exit.")

    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
    <-c

    ExitGracefully("Bot has been terminated gracefully!")
}
