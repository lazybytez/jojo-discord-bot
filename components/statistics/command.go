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
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/lazybytez/jojo-discord-bot/api"
	"io"
	"os"
	"runtime"
	"strconv"
	"text/tabwriter"
	"time"
)

// statsCommand registers the alias /stats
var statsCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "stats",
		Description: "Show information of the bot and runtime statistics.",
	},
	Handler: handleStats,
}

// infoCommand registers the alias /info
var infoCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "info",
		Description: "Show information of the bot and runtime statistics.",
	},
	Handler: handleStats,
}

// With the m variable the command can access memory runtime statistics
var m runtime.MemStats

// botStartTime returns the current local time to calculate the uptime of the bot instance
var botStartTime = time.Now()

// handleStats gets called when executing the /stats or /info command
func handleStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	runtime.ReadMemStats(&m)
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: buildInfoEmbed(s),
		},
	})
}

// buildInfoEmbed the embed is build in a different function to secure readability
func buildInfoEmbed(s *discordgo.Session) []*discordgo.MessageEmbed {
	return []*discordgo.MessageEmbed{
		{
			Title: "Info",
			Color: 0x5D397C,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: s.State.User.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Bot",
					Value: buildInfoOutput(),
				},
				{
					Name:  "Stats",
					Value: buildStatOutput(),
				},
				{
					Name:  "Links",
					Value: "[GitHub](https://github.com/lazybytez/jojo-discord-bot)",
				},
			},
		},
	}
}

// buildInfoOutput generates a big text with some general bot information like how many slash commands it has etc.
func buildInfoOutput() string {
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	w.Init(buf, 0, 4, 0, ' ', 0)

	appendStatLine(w, "Slash Commands: **%v**\n", C.SlashCommandManager().GetCommandCount())

	err := w.Flush()
	if err != nil {
		C.Logger().Err(err, "Could not flush statistics embed text write buffer.")
	}

	return buf.String()
}

// buildStatOutput generates a big string with all runtime statistics
func buildStatOutput() string {
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	count, err := C.EntityManager().Guilds().Count()
	countMsg := strconv.FormatInt(count, 10)
	if nil != err {
		countMsg = "Error"
	}
	cluster, err := os.Hostname()
	if nil != err {
		cluster = "Error"
	}

	w.Init(buf, 0, 4, 0, ' ', 0)

	appendStatLine(w, "Uptime: **%v**\n", getDurationString(time.Since(botStartTime)))
	appendStatLine(w, "Memory used: **%s / %s**\n", humanize.Bytes(m.Alloc), humanize.Bytes(m.Sys))
	appendStatLine(w, "Garbage collected: **%s**\n", humanize.Bytes(m.TotalAlloc))
	appendStatLine(w, "Threads: **%s**\n", humanize.Comma(int64(runtime.NumGoroutine())))
	appendStatLine(w, "Connected Servers: **%v**\n", countMsg)
	appendStatLine(w, "Cluster ID: **%s**\n", cluster)

	err = w.Flush()
	if err != nil {
		C.Logger().Err(err, "Could not flush statistics embed text write buffer.")
	}

	return buf.String()
}

// getDurationString transforms duration into a readable string
func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

// appendStatLine adds another line to a single text field in the embed
func appendStatLine(w io.Writer, msg string, values ...interface{}) {
	_, err := fmt.Fprintf(w, msg, values...)
	if nil != err {
		C.Logger().Err(err, "Failed to generate statistics embed text write buffer.")
	}
}
