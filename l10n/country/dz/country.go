package dz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "DZ", A3: "DZA", Num: "012",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?213|0)(?:5|6|7)[0-9]{8}$`),
	})
}
