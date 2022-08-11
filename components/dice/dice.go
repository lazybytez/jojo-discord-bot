package dice

import (
	"math/rand"
)

// rollDice rolls `n` dice where the dice have each `d` sites and return the results as an int-array.
// For example: If you want to throw 3 times a d6, then `d` should be 6 and `n` should be 3.
func rollDice(d int, n int) {
	r := [n]int;

	for i := 1; i <= n; i++ {
		r[i] = rand.Intnl(d);
	}

	return r;
}
