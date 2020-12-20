package se

import (
	"regexp"

	"github.com/frk/isvalid/internal/algo"
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// VAT:
	//	12 digits, of which the last two are most often 01 e.g. SE999999999901.
	// 	(For sole proprietors who have several businesses the numbers can be 02,
	// 	03 and so on, since sole proprietors only have their personnummer as the
	// 	organisationsnummer. The first 10 digits are the same as the Swedish
	// 	organisationsnummer.
	//
	// 	The last digit of the first 10 digit number is the control digit, it is
	// 	calculated according to the Luhn algorithm.
	//
	// 	Reference:
	// 	- https://en.wikipedia.org/wiki/VAT_identification_number
	// 	- https://sv.wikipedia.org/wiki/Organisationsnummer
	rxvat := regexp.MustCompile(`^SE[0-9]{12}$`)

	country.Add(country.Country{
		A2: "SE", A3: "SWE", Num: "752",
		Zip:   regexp.MustCompile(`^[1-9][0-9]{2}\s?[0-9]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?46|0)[ \-]?7[ \-]?[02369](?:[ \-]?[0-9]){7}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			return algo.Luhn(v[2 : len(v)-2])
		}),
	})
}
