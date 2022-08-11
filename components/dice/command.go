
package dice

import (
	"github.com/bwmarrin/discordgo"
)

func handleDice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := getOptionsAsMap(i)
	n := getIntOption(options, "number-dice", 1)
	d := getIntOption(options, "die-sites-number", 6)

	r := rollDice(n, d)
	sendAnswerToUser(n, d, r, s, i)
}

// create a map and insert the command options
func getOptionsAsMap(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

// Create the option as int
func getIntOption(options map[string]*discordgo.ApplicationCommandInteractionDataOption, name string, defaultValue int) int {
	o := defaultValue
	if opt, ok := options[name]; ok {
		o = int(opt.IntValue())
	}

	return o
}
