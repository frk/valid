package sk

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 'SK'+10 digits (number must be divisible by 11)
	rxvat := regexp.MustCompile(`^SK[0-9]{10}$`)

	country.Add(country.Country{
		A2: "SK", A3: "SVK", Num: "703",
		Zip:   regexp.MustCompile(`^[0-9]{3}\s?[0-9]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?421)?[ ]?[1-9][0-9]{2}(?:[ ]?[0-9]{3}){2}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}

			n, _ := strconv.Atoi(v[2:])
			return (n % 11) == 0
		}),
	})
}
