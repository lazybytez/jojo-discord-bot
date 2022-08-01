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

package pingpong

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// This component is an example on how to create a component
// The variable C holds the api.Component itself.
// The init function initializes the component with all required
// metadata.
// To make a component autoload properly, it is necessary to:
//  1. Enable it in the api.State of the component
//  2. Add an api.Lifecycle with a valid callback
//     that initializes the component

// C is the instance of the ping pong component.
// Can be used to register the component or get information about it.
var C *api.Component

// init initializes the component with its metadata
func init() {
	C = &api.Component{
		// Metadata
		Name:         "Ping Pong Component",
		Description:  "This Component plays pingpong with you and returns Latency (maybe)",
		DmPermission: true,

		State: api.State{
			Enabled: true,
		},

		Lifecycle: api.LifecycleHooks{
			LoadComponent: LoadComponent,
		},
	}
}

// LoadComponent loads the Ping-Pong Component
func LoadComponent(discord *discordgo.Session) error {
	// Register the messageCreate func as a callback for MessageCreate events.
	_, _ = C.HandlerManager().Register("pingpong", onMessageCreate)

	return nil
}

// onMessageCreate listens for new messages and replies with
// "Ping!" or "Pong!" depending on the received message.
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if nil != err {
			C.Logger().Warn("Failed to deliver \"Pong!\" message: %v", err.Error())

			return
		}
		C.Logger().Info("Send \"Pong!\" into channel with ID \"%v\"", m.ChannelID)
	}

	if m.Content == "pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if nil != err {
			C.Logger().Warn("Failed to deliver \"Ping!\" message: %v", err.Error())

			return
		}
		C.Logger().Info("Send \"Ping!\" into channel with ID \"%v\"", m.ChannelID)
	}
}
