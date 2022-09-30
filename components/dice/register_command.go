package dice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// C is the instance of the dice component.
// Can be used to register the component or get information about it.
var C = api.Component{
	// Metadata
	Code:        "dice",
	Name:        "Dice Component",
	Description: "This Component throws one or multiple dice",

	State: &api.State{
		DefaultEnabled: true,
	},
}

// init initializes the component with its metadata
func init() {
	api.RegisterComponent(&C, LoadComponent)
}

// LoadComponent loads the Dice Component
func LoadComponent(discord *discordgo.Session) error {
	// Register the messageCreate func as a callback for MessageCreate events.
	_ = C.SlashCommandManager().Register(diceCommand)

	return nil
}
