package statistics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// C is the instance of the statistics component
var C = api.Component{
	// Metadata
	Name:        "Statistics Component",
	Description: "This Component returns statistics about the bot and the runtime.",

	State: &api.State{
		DefaultEnabled: true,
	},
}

// init initializes the component with its metadata
func init() {
	api.RegisterComponent(&C, LoadComponent)
}

// LoadComponent loads the two registered slash commands
func LoadComponent(_ *discordgo.Session) error {
	_ = C.SlashCommandManager().Register(statsCommand)
	_ = C.SlashCommandManager().Register(infoCommand)

	return nil
}
