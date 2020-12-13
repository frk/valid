package kw

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KW", A3: "KWT", Num: "414",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?965)[569][0-9]{7}$`),
	})
}
