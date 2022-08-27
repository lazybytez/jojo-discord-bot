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
	"text/tabwriter"
	"time"
)

var statsCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "stats",
		Description: "Show information of the bot and runtime statistics.",
	},
	Handler: handleStats,
}

var infoCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name:        "info",
		Description: "Show information of the bot and runtime statistics.",
	},
	Handler: handleStats,
}

var m runtime.MemStats

var statsStartTime = time.Now()

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
			Title: "Information",
			Color: 0x5D397C,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: s.State.User.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
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

// buildStatOutput generates a big string with all runtime statistics
func buildStatOutput() string {
	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	count, _ := C.EntityManager().Guilds().Count()
	cluster, _ := os.Hostname()

	w.Init(buf, 0, 4, 0, ' ', 0)
	appendStatLine(w, "Uptime: **%v**\n", getDurationString(time.Since(statsStartTime)))
	appendStatLine(w, "Memory used: **%s / %s**\n", humanize.Bytes(m.Alloc), humanize.Bytes(m.Sys))
	appendStatLine(w, "Garbage collected: **%s**\n", humanize.Bytes(m.TotalAlloc))
	appendStatLine(w, "Threads: **%s**\n", humanize.Comma(int64(runtime.NumGoroutine())))
	appendStatLine(w, "Connected Servers: **%v**\n", count)
	appendStatLine(w, "Cluster ID: **%s**\n", cluster)
	appendStatLine(w, "Registered Slash Commands: **%v**\n", C.SlashCommandManager().GetCommandCount())

	err := w.Flush()
	if nil != err {
		return ""
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

func appendStatLine(w io.Writer, msg string, values ...interface{}) {
	_, err := fmt.Fprintf(w, msg, values...)
	if nil != err {
		C.Logger().Err(err, "Failed to generate statistics embed text write buffer.")
	}
}
