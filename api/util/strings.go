package util

import (
    "regexp"
    "strings"
)

// snakeCaseFirstCap matches the first char and can be used to turn it into lowercase
// and add an underscore between a leading special character and the rest of the string
var snakeCaseFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")

// snakeCaseAllCap matches the capitals coming after a capital character or digit.
var snakeCaseAllCap = regexp.MustCompile("([a-z\\d])([A-Z])")

// StringToSnakeCase turns the given string into snake-case.
//
// Examples:
//   - ThisIsATest => this_is_a_test
//   - $ThisIsATest => $_this_is_a_test
func StringToSnakeCase(str string) string {
    snake := snakeCaseFirstCap.ReplaceAllString(str, "${1}_${2}")
    snake = snakeCaseAllCap.ReplaceAllString(snake, "${1}_${2}")

    return strings.ToLower(snake)
}
