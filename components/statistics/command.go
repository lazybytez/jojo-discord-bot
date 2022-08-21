package statistics

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/lazybytez/jojo-discord-bot/api"
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

	w.Init(buf, 0, 4, 0, ' ', 0)
	fmt.Fprintf(w, "Uptime: **%v**\n", getDurationString(time.Since(statsStartTime)))
	fmt.Fprintf(w, "Memory used: **%s / %s**\n", humanize.Bytes(m.Alloc), humanize.Bytes(m.Sys))
	fmt.Fprintf(w, "Garbage collected: **%s**\n", humanize.Bytes(m.TotalAlloc))
	fmt.Fprintf(w, "Threads: **%s**\n", humanize.Comma(int64(runtime.NumGoroutine())))

	w.Flush()
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
