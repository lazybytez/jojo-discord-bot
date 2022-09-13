package dice

import (
	"reflect"
	"testing"
)

func TestArrayIntToArrayString(t *testing.T) {
	p := [5]int{4, 20, -13, 0, -6}
	e := [5]string{"4", "20", "-13", "0", "-6"}

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

func TestCreateAnswerResultContent(t *testing.T) {
	a := [10]string{"What", "is", "'Courage'?", "Courage", "is", "owning", "your", "feat!", "-", "Zeppeli"}
	e := "What, is, 'Courage'?, Courage, is, owning, your, feat!, -, Zeppeli"

	r := createAnswerResultContent(a[:])

	checkIfTwoStringsAreTheSame(t, e, r, "contents")
}

func TestCreateAnswerResultTitleOne(t *testing.T) {
	e := "The Result is"

	r := createAnswerResultTitle(1)

	checkIfTwoStringsAreTheSame(t, e, r, "titles (one)")
}

func TestCreateAnswerResultTitleMultiple(t *testing.T) {
	e := "The Results are"

	r := createAnswerResultTitle(2)

	checkIfTwoStringsAreTheSame(t, e, r, "titles (multiple)")
}

func TestCreateAnswerTitle(t *testing.T) {
	e := "You rolled 3 d6"

	r := createAnswerTitle(3, 6)

	checkIfTwoStringsAreTheSame(t, e, r, "titles")
}

// Check if two strings are the same, say in the error message what (n) was not the same
func checkIfTwoStringsAreTheSame(t *testing.T, e string, g string, n string) {
	if e != g {
		t.Errorf("The %s are not the same, expected: \"%s\", got: \"%s\"", n, e, g)
	}
}

// Check if two ints are the same, say in the error message what (n) was not the same
func checkIfTwoIntAreTheSame(t *testing.T, e int, g int, n string) {
	if e != g {
		t.Errorf("The %s are not the same, expected: %d, got %d", n, e, g)
	}
}

// Check if two boolean are the same, say in the error message what (n) was not the same
func checkIfTwoBoolAreTheSame(t *testing.T, e bool, g bool, n string) {
	if e != g {
		t.Errorf("The %s are not the same, expected: %t, got %t", n, e, g)
	}
}

// Check if two objects are the same in depth. Say in the error message what (n) the objects are
func checkIfTwoObjectsAreTheSameInDepth(t *testing.T, e, g any, n string) {
	if !reflect.DeepEqual(e, g) {
		t.Errorf("The two %s are not the same", n)
	}
}
