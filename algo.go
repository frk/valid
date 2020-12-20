package isvalid

// Luhn validates the given string v using the Luhn (mod 10) algorithm, the string
// v is expected to contain only digits. (https://en.wikipedia.org/wiki/Luhn_algorithm)
func Luhn(v string) bool {
	var sum int
	var double bool
	for i := len(v) - 1; i >= 0; i-- {
		num := btoi(v[i])

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
