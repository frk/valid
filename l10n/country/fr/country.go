package fr

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/internal/algo"
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// - 'FR'+ 2 digits (as validation key ) + 9 digits (as SIREN), the first and/or the
	//   second value can also be a character (any except O or I) - e.g. FRXX999999999
	// References:
	// - https://en.wikipedia.org/wiki/VAT_identification_number
	// - https://www.gov.uk/guidance/vat-eu-country-codes-vat-numbers-and-vat-in-other-languages
	rxvat := regexp.MustCompile(`^FR[A-HJ-NP-Z0-9]{2}[0-9]{9}$`)

	country.Add(country.Country{
		A2: "FR", A3: "FRA", Num: "250",
		Zip:   regexp.MustCompile(`^[0-9]{2}\s?[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:\+?33|0)[67][0-9]{8}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			v = v[2:]

			// The SIREN number ought to be valid Luhn: https://en.wikipedia.org/wiki/SIREN_code
			if !algo.Luhn(v[2:]) {
				return false
			}

			key, err := strconv.ParseInt(v[:2], 10, 32)
			if err != nil {
				return false
			}

			siren, err := strconv.ParseInt(v[2:], 10, 32)
			if err != nil {
				return false
			}

			// The validation key is calculated as follows:
			// [ 12 + 3 * ( SIREN modulo 97 ) ] modulo 97
			return key == (12+(3*(siren%97)))%97

			// NOTE(mkopriva): couldn't find anything official and/or
			// open source that demonstrates how to handle a validation
			// key that contains letters ... this will fail if such a
			// key is provided even if it is valid.
		}),
	})
}
