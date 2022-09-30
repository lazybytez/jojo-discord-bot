package dice

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// createAnswerEmbedMessage create the answer parts and put them in a embedMessage
func createAnswerEmbedMessage(n int, d int, r []int) discordgo.MessageEmbed {
	st := createAnswerTitle(n, d)
	srt := createAnswerResultTitle(n)

	s := arrayIntToArrayString(r)
	src := createAnswerResultContent(s)

	e := createMessageEmbed(st, srt, src)

	return e
}

// createAnswerTitle creates the answer title
func createAnswerTitle(n int, d int) string {
	answer := fmt.Sprintf("You rolled %d d%d", n, d)

	return answer
}

// createAnswerResultTitle creates the answer result title
func createAnswerResultTitle(n int) string {
	r := "The Result"
	if n > 1 {
		r += "s are"
	} else {
		r += " is"
	}

	return r
}

// createAnswerResultContent creates the answer result content
func createAnswerResultContent(s []string) string {
	return strings.Join(s, ", ")
}

// arrayIntToArrayString converts an array of int to an array of string
func arrayIntToArrayString(ints []int) []string {
	strings := make([]string, len(ints))

	for k, i := range ints {
		strings[k] = strconv.Itoa(i)
	}

	return strings
}
