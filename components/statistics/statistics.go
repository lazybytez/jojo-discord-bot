package statistics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// C is the instance of the statistics component.
var C *api.Component

// init initializes the component with its metadata
func init() {
	C = &api.Component{
		// Metadata
		Name:        "Statistics Component",
		Description: "This Component returns statistics about the bot and the runtime.",

		State: &api.State{
			DefaultEnabled: true,
		},

		Lifecycle: api.LifecycleHooks{
			LoadComponent: LoadComponent,
		},
	}
}

// LoadComponent loads the Ping-Pong Component
func LoadComponent(discord *discordgo.Session) error {
	_ = C.SlashCommandManager().Register(statsCommand)
	_ = C.SlashCommandManager().Register(infoCommand)

	return nil
}
