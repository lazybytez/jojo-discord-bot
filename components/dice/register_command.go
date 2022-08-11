package dice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api"
)

// C is the instance of the dice component.
// Can be used to register the component or get information about it.
var C *api.Component

var diceCommand = &api.Command{
	Cmd: &discordgo.ApplicationCommand{
		Name: "dice",
		Description: "throw one or more dice of your wished type.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name: "die-sites-number",
				Description: "The number of how many sites the die has, default is 6",
				Type: discordgo.ApplicationCommandOptionInteger,
				MinValue: 2,
				Required: false,
			},
			{
				Name: "number-dice",
				Description: "How many dice you want to throw, default is 1",
				Type: discordgo.ApplicationCommandOptionInteger,
				MinValue: 2,
				Required: false
			},
		},
	},
	Handler: handleDice,
}

// init initializes the component with its metadata
func init() {
	C = &api.Component{
		// Metadata
		Name:	"Dice Component",
		Description: "This Component throws one or multiple dice",
		DmPermission: true,

		State: api.State{
			Enabled: true,
		},

		Lifecycle: api.LifecycleHooks{
			LoadComponent: LoadComponent,
		},
	}
}

// LoadComponent loads the Dice Component
func LoadComponent(discord *discordgo.Session) error {
	// Register the messageCreate func as a callback for MessageCreate events.
	_ = C.SlashCommandManager().Register(diceCommand)
}
