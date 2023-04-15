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
	"github.com/lazybytez/jojo-discord-bot/services/cache"
)

// ToggleModuleRateLimit is the maximum count how
// often modules can be enabled and disabled in 10 minutes
// before the user has to wait for the rate limit to expire.
const ToggleModuleRateLimit = 10

// getComponentToggleCountCacheKey returns the cache key used to store
// the count of how often the module toggle commands were used in the last
// 10 minutes.
func getComponentToggleCountCacheKey(guildID uint64) string {
	return fmt.Sprintf("component_toggle_counter_%d", guildID)
}

// isModuleToggleRateLimited returns whether the current guild is
// right now rate limited or not.
func isModuleToggleRateLimited(guild *entities.Guild) bool {
	cacheKey := getComponentToggleCountCacheKey(guild.GuildID)
	toggleCount, ok := cache.Get(cacheKey, 0)
	if ok && toggleCount > ToggleModuleRateLimit {
		return true
	}

	return false
}

// increaseRateLimitCount increases the count of module
// toggles in the cache by one.
// The function returns whether increasing the rate limit count
// worked or not.
func increaseRateLimitCount(guild *entities.Guild) bool {
	cacheKey := getComponentToggleCountCacheKey(guild.GuildID)
	toggleCount, ok := cache.Get(cacheKey, 0)
	if !ok {
		return false
	}

	toggleCount += 1
	err := cache.Update(cacheKey, toggleCount)
	if nil != err {
		C.Logger().Err(err, fmt.Sprintf(
			"Failed to store incremented module toggle rate limit count in cache for guild %d",
			guild.GuildID))
	}

	return true
}

// respondWithRateLimited responds with a message telling the user that the has to wait
// a few minutes before he uses the command again.
func respondWithRateLimited(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	resp *discordgo.InteractionResponseData,
) {
	embeds := []*discordgo.MessageEmbedField{
		{
			Name: ":x: Slow down my friend!",
			Value: fmt.Sprintf("The `/jojo module enable` and `/jojo module disable` "+
				"commands can only be used up to %d times in 10 minutes per guild!",
				ToggleModuleRateLimit),
		},
	}

	resp.Embeds[0].Fields = embeds

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}
