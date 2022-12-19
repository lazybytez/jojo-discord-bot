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

package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// DiscordGoApiWrapper is a wrapper around some crucial discordgo
// functions. It provides functions that might be
// frequently used without an ongoing event.
type DiscordGoApiWrapper struct {
	owner *Component
}

// DiscordApiWrapper provides some useful functions of the discord API,
// that might be needed without an ongoing event.
// An example would be the WebAPI.
//
// Also, the functions of this interface are sharding compatible.
// As soon as sharding is enabled, functions like GuildCount will still
// output the proper value.
type DiscordApiWrapper interface {
	// GuildCount returns the number of guilds the bot is currently on.
	GuildCount() int
	// SetBotStatus updates the status of the bot according to the passed
	// SimpleBotStatus data.
	SetBotStatus(status SimpleBotStatus) error
}

// DiscordApi is used to obtain the components slash DiscordApiWrapper management
//
// On first call, this function initializes the private Component.discordAPi
// field. On consecutive calls, the already present DiscordGoApiWrapper will be used.
func (c *Component) DiscordApi() DiscordApiWrapper {
	if nil == c.discordApi {
		c.discordApi = &DiscordGoApiWrapper{owner: c}
	}

	return c.discordApi
}

// GuildCount returns the number of guilds the bot is currently on.
//
// TODO: As soon as sharding support is implemented, the guild count needs to be computed from data collected across all shards
func (dgw *DiscordGoApiWrapper) GuildCount() int {
	return len(dgw.owner.discord.State.Guilds)
}

// SimpleBotStatus is a simplified version of discordgo.UpdateStatusData
// that can be used to simply change the status of the bot to something else.
// Note that the URL should be only set for discordgo.ActivityTypeStreaming.
type SimpleBotStatus struct {
	ActivityType discordgo.ActivityType
	Content      string
	Url          string
}

// SetBotStatus updates the status of the bot according to the passed
// SimpleBotStatus data.
func (dgw *DiscordGoApiWrapper) SetBotStatus(status SimpleBotStatus) error {
	switch status.ActivityType {
	case discordgo.ActivityTypeGame:
		return dgw.owner.discord.UpdateGameStatus(0, status.Content)
	case discordgo.ActivityTypeStreaming:
		return dgw.owner.discord.UpdateStreamingStatus(0, status.Content, status.Url)
	case discordgo.ActivityTypeListening:
		return dgw.owner.discord.UpdateListeningStatus(status.Content)
	default:
		return fmt.Errorf("tried to update bot status to activity type \"%d\", which is not supported",
			status.ActivityType)
	}
}
