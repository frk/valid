package algo

import (
	"strconv"
)

// Luhn validates the given string v using the Luhn (mod 10) algorithm.
// - https://en.wikipedia.org/wiki/Luhn_algorithm
//
// The string v is assumed to contain only digits.
func Luhn(v string) bool {
	var sum int
	var double bool
	for i := len(v) - 1; i >= 0; i-- {
		num, _ := strconv.Atoi(string(v[i]))

		if double {
			num *= 2
			if num > 9 {
				num = (num % 10) + 1
			}
		}
		double = !double

		sum += num
	}
	return sum%10 == 0
}
