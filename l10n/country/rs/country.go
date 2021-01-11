package rs

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 9 digits (ex. 123456788) of which the first 8 are the actual ID number,
	// and the last digit is a checksum digit, calculated according to ISO 7064,
	// MOD 11-10
	rxvat := regexp.MustCompile(`^[0-9]{9}$`)

	country.Add(country.Country{
		A2: "RS", A3: "SRB", Num: "688",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+3816|06)[\- 0-9]{5,9}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			return country.ISO7064_MOD11_10(v)
		}),
	})
}
