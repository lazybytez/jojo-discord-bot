package dice

import "testing"

func TestDice(t *testing.T) {
	tables := getTestStruct()

	for _, table := range tables {
		dice := rollDice(table.d, table.n)

		checkIfMoreOrLessDiceWhereRolled(dice, table.n, t)		
		
		for _, die := range dice {
			checkIfDieLowerThanOne(die, t)
			checkIfDieHigherThanExpected(die, table.d, t)
		}
	}
}

func getTestStruct() []struct{d int; n int} {
	tables := []struct {
		d int
		n int
	} {
		{2, 3},
		{6, 1},
		{12, 0},
		{20, 5},
	}

	return tables
}

// Check if the dice-array (dice) has more or lessed thrown dice than wanted (n)
func checkIfMoreOrLessDiceWhereRolled(dice []int, n int, t *testing.T) {
	if len(dice) != n {
		t.Errorf("The number of rolled dices is incorrect, rolled: %d, want: %d", len(dice), n)
	}
}

// Check if the dice (d) has a lower value than 1
func checkIfDieLowerThanOne(d int, t *testing.T) {
	if d < 1 {
		t.Errorf("The dice should not be lower 1, got %d", d)
	}
}

// Check if the dice (d) has a higher value than expecet (e)
func checkIfDieHigherThanExpected(d int, e int, t *testing.T) {
	if d > e {
		t.Errorf("The dice should not be higher, expecet highest %d, got %d", e, d)
	}
}
