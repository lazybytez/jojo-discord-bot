package dice

import "testing"

func TestArrayIntToArrayString(t *testing.T) {
	p := [5]int {4, 20, -13, 0, -6}
	e := [5]string {"4", "20", "-13", "0", "-6"}

	r := arrayIntToArrayString(p[:])

	checkIfArraysHaveSameLength(t, e[:], r[:])
	checkIfParameterOfArraysAreTheSame(t, e[:], r[:])
}

// Check if the array we got (g) is as long as the array we expect (e)
func checkIfArraysHaveSameLength(t *testing.T, e []string, g []string) {
	if len(e) != len(g) {
		t.Errorf("The resulted array has not the correct length, expected: %d, got: %d", len(e), len(g))
	}
}

// Check if the two arrays have the same values
func checkIfParameterOfArraysAreTheSame(t *testing.T, expect []string, g []string) {
	for k, e := range expect {
		if e != g[k] {
			t.Errorf("The value of element %d is not what we expect, expected %s, got %s", k, e, g[k])
		}
	}
}

func TestImplode(t *testing.T) {
	a := [10]string {"What", "is", "'Courage'?", "Courage", "is", "owning", "your", "feat!", "-" ,"Zeppeli" }
	s := " "
	e := "What is 'Courage'? Courage is owning your feat! - Zeppeli"

	g := implode(s, a[:])

	checkIfTwoStringsAreTheSame(t, e, g)
}

func TestCreateAnswerResultContent(t *testing.T) {
	a := [4]string {"It", "was", "me", "Dio!"}
	e := "It, was, me, Dio!"

	r := createAnswerResultContent(a[:])

	checkIfTwoStringsAreTheSame(t, e, r)
}

func TestCreateAnswerResultTitleOne(t *testing.T) {
	e := "The Result is:\n"

	r := createAnswerResultTitle(1)

	checkIfTwoStringsAreTheSame(t, e, r)
}

func TestCreateAnswerResultTitleMultiple(t *testing.T) {
	e := "The Results are:\n"

	r := createAnswerResultTitle(2)

	checkIfTwoStringsAreTheSame(t, e, r)
}

func TestCreateAnswerTitle(t *testing.T) {
	e := "You rolled 3 d6,\n"

	r := createAnswerTitle(3, 6)

	checkIfTwoStringsAreTheSame(t, e, r)
}

func TestCreateAnswerText(t *testing.T) {
	e := "You rolled 3 d6,\nThe Results are:\n3, 5, 2"
	a := [3]int {3, 5, 2}

	r := createAnswerText(3, 6, a[:])

	checkIfTwoStringsAreTheSame(t, e, r)
}

// Check if two strings are the same
func checkIfTwoStringsAreTheSame(t *testing.T, e string, g string) {
	if e != g {
		t.Errorf("The strings are not the same, expected: \"%s\", got: \"%s\"", e, g)
	}
}
