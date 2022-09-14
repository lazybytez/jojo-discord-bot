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

func getTestStruct() []struct {
	d int
	n int
} {
	tables := []struct {
		d int
		n int
	}{
		{2, 3},
		{6, 1},
		{12, 0},
		{20, 5},
	}

	return tables
}

func checkIfMoreOrLessDiceWhereRolled(dice []int, n int, t *testing.T) {
	if len(dice) != n {
		t.Errorf("The number of rolled dices is incorrect, rolled: %d, want: %d", len(dice), n)
	}
}

func checkIfDieLowerThanOne(d int, t *testing.T) {
	if d < 1 {
		t.Errorf("The dice should not be lower 1, got %d", d)
	}
}

func checkIfDieHigherThanExpected(d int, e int, t *testing.T) {
	if d > e {
		t.Errorf("The dice should not be higher, expecet highest %d, got %d", e, d)
	}
}

func TestDiceForAllValues(t *testing.T) {
	throws := 1000
	d := 3
	dice := rollDice(d, throws)
	one := false
	two := false
	three := false

	for _, die := range dice {
		if die == 1 {
			one = true
		}
		if die == 2 {
			two = true
		}
		if die == 3 {
			three = true
		}

		if one && two && three {
			return
		}
	}

	var firstFailed string
	switch {
	case !one:
		firstFailed = "one"
	case !two:
		firstFailed = "two"
	case !three:
		firstFailed = "three"
	}

	t.Errorf("There was no die with result %s by %d throws", firstFailed, throws)
}
