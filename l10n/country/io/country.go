package io

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IO", A3: "IOT", Num: "086",
		Zip: regexp.MustCompile(`^BBND 1ZZ$`),
	})
}
