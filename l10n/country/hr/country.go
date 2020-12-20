package hr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 'HR'+ 11 digit number, e.g. HR12345678901
	rxvat := regexp.MustCompile(`^HR[0-9]{11}$`)

	country.Add(country.Country{
		A2: "HR", A3: "HRV", Num: "191",
		Zip: regexp.MustCompile(`^(?:[1-5][0-9]{4}$)`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			return country.ISO7064_MOD11_10(v[2:])
		}),
	})
}
