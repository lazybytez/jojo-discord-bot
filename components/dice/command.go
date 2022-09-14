package dice

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/slash_commands"
)

// handleDice handles the dice slash command
func handleDice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := getOptionsAsMap(i)
	n := getIntOption(options, "number-dice", 1)
	d := getIntOption(options, "die-sites-number", 6)

	r := rollDice(d, n)
	e := createAnswerEmbedMessage(n, d, r)
	eSlice := [1]*discordgo.MessageEmbed{&e}
	sendAnswer(s, i, eSlice[:])
}

// getOptionsAsMap create a map and insert the command options
func getOptionsAsMap(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

// getIntOption returns the option as int
func getIntOption(options map[string]*discordgo.ApplicationCommandInteractionDataOption, name string, defaultValue int) int {
	o := defaultValue
	if opt, ok := options[name]; ok {
		o = int(opt.IntValue())
	}

	return o
}

// sendAnser sends the Answer
func sendAnswer(s *discordgo.Session, i *discordgo.InteractionCreate, e []*discordgo.MessageEmbed) {
	resp := &discordgo.InteractionResponseData{
		Embeds: e,
	}
	slash_commands.Respond(C, s, i, resp)
}
