package dice

import (
	"math/rand"
)

// rollDice rolls `n` dice where the dice have each `d` sites and return the results as an int-array.
// For example: If you want to throw 3 times a d6, then `d` should be 6 and `n` should be 3.
func rollDice(d int, n int) []int {
	var r = make([]int, n)

	for i := 0; i < n; i++ {
		r[i] = rand.Intn(d - 1) + 1
	}

	return r
}
