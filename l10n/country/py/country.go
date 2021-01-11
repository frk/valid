package py

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 6 to 8 digits, 1 dash, 1 check sum digit
	rxvat := regexp.MustCompile(`^[0-9]{6,8}-[0-9]$`)

	country.Add(country.Country{
		A2: "PY", A3: "PRY", Num: "600",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?595|0)9[9876][0-9]{7}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			// TODO
			return false
		}),
	})
}
