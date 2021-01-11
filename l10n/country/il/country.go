package il

import (
	"regexp"

	"github.com/frk/isvalid/internal/algo"
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 9 digit number. If the number of digits is less than 9, then zeros
	// should be padded to the left side. The leftmost digit is 5 for corporations.
	// Other leftmost digits are used for individuals. The rightmost digit
	// is a check digit (using Luhn algorithm).
	rxvat := regexp.MustCompile(`^[0-9]{9}$`)

	country.Add(country.Country{
		A2: "IL", A3: "ISR", Num: "376",
		Zip:   regexp.MustCompile(`^(?:[0-9]{5}|[0-9]{7})$`),
		Phone: regexp.MustCompile(`^(?:\+972|0)(?:[23489]|5[012345689]|77)[1-9][0-9]{6}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			return algo.Luhn(v)
		}),
	})
}
