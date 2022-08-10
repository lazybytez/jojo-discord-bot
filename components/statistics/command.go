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

package statistics

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
	"runtime"
)

var statsCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "stats",
		Description: "Show statistics of the bot and its runtime.",
	},
	Handler: handleStats,
}

var m runtime.MemStats

func handleStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	runtime.ReadMemStats(&m)
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Discord Bot Statistics",
					Description: "Overview of runtime statistics",
					Color:       0x5D397C,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Memory",
							Value:  fmt.Sprintf("%v MB / %v MB", bToMb(m.Alloc), bToMb(m.TotalAlloc)),
							Inline: true,
						},
						{
							Name:   "Maximum Memory",
							Value:  fmt.Sprintf("System has %v MB", bToMb(m.Sys)),
							Inline: true,
						},
						{
							Name:   "Garbage Collector Cycles",
							Value:  fmt.Sprintf("Run %v Times", m.NumGC),
							Inline: true,
						},
					},
				},
			},
		},
	})
}

// bToMb calculates Bytes to Megabytes (not Mebibyte, that would be 1024 / 1024)
func bToMb(b uint64) uint64 {
	return b / 1024 / 1000
}
