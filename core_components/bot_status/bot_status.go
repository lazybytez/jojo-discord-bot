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

package bot_status

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"time"
)

// BotStatusRotationTime is the time between the status
// that are registered in the DiscordGoStatusManager.
const BotStatusRotationTime = 5 * time.Minute

var C = api.Component{
	// Metadata
	Code:         "bot_status",
	Name:         "Bot Status",
	Description:  "This component handles automated rotation and setting of the bot status in Discord.",
	LoadPriority: -1000, // Be the last core component, so others can register initial status

	State: &api.State{
		DefaultEnabled: true,
	},
}

// botStatusRotationTicker id the ticker used to periodically
// change the bots' status.
var botStatusRotationTicker *time.Ticker

func init() {
	api.RegisterComponent(&C, LoadComponent)
}

// LoadComponent loads the bot core component
// and handles migration of core entities
// and registration of important core event handlers.
func LoadComponent(_ *discordgo.Session) error {
	C.HandlerManager().RegisterOnce("start_status_rotation", onBotReady)

	return nil
}

// onBotReady starts the bot status rotation.
// At this point, discordgo is fully initialized and connected.
func onBotReady(_ *discordgo.Session, _ *discordgo.Ready) {
	startBotStatusRotation()
}

// startBotStatusRotation starts a routine that handles the automated rotation
// of the bots' status.
func startBotStatusRotation() {
	botStatusRotationTicker = time.NewTicker(BotStatusRotationTime)

	// Initial status rotation
	rotateStatus()

	// Continues status rotation
	go func() {
		for range botStatusRotationTicker.C {
			rotateStatus()
		}
	}()
}

// rotateStatus updates the status of the bot by rotating it.
func rotateStatus() {
	status := C.BotStatusManager().Next()

	if nil == status {
		C.Logger().Info("Not updating status, as no status are registered!")
	}

	err := C.DiscordApi().SetBotStatus(*status)

	if nil != err {
		C.Logger().Err(err, "Could not update the status of the bot due to an unexpected error!")

		return
	}

	C.Logger().Info("Updated bot status to content \"%s\" and url \"%s\" with activity type \"%d\"",
		status.Content,
		status.Url,
		status.ActivityType)
}
