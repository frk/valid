package tz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TZ", A3: "TZA", Num: "834",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?255|0)?[67][0-9]{8}$`),
	})
}
