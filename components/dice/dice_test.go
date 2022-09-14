package dice

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DiceTestSuite struct{ suite.Suite }

func TestDice(t *testing.T) {
	suite.Run(t, new(DiceTestSuite))
}

func (suite *DiceTestSuite) TestDice() {
	tables := getTestStruct()

	for _, table := range tables {
		dice := rollDice(table.d, table.n)

		suite.Len(dice, table.n)

		for _, die := range dice {
			suite.Greater(die, 0)
			suite.LessOrEqual(die, table.d)
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

func (suite *DiceTestSuite) TestDiceForAllValues() {
	throws := 100000
	d := 3

	one := false
	two := false
	three := false
	for i := 1; i < throws; i++ {
		die := rollDice(d, 1)
		switch die[0] {
		case 1:
			one = true
		case 2:
			two = true
		case 3:
			three = true
		}

		if one && two && three {
			return
		}
	}

	suite.FailNow("One of the following numbers was not thrown: 1, 2, 3")
}
