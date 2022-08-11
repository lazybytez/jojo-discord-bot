package dice

import (
	"fmt"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

func sendAnswerToUser(n int, d int, r []int, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sArray := arrayIntToArrayString(r)
	a := createAnswer(n, d, sArray)
	sendAnswer(s, i, a)
}

// Create the answer to send
func createAnswer(n int, d int, rolledDice []string) string {
	answer := fmt.Sprintf("You rolled %d d%d, The Results are:\n", n, d)
	
	answer += implode(", ", rolledDice)

	return answer
}

// Send the Answer
func sendAnswer(s *discordgo.Session, i *discordgo.InteractionCreate, answerText string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: answerText,
		},
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
	strings := make([]string, len(ints))

	for k, i := range ints {
		strings[k] = strconv.Itoa(i)
	}

	return strings
}