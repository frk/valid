package be

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 'BE'+ 8 digits + 2 check digits – e.g. BE09999999XX
	var rxvat = regexp.MustCompile(`^BE[01][0-9]{9}$`)

	country.Add(country.Country{
		A2: "BE", A3: "BEL", Num: "056",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?32|0)4?[0-9]{8}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			v = v[2:]

			// mod 97
			num, _ := strconv.Atoi(v[:8])
			chk, _ := strconv.Atoi(v[8:])
			return num%97 == 97-chk
		}),
	})
}
