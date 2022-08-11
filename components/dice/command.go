
package dice

import (
	"fmt"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

func handleDice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := getOptionsAsMap(i)
	n := getIntOption(i, "number-dice", 1)
	d := getIntOption(i, "die-sites-number", 6)

	r := rollDice(n, d)
	s := arrayIntToArrayString(r)
	a := createAnswer(n, d, s)
	sendAnswer(a)
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
		o = opt.IntValue()
	}

	return 0
}

// Create the answer to send
func createAnswer(n int, d int, rolledDice []int) string {
	answer := fmt.Sprintf("You rolled %d d%d, The Results are:\n", n, d)
	
	answer += implode(", ", rolledDice)

	return answer
}

// Send the Answer
func sendAnswer(answerText string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: answerText,
		}
	})
}

// Create one string with the array of string seperated with seperator (s)
func implode(s string, array []string) string {
	first := true
	r := ""

	for _, a := range array {
		if (first) {
			first = false
			r += a
		} else {
			r += s + a
		}
	}

	return r
}

// Convert an array of int to an array of string
func arrayIntToArrayString(ints []int) []string {
	strings := []string {}

	for k, i := range ints {
		strings[k] = strconf.Itoa(i)
	}

	return strings
}