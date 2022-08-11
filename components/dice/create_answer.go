package dice

import (
	"fmt"
	"strconv"
)

func createAnswerText(n int, d int, r []int) string {
	s := arrayIntToArrayString(r)

	st := createAnswerTitle(n, d)
	srt := createAnswerResultTitle(n)
	src := createAnswerResultContent(s)

	return st + srt + src
}

// Create the answer title
func createAnswerTitle(n int, d int) string {
	answer := fmt.Sprintf("You rolled %d d%d,\n", n, d)

	return answer
}

// Create the answer result title
func createAnswerResultTitle(n int) string {
	r := "The Result"
	if n > 1 {
		r += "s are:\n"
	} else {
		r += " is:\n"
	}

	return r
}

// Create the answer result content
func createAnswerResultContent(s []string) string {
	return implode(", ", s)
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