package uz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 9 digits
	// Companies: 20000000X-29999999X
	// People: 40000000X-79999999X
	//
	// The taxpayer identification number consists of nine digits, the first
	// 8 digits are the taxpayer’s own number, and the last digit is the control
	// number.  The control number is formed by a certain algorithm in the process
	// of assigning a taxpayer identification number by a computer in the State
	// Tax Committee of the Republic of Uzbekistan.
	// - https://www.lex.uz/acts/499945
	rxvat := regexp.MustCompile(`^[24-7][0-9]{8}$`)

	country.Add(country.Country{
		A2: "UZ", A3: "UZB", Num: "860",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?998)?(?:6[125-79]|7[1-69]|88|9[0-9])[0-9]{7}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}

			// TODO if possible figure out what algorithm is used for the control digit

			return true
		}),
	})
}
