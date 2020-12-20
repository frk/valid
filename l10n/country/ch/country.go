package ch

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 6 digits (up to 31 December 2013). CHE 9 numeric digits plus TVA/MWST/IVA
	// e.g. CHE-123.456.788 TVA[20] The last digit is a MOD11 checksum digit build
	// with weighting pattern: 5,4,3,2,7,6,5,4
	// - https://en.wikipedia.org/wiki/VAT_identification_number
	rxvat := regexp.MustCompile(``) // TODO

	country.Add(country.Country{
		A2: "CH", A3: "CHE", Num: "756",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+41|0)7[5-9][0-9]{1,7}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			// TODO
			return false
		}),
	})
}
